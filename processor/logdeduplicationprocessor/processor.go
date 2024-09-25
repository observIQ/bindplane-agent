// Copyright  observIQ, Inc.
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

package logdeduplicationprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/observiq/bindplane-agent/expr"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

// logDedupProcessor is a logDedupProcessor that counts duplicate instances of logs.
type logDedupProcessor struct {
	emitInterval    time.Duration
	condition       *expr.OTTLCondition[ottllog.TransformContext]
	conditionString string
	aggregator      *logAggregator
	remover         *fieldRemover
	consumer        consumer.Logs
	logger          *zap.Logger
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	mux             sync.Mutex
}

func newProcessor(cfg *Config, condition *expr.OTTLCondition[ottllog.TransformContext], consumer consumer.Logs, logger *zap.Logger) (*logDedupProcessor, error) {
	// This should not happen due to config validation but we check anyways.
	timezone, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %w", err)
	}

	return &logDedupProcessor{
		emitInterval:    cfg.Interval,
		condition:       condition,
		conditionString: cfg.Condition,
		aggregator:      newLogAggregator(cfg.LogCountAttribute, timezone),
		remover:         newFieldRemover(cfg.ExcludeFields),
		consumer:        consumer,
		logger:          logger,
	}, nil
}

// Start starts the processor.
func (p *logDedupProcessor) Start(_ context.Context, _ component.Host) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	p.wg.Add(1)
	go p.handleExportInterval(ctx)

	return nil
}

// Capabilities returns the consumer's capabilities.
func (p *logDedupProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Shutdown stops the processor.
func (p *logDedupProcessor) Shutdown(ctx context.Context) error {
	if p.cancel != nil {
		p.cancel()
	}

	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		p.wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-doneChan:
		return nil
	}
}

// ConsumeLogs processes the logs.
func (p *logDedupProcessor) ConsumeLogs(ctx context.Context, pl plog.Logs) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	for i := 0; i < pl.ResourceLogs().Len(); i++ {
		resourceLogs := pl.ResourceLogs().At(i)
		resourceAttrs := resourceLogs.Resource().Attributes()
		for j := 0; j < resourceLogs.ScopeLogs().Len(); j++ {
			scope := resourceLogs.ScopeLogs().At(j)
			logs := scope.LogRecords()
			logs.RemoveIf(func(logRecord plog.LogRecord) bool {
				var match bool
				if p.conditionString == "true" || p.conditionString == "" {
					match = true
				} else {
					logCtx := ottllog.NewTransformContext(
						logRecord,
						scope.Scope(),
						resourceLogs.Resource(),
						scope,
						resourceLogs,
					)
					logMatch, err := p.condition.Match(ctx, logCtx)
					match = err == nil && logMatch
				}
				// only aggregate logs that match condition
				if match {
					// Remove excluded fields if any
					p.remover.RemoveFields(logRecord)

					// Add the log to the aggregator
					p.aggregator.Add(resourceAttrs, logRecord)
				}
				return match
			})
		}
	}

	// immediately consume any logs that didn't match the condition
	if pl.LogRecordCount() > 0 {
		err := p.consumer.ConsumeLogs(ctx, pl)
		if err != nil {
			p.logger.Error("failed to consume logs", zap.Error(err))
		}
	}

	return nil
}

// handleExportInterval sends metrics at the configured interval.
func (p *logDedupProcessor) handleExportInterval(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.emitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.mux.Lock()

			logs := p.aggregator.Export()
			// Only send logs if we have some
			if logs.LogRecordCount() > 0 {
				err := p.consumer.ConsumeLogs(ctx, logs)
				if err != nil {
					p.logger.Error("failed to consume logs", zap.Error(err))
				}
			}
			p.aggregator.Reset()
			p.mux.Unlock()
		}
	}
}
