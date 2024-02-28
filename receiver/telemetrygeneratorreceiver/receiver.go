// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetrygeneratorreceiver //import "github.com/observiq/bindplane-agent/receiver/telemetrygeneratorreceiver"

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"

	"go.uber.org/zap"
)

type generator interface {
	generate() error
}
type telemetryGeneratorReceiver struct {
	logger             *zap.Logger
	id                 component.ID
	cfg                *Config
	supportedTelemetry component.DataType
	doneChan           chan struct{}
	ctx                context.Context
	cancelFunc         context.CancelCauseFunc
	generator          generator
}

// newTelemetryGeneratorReceiver creates a new rehydration receiver
func newTelemetryGeneratorReceiver(ctx context.Context, id component.ID, logger *zap.Logger, cfg *Config, g generator) (telemetryGeneratorReceiver, error) {
	ctx, cancel := context.WithCancelCause(ctx)

	return telemetryGeneratorReceiver{
		logger:     logger,
		id:         id,
		cfg:        cfg,
		doneChan:   make(chan struct{}),
		ctx:        ctx,
		cancelFunc: cancel,
		generator:  g,
	}, nil
}

// Shutdown shuts down the telemetry generator receiver
func (r *telemetryGeneratorReceiver) Shutdown(ctx context.Context) error {
	r.cancelFunc(errors.New("shutdown"))
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-r.doneChan:
	}

	return err
}

// Start starts the logsGeneratorReceiver receiver
func (r *telemetryGeneratorReceiver) Start(ctx context.Context, _ component.Host) error {

	go func() {
		defer close(r.doneChan)

		ticker := time.NewTicker(time.Second / time.Duration(r.cfg.PayloadsPerSecond))
		defer ticker.Stop()

		// Call once before the loop to ensure we do a collection before the first ticker
		r.generator.generate()
		for {
			select {
			case <-ctx.Done():
				return
			case <-r.ctx.Done():
				return
			case <-ticker.C:
				r.generator.generate()
			}
		}
	}()
	return nil
}
