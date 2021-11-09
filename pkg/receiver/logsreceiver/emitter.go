// Copyright The OpenTelemetry Authors
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

package logsreceiver

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/open-telemetry/opentelemetry-log-collection/entry"
	"github.com/open-telemetry/opentelemetry-log-collection/operator"
	"github.com/open-telemetry/opentelemetry-log-collection/operator/helper"
	"go.uber.org/zap"
)

// LogEmitter is a stanza operator that emits log entries to a channel
type LogEmitter struct {
	helper.OutputOperator
	logChan               chan []*entry.Entry
	stopOnce              sync.Once
	cancel                context.CancelFunc
	batchMux              sync.Mutex
	batch                 []*entry.Entry
	wg                    sync.WaitGroup
	flushTriggerAmount    uint
	flushInterval         time.Duration
	entryBatchInitialSize int
}

var (
	defaultFlushInterval           = 100 * time.Millisecond
	defaultFlushTriggerAmount uint = 200
)

// NewLogEmitter creates a new receiver output
// TODO: Convert args here to functional options
func NewLogEmitter(logger *zap.SugaredLogger, flushInterval time.Duration, flushTriggerAmount uint) *LogEmitter {
	if flushInterval == 0 {
		flushInterval = defaultFlushInterval
	}

	if flushTriggerAmount == 0 {
		flushTriggerAmount = defaultFlushTriggerAmount
	}

	entryBatchInitialSize := int(math.Max(float64(flushTriggerAmount), float64(flushTriggerAmount)*1.1))

	return &LogEmitter{
		OutputOperator: helper.OutputOperator{
			BasicOperator: helper.BasicOperator{
				OperatorID:    "log_emitter",
				OperatorType:  "log_emitter",
				SugaredLogger: logger,
			},
		},
		logChan:               make(chan []*entry.Entry),
		batch:                 make([]*entry.Entry, 0, entryBatchInitialSize),
		flushInterval:         flushInterval,
		flushTriggerAmount:    flushTriggerAmount,
		entryBatchInitialSize: entryBatchInitialSize,
	}
}

func (e *LogEmitter) Start(_ operator.Persister) error {
	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel

	e.wg.Add(1)
	go e.flusher(ctx)
	return nil
}

// Process will emit an entry to the output channel
func (e *LogEmitter) Process(ctx context.Context, ent *entry.Entry) error {
	e.batchMux.Lock()

	e.batch = append(e.batch, ent)
	if uint(len(e.batch)) >= e.flushTriggerAmount {
		// flushTriggerAmount triggers a flush
		e.batchMux.Unlock()
		e.flush(ctx)
		return nil
	}

	e.batchMux.Unlock()
	return nil
}

func (e *LogEmitter) flusher(ctx context.Context) {
	defer e.wg.Done()

	ticker := time.NewTicker(e.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.flush(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (e *LogEmitter) flush(ctx context.Context) {
	var batch []*entry.Entry
	e.batchMux.Lock()

	if len(e.batch) == 0 {
		e.batchMux.Unlock()
		return
	}
	batch = e.batch
	e.batch = make([]*entry.Entry, 0, e.entryBatchInitialSize)

	e.batchMux.Unlock()

	select {
	case e.logChan <- batch:
	case <-ctx.Done():
	}
}

// Stop will close the log channel
func (e *LogEmitter) Stop() error {
	e.stopOnce.Do(func() {
		close(e.logChan)
	})

	if e.cancel != nil {
		e.cancel()
		e.cancel = nil
	}

	e.wg.Wait()
	return nil
}
