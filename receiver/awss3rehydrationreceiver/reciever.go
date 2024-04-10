// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package awss3rehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/awss3rehydrationreceiver"

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/observiq/bindplane-agent/receiver/awss3rehydrationreceiver/internal/s3"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/extension/experimental/storage"
	"go.uber.org/zap"
)

// checkpointStorageKey is the key used for storing the checkpoint
const checkpointStorageKey = "aws_s3_checkpoint"

// newAWSS3Client is the function used to create new AWS S3 clients.
// Meant to be overwritten for tests
var newAWSS3Client = s3.NewAWSClient

type rehydrationReceiver struct {
	logger             *zap.Logger
	id                 component.ID
	cfg                *Config
	awsClient          s3.S3Client
	supportedTelemetry component.DataType
	storageClient      storage.Client

	startingTime time.Time
	endingTime   time.Time

	doneChan   chan struct{}
	started    bool
	ctx        context.Context
	cancelFunc context.CancelCauseFunc
}

// newMetricsReceiver creates a new metrics specific receiver.
func newMetricsReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Metrics) (*rehydrationReceiver, error) {
	r, err := newRehydrationReceiver(id, logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeMetrics
	// r.consumer = newMetricsConsumer(nextConsumer)

	return r, nil
}

// newLogsReceiver creates a new logs specific receiver.
func newLogsReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Logs) (*rehydrationReceiver, error) {
	r, err := newRehydrationReceiver(id, logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeLogs
	// r.consumer = newLogsConsumer(nextConsumer)

	return r, nil
}

// newTracesReceiver creates a new traces specific receiver.
func newTracesReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Traces) (*rehydrationReceiver, error) {
	r, err := newRehydrationReceiver(id, logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeTraces
	// r.consumer = newTracesConsumer(nextConsumer)

	return r, nil
}

// newRehydrationReceiver creates a new rehydration receiver
func newRehydrationReceiver(id component.ID, logger *zap.Logger, cfg *Config) (*rehydrationReceiver, error) {
	awsClient := newAWSS3Client(cfg.Region, cfg.RoleArn)

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
		logger:        logger,
		id:            id,
		cfg:           cfg,
		awsClient:     awsClient,
		doneChan:      make(chan struct{}),
		storageClient: storage.NewNopClient(),
		startingTime:  startingTime,
		endingTime:    endingTime,
		ctx:           ctx,
		cancelFunc:    cancel,
	}, nil
}

// Start starts the rehydration receiver
func (r *rehydrationReceiver) Start(ctx context.Context, host component.Host) error {

	// if r.cfg.StorageID != nil {
	// 	storageClient, err := getStorageClient(ctx, host, r.cfg.StorageID, r.id, r.supportedTelemetry)
	// 	if err != nil {
	// 		return fmt.Errorf("getStorageClient: %w", err)
	// 	}

	// 	r.storageClient = storageClient
	// }

	// r.started = true
	// go r.scrape()
	return nil
}

// Shutdown shuts down the rehydration receiver
func (r *rehydrationReceiver) Shutdown(ctx context.Context) error {
	r.cancelFunc(errors.New("shutdown"))
	var err error

	// If we have called started then close and wait for goroutine to finish
	if r.started {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-r.doneChan:
		}
	}

	err = errors.Join(err, r.storageClient.Close(ctx))

	return err
}
