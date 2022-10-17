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

package logsummaryprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type processor struct {
	interval time.Duration
	cancel   context.CancelFunc
	summary  *summary
	mux      sync.Mutex

	logger   *zap.Logger
	config   *Config
	exporter component.MetricsExporter
}

// newProcessor returns a new processor.
func newProcessor(logger *zap.Logger, config *Config) *processor {
	return &processor{
		summary: newSummary(),
		logger:  logger,
		config:  config,
	}
}

// Start starts the processor.
func (p *processor) Start(_ context.Context, host component.Host) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	exporters, ok := host.GetExporters()[config.MetricsDataType]
	if !ok {
		return fmt.Errorf("exporters of type %s do not exist", config.MetricsDataType)
	}

	exporter, ok := exporters[p.config.Exporter]
	if !ok {
		return fmt.Errorf("exporter with id %s does not exist", p.config.Exporter)
	}

	p.exporter = exporter.(component.MetricsExporter)
	go p.sendMetrics(ctx)

	return nil
}

// Capabilities returns the consumer's capabilities.
func (p *processor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// Shutdown stops the processor.
func (p *processor) Shutdown(_ context.Context) error {
	p.cancel()
	return nil
}

// handleMetrics handles sending metrics on the configured interval.
func (p *processor) handleMetrics(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.sendMetrics(ctx)
		}
	}
}

// sendMetrics sends the summary as metrics to the configured exporter.
func (p *processor) sendMetrics(ctx context.Context) {
	p.mux.Lock()
	defer p.mux.Unlock()

	metrics := p.summary.toMetrics()
	p.summary.reset()

	if metrics.DataPointCount() > 0 {
		_ = p.exporter.ConsumeMetrics(ctx, metrics)
	}
}

// ConsumeLogs consumes logs and summarizes them as metrics.
func (p *processor) ConsumeLogs(_ context.Context, pl plog.Logs) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	err := p.summary.update(pl)
	return err
}
