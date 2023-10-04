package azureblobrehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver"

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

// azurePathSeparator the path separator that Azure storage uses
const azurePathSeparator = "/"

var errInvalidBlobPath = errors.New("invalid blob path")

type rehydrationReceiver struct {
	logger             *zap.Logger
	cfg                *Config
	azureClient        blobClient
	supportedTelemetry component.DataType
	consumer           blobConsumer

	startingTime time.Time
	endingTime   time.Time

	doneChan   chan struct{}
	ctx        context.Context
	cancelFunc context.CancelCauseFunc
}

// newMetricsReceiver creates a new metrics specific receiver.
func newMetricsReceiver(logger *zap.Logger, cfg *Config, nextConsumer consumer.Metrics) (*rehydrationReceiver, error) {
	r, err := newRehydrationReceiver(logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeMetrics
	r.consumer = newMetricsConsumer(nextConsumer)

	return r, nil
}

// newLogsReceiver creates a new logs specific receiver.
func newLogsReceiver(logger *zap.Logger, cfg *Config, nextConsumer consumer.Logs) (*rehydrationReceiver, error) {
	r, err := newRehydrationReceiver(logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeLogs
	r.consumer = newLogsConsumer(nextConsumer)

	return r, nil
}

// newTracesReceiver creates a new traces specific receiver.
func newTracesReceiver(logger *zap.Logger, cfg *Config, nextConsumer consumer.Traces) (*rehydrationReceiver, error) {
	r, err := newRehydrationReceiver(logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeTraces
	r.consumer = newTracesConsumer(nextConsumer)

	return r, nil
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

// Start starts the rehydration receiver
func (r *rehydrationReceiver) Start(_ context.Context, host component.Host) error {
	go r.rehydrateBlobs()
	return nil
}

// Shutdown shuts down the rehydration receiver
func (r *rehydrationReceiver) Shutdown(ctx context.Context) error {
	r.cancelFunc(errors.New("shutdown"))
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.doneChan:
		return nil
	}
}

// rehydrateBlobs pulls blob paths from the UI and if they are within the specified
// time range then the blobs will be downloaded and rehydrated.
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
			blobs, nextMarker, err := r.azureClient.ListBlobs(r.ctx, r.cfg.Container, prefix, marker)
			if err != nil {
				r.logger.Error("Failed to list blobs", zap.Error(err))
				continue
			}

			marker = nextMarker

			for _, blob := range blobs {
				blobTime, telemetryType, err := r.parseBlobPath(prefix, blob.Name)
				switch {
				case errors.Is(err, errInvalidBlobPath):
					r.logger.Debug("Skipping Blob, non-matching blob path", zap.String("blob", blob.Name))
				case err != nil:
					r.logger.Error("Error processing blob path", zap.String("blob", blob.Name), zap.Error(err))
				default:
					// if the blob is not in the specified time range or not of the telemetry type supported by this receiver
					// then skip consuming it.
					if !r.isInTimeRange(*blobTime) || telemetryType != r.supportedTelemetry {
						continue
					}

					// Process and consume the blob at the given path
					if err := r.processBlob(blob); err != nil {
						r.logger.Error("Error consuming blob", zap.String("blob", blob.Name), zap.Error(err))
					}
				}
			}

		}
	}
}

// constants for blob path parts
const (
	year   = "year="
	month  = "month="
	day    = "day="
	hour   = "hour="
	minute = "minute="
)

// strings that indicate what type of telemetry is in a blob
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
	for ; i < len(parts)-1; i++ {
		part := parts[i]

		if !strings.HasPrefix(part, nextExpectedPart) {
			err = errInvalidBlobPath
			return
		}

		val := strings.TrimPrefix(part, nextExpectedPart)

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

	// Special case when using hour granularity.
	// There won't be a minute=XX part of the path if we've exited the loop
	// and we still expect minutes just write ':00' for minutes.
	if nextExpectedPart == minute {
		tsBuilder.WriteString(":00")
	}

	// Parse the expected format
	parsedTime, timeErr := time.Parse(timeFormat, tsBuilder.String())
	if err != nil {
		err = fmt.Errorf("parse blob time: %w", timeErr)
		return
	}
	blobTime = &parsedTime

	// For the last part of the path parse the telemetry type
	lastPart := parts[len(parts)-1]
	switch {
	case strings.Contains(lastPart, metricBlobSignifier):
		telemetryType = component.DataTypeMetrics
	case strings.Contains(lastPart, logsBlobSignifier):
		telemetryType = component.DataTypeLogs
	case strings.Contains(lastPart, tracesBlobSignifier):
		telemetryType = component.DataTypeTraces
	}

	return
}

// isInTimeRange returns true if startingTime <= blobTime <= endingTime
func (r *rehydrationReceiver) isInTimeRange(blobTime time.Time) bool {
	return (blobTime.Equal(r.startingTime) || blobTime.After(r.startingTime)) &&
		(blobTime.Equal(r.endingTime) || blobTime.Before(r.endingTime))
}

// processBlob does the following:
// 1. Downloads the blob
// 2. Decompresses the blob if applicable
// 3. Pass the blob to the consumer
func (r *rehydrationReceiver) processBlob(blob *blobInfo) error {
	// Allocate a buffer the size of the blob. If the buffer isn't big enough download errors.
	blobBuffer := make([]byte, blob.Size)

	size, err := r.azureClient.DownloadBlob(r.ctx, r.cfg.Container, blob.Name, blobBuffer)
	if err != nil {
		return fmt.Errorf("download blob: %w", err)
	}

	// Check file extension see if we need to decompress
	ext := filepath.Ext(blob.Name)
	switch {
	case ext == ".gz":
		blobBuffer, err = gzipDecompress(blobBuffer[:size])
		if err != nil {
			return fmt.Errorf("gzip: %w", err)
		}
	case ext == ".json":
		// Do nothing for json files
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	if err := r.consumer.Consume(r.ctx, blobBuffer); err != nil {
		return fmt.Errorf("consume: %w", err)
	}
	return nil
}

// gzipDecompress does a gzip decompression on the passed in contents
func gzipDecompress(contents []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(contents))
	if err != nil {
		return nil, fmt.Errorf("new reader: %w", err)
	}

	result, err := io.ReadAll(gr)
	if err != nil {
		return nil, fmt.Errorf("decompression: %w", err)
	}

	return result, nil
}
