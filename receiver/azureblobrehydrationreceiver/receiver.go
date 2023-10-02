package azureblobrehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver"

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

var (
	errMissingHost = errors.New("nil host")
)

type rehydrationReceiver struct {
	logger      *zap.Logger
	cfg         *Config
	azureClient blobClient

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

	ctx, cancel := context.WithCancelCause(context.Background())

	return &rehydrationReceiver{
		logger:      logger,
		cfg:         cfg,
		azureClient: azureClient,
		doneChan:    make(chan struct{}),
		ctx:         ctx,
		cancelFunc:  cancel,
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
			blobs, nextMarker, err := r.azureClient.ListBlobs(r.ctx, r.cfg.Container, prefix, marker)
			if err != nil {
				r.logger.Error("Failed while listing blobs", zap.Error(err))
			}

			for _, blob := range blobs {
				r.logger.Info(blob)
			}

			marker = nextMarker
		}
	}
}
