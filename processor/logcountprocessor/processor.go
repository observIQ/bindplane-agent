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

package logcountprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type processor struct {
	config    *Config
	evaluator *Evaluator
	counter   *Counter
	exporter  component.MetricsExporter
	consumer  consumer.Logs
	logger    *zap.Logger
	cancel    context.CancelFunc
	mux       sync.Mutex
}

// newProcessor returns a new processor.
func newProcessor(config *Config, consumer consumer.Logs, logger *zap.Logger) *processor {
	return &processor{
		config:   config,
		counter:  NewCounter(logger),
		consumer: consumer,
		logger:   logger,
	}
}

// Start starts the processor.
func (p *processor) Start(_ context.Context, host component.Host) error {
	exporter, err := p.getExporter(host)
	if err != nil {
		return fmt.Errorf("failed to get configured exporter: %w", err)
	}
	p.exporter = exporter

	evaluator, err := p.createEvaluator()
	if err != nil {
		return fmt.Errorf("failed to create evaluator: %w", err)
	}
	p.evaluator = evaluator

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	go p.handleMetricInterval(ctx)

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

// ConsumeLogs adds log records to the counter if they match the defined expression.
func (p *processor) ConsumeLogs(ctx context.Context, pl plog.Logs) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	for i := 0; i < pl.ResourceLogs().Len(); i++ {
		resourceLogs := pl.ResourceLogs().At(i)
		resource := resourceLogs.Resource()
		for j := 0; j < resourceLogs.ScopeLogs().Len(); j++ {
			logs := resourceLogs.ScopeLogs().At(j).LogRecords()
			for k := 0; k < logs.Len(); k++ {
				p.evaluateLog(resource, logs.At(k))
			}
		}
	}

	return p.consumer.ConsumeLogs(ctx, pl)
}

// evaluateLog evaluates an incoming log and updates the counter.
func (p *processor) evaluateLog(resource pcommon.Resource, log plog.LogRecord) {
	if !p.evaluator.MatchesLog(resource, log) {
		return
	}

	resourceAttributes := resource.Attributes().AsRaw()
	attributes := p.evaluator.GetAttributes(resource, log)
	record := NewRecord(resourceAttributes, attributes)
	p.counter.Add(record)
}

// handleMetricInterval handles sending metrics on the configured interval.
func (p *processor) handleMetricInterval(ctx context.Context) {
	ticker := time.NewTicker(p.config.Interval)
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

	if len(p.counter.counts) == 0 {
		return
	}

	metrics := p.getMetrics()
	p.counter.Reset()

	err := p.exporter.ConsumeMetrics(ctx, metrics)
	if err != nil {
		p.logger.Error("Failed to export metrics", zap.Error(err))
		return
	}
}

// getMetrics creates metrics from the current counter.
func (p *processor) getMetrics() pmetric.Metrics {
	metrics := pmetric.NewMetrics()
	for _, count := range p.counter.counts {
		resourceMetrics := metrics.ResourceMetrics().AppendEmpty()
		resourceMetrics.Resource().Attributes().FromRaw(count.record.Resource)
		scopeMetrics := resourceMetrics.ScopeMetrics().AppendEmpty()
		scopeMetrics.Scope().SetName(typeStr)

		metrics := scopeMetrics.Metrics().AppendEmpty()
		metrics.SetName(p.config.MetricName)
		metrics.SetUnit(p.config.MetricUnit)
		metrics.SetEmptyGauge()

		gauge := metrics.Gauge().DataPoints().AppendEmpty()
		gauge.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		gauge.SetIntValue(int64(count.value))
		gauge.Attributes().FromRaw(count.record.Attributes)
	}

	return metrics
}

// getExporter retrieves the configured exporter from the host.
func (p *processor) getExporter(host component.Host) (component.MetricsExporter, error) {
	exporters, ok := host.GetExporters()[config.MetricsDataType]
	if !ok {
		return nil, fmt.Errorf("exporter with id %s is not configured in a metric pipeline", p.config.Exporter)
	}

	exporter, ok := exporters[p.config.Exporter]
	if !ok {
		return nil, fmt.Errorf("exporter with id %s does not exist", p.config.Exporter)
	}

	return exporter.(component.MetricsExporter), nil
}

// createEvaluator creates an evaluator from the configured expressions.
func (p *processor) createEvaluator() (*Evaluator, error) {
	matchExpr, err := expr.Compile(p.config.Match, expr.AsBool(), expr.AllowUndefinedVariables())
	if err != nil {
		return nil, fmt.Errorf("failed to compile match expression: %w", err)
	}

	attrsExpr := map[string]*vm.Program{}
	for key, value := range p.config.Attributes {
		expr, err := expr.Compile(value, expr.AllowUndefinedVariables())
		if err != nil {
			return nil, fmt.Errorf("failed to compile attribute expression (%s): %w", key, err)
		}
		attrsExpr[key] = expr
	}

	return NewEvaluator(matchExpr, attrsExpr, p.logger), nil
}
