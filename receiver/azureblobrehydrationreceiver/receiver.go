package azureblobrehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver"

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// azurePathSeparator the path separator that Azure storage uses
const azurePathSeparator = "/"

var (
	errMissingHost     = errors.New("nil host")
	errInvalidBlobPath = errors.New("invalid blob path")
)

type rehydrationReceiver struct {
	logger      *zap.Logger
	cfg         *Config
	azureClient blobClient

	startingTime time.Time
	endingTime   time.Time

	doneChan   chan struct{}
	ctx        context.Context
	cancelFunc context.CancelCauseFunc
}

// newRehydrationReceiver creates a new rehydration receiver
func newRehydrationReceiver(logger *zap.Logger, cfg *Config) (*rehydrationReceiver, error) {
	azureClient, err := newAzureBlobClient(cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("new Azure client: %w", err)
	}

	// We should not get an error for either of these time parsings as we check in config validate.
	// Doing error checking anyways just in case.
	startingTime, err := time.Parse(timeFormat, cfg.StartingTime)
	if err != nil {
		return nil, fmt.Errorf("invalid starting_time timestamp: %w", err)
	}

	endingTime, err := time.Parse(timeFormat, cfg.EndingTime)
	if err != nil {
		return nil, fmt.Errorf("invalid ending_time timestamp: %w", err)
	}

	ctx, cancel := context.WithCancelCause(context.Background())

	return &rehydrationReceiver{
		logger:       logger,
		cfg:          cfg,
		azureClient:  azureClient,
		doneChan:     make(chan struct{}),
		startingTime: startingTime,
		endingTime:   endingTime,
		ctx:          ctx,
		cancelFunc:   cancel,
	}, nil
}

func (r *rehydrationReceiver) Start(_ context.Context, host component.Host) error {
	if host == nil {
		return errMissingHost
	}

	go r.rehydrateBlobs()
	return nil
}

func (r *rehydrationReceiver) Shutdown(ctx context.Context) error {
	r.cancelFunc(errors.New("shutdown"))
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.doneChan:
		return nil
	}
}

func (r *rehydrationReceiver) rehydrateBlobs() {
	defer close(r.doneChan)

	var prefix *string
	if r.cfg.RootFolder != "" {
		prefix = &r.cfg.RootFolder
	}

	var marker *string
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			blobPaths, nextMarker, err := r.azureClient.ListBlobs(r.ctx, r.cfg.Container, prefix, marker)
			if err != nil {
				r.logger.Error("Failed to list blobs", zap.Error(err))
				continue
			}

			for _, blobPath := range blobPaths {
				blobTime, _, err := r.parseBlobPath(prefix, blobPath)
				switch {
				case errors.Is(err, errInvalidBlobPath):
					r.logger.Debug("Skipping Blob, non-matching blob path", zap.String("blob", blobPath))
				case err != nil:
					r.logger.Error("Error processing blob path", zap.String("blob", blobPath), zap.Error(err))
				default:
					if r.isInTimeRange(*blobTime) {

					}
				}
			}

			marker = nextMarker
		}
	}
}

// constants for blob path parts
const (
	year   = "year"
	month  = "month"
	day    = "day"
	hour   = "hour"
	minute = "minute"
)

// Strings that indicate what type of telemetry is in a blob
const (
	metricBlobSignifier = "metrics_"
	logsBlobSignifier   = "logs_"
	tracesBlobSignifier = "traces_"
)

// parseBlobPath returns true if the blob is within the existing time range
func (r *rehydrationReceiver) parseBlobPath(prefix *string, blobName string) (blobTime *time.Time, telemetryType component.DataType, err error) {
	parts := strings.Split(blobName, azurePathSeparator)

	if len(parts) == 0 {
		err = errInvalidBlobPath
		return
	}

	// Build timestamp in 2006-01-02T15:04 format
	tsBuilder := strings.Builder{}

	i := 0
	// If we have a prefix start looking at the second part of the path
	if prefix != nil {
		i = 1
	}

	nextExpectedPart := year
	for ; i < len(parts); i++ {
		part := parts[i]

		// For the last part of the path parse the telemetry type
		if i == len(parts)-1 {
			// Special case to catch when using hour granularity.
			// There won't be a minute=XX part of the path so if we're on the last
			// part and we still expect minutes just write ':00' for minutes.
			if nextExpectedPart == minute {
				tsBuilder.WriteString(":00")
			}

			switch {
			case strings.Contains(part, metricBlobSignifier):
				telemetryType = component.DataTypeMetrics
			case strings.Contains(part, logsBlobSignifier):
				telemetryType = component.DataTypeLogs
			case strings.Contains(part, tracesBlobSignifier):
				telemetryType = component.DataTypeTraces
			}
		}

		if !strings.HasPrefix(part, nextExpectedPart) {
			err = errInvalidBlobPath
			return
		}

		val := strings.TrimPrefix(part, fmt.Sprintf("%s=", nextExpectedPart))

		switch nextExpectedPart {
		case year:
			nextExpectedPart = month
		case month:
			tsBuilder.WriteString("-")
			nextExpectedPart = day
		case day:
			tsBuilder.WriteString("-")
			nextExpectedPart = hour
		case hour:
			tsBuilder.WriteString("T")
			nextExpectedPart = minute
		case minute:
			tsBuilder.WriteString(":")
			nextExpectedPart = ""
		}

		tsBuilder.WriteString(val)
	}

	// Parse the expected format
	*blobTime, err = time.Parse(timeFormat, tsBuilder.String())
	if err != nil {
		err = fmt.Errorf("parse blob time: %w", err)
		return
	}

	return
}

// isInTimeRange returns true if startingTime <= blobTime <= endingTime
func (r *rehydrationReceiver) isInTimeRange(blobTime time.Time) bool {
	return (blobTime.Equal(r.startingTime) || blobTime.After(r.startingTime)) &&
		(blobTime.Equal(r.endingTime) || blobTime.Before(r.endingTime))
}
