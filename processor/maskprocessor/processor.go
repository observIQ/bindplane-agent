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

package maskprocessor

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

const (
	resourceField   = "resource"
	attributesField = "attributes"
	bodyField       = "body"
)

// processor is the processor used to mask data.
type processor struct {
	logger *zap.Logger
	cfg    *Config
	rules  map[string]*regexp.Regexp
}

// newProcessor creates a new mask processor.
func newProcessor(logger *zap.Logger, cfg *Config) *processor {
	return &processor{
		logger: logger,
		cfg:    cfg,
	}
}

// start is used to start the processor.
func (p *processor) start(context.Context, component.Host) error {
	rules, err := p.cfg.CompileRules()
	if err != nil {
		return err
	}

	p.rules = rules
	return nil
}

// processLogs masks incoming logs.
func (p *processor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resource := ld.ResourceLogs().At(i)
		resourceMaskFunc := p.createMaskFunc(resourceField)
		resource.Resource().Attributes().Range(resourceMaskFunc)
		for j := 0; j < resource.ScopeLogs().Len(); j++ {
			scope := resource.ScopeLogs().At(j)
			for k := 0; k < scope.LogRecords().Len(); k++ {
				logs := scope.LogRecords().At(k)
				attrMaskFunc := p.createMaskFunc(attributesField)
				logs.Attributes().Range(attrMaskFunc)
				p.maskValue(bodyField, logs.Body())
			}
		}
	}

	return ld, nil
}

// processTraces masks incoming traces.
func (p *processor) processTraces(_ context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	for i := 0; i < td.ResourceSpans().Len(); i++ {
		resource := td.ResourceSpans().At(i)
		resourceMaskFunc := p.createMaskFunc(resourceField)
		resource.Resource().Attributes().Range(resourceMaskFunc)
		for j := 0; j < resource.ScopeSpans().Len(); j++ {
			scope := resource.ScopeSpans().At(j)
			for k := 0; k < scope.Spans().Len(); k++ {
				spans := scope.Spans().At(k)
				attrMaskFunc := p.createMaskFunc(attributesField)
				spans.Attributes().Range(attrMaskFunc)
			}
		}
	}

	return td, nil
}

// processMetrics masks incoming metrics.
func (p *processor) processMetrics(_ context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		resource := md.ResourceMetrics().At(i)
		resourceMaskFunc := p.createMaskFunc(resourceField)
		resource.Resource().Attributes().Range(resourceMaskFunc)
		for j := 0; j < resource.ScopeMetrics().Len(); j++ {
			scope := resource.ScopeMetrics().At(j)
			for k := 0; k < scope.Metrics().Len(); k++ {
				metrics := scope.Metrics().At(k)
				switch metrics.Type() {
				case pmetric.MetricTypeSum:
					p.processSum(metrics.Sum())
				case pmetric.MetricTypeGauge:
					p.processGauge(metrics.Gauge())
				case pmetric.MetricTypeSummary:
					p.processSummary(metrics.Summary())
				}
			}
		}
	}

	return md, nil
}

// processSum masks a sum metric.
func (p *processor) processSum(sum pmetric.Sum) {
	for i := 0; i < sum.DataPoints().Len(); i++ {
		maskFunc := p.createMaskFunc(attributesField)
		sum.DataPoints().At(i).Attributes().Range(maskFunc)
	}
}

// processGauge masks a gauge metric.
func (p *processor) processGauge(gauge pmetric.Gauge) {
	for i := 0; i < gauge.DataPoints().Len(); i++ {
		maskFunc := p.createMaskFunc(attributesField)
		gauge.DataPoints().At(i).Attributes().Range(maskFunc)
	}
}

// processSummary masks a summary metrics.
func (p *processor) processSummary(summary pmetric.Summary) {
	for i := 0; i < summary.DataPoints().Len(); i++ {
		maskFunc := p.createMaskFunc(attributesField)
		summary.DataPoints().At(i).Attributes().Range(maskFunc)
	}
}

// maskValue masks a pcommon.Value.
func (p *processor) maskValue(field string, value pcommon.Value) {
	for _, excludeField := range p.cfg.Exclude {
		if strings.Contains(field, excludeField) {
			return
		}
	}

	switch value.Type() {
	case pcommon.ValueTypeMap:
		maskFunc := p.createMaskFunc(field)
		value.Map().Range(maskFunc)
	case pcommon.ValueTypeStr:
		p.maskString(value)
	}
}

// maskString masks a pcommon string.
func (p *processor) maskString(value pcommon.Value) {
	strValue := value.Str()

	for mask, rule := range p.rules {
		if !rule.MatchString(strValue) {
			continue
		}

		strValue = rule.ReplaceAllString(strValue, mask)
	}

	value.SetStr(strValue)
}

// createMaskFunc creates a func for ranging through a pcommon.Map and masking its values.
func (p *processor) createMaskFunc(field string) func(k string, v pcommon.Value) bool {
	return func(k string, v pcommon.Value) bool {
		childField := fmt.Sprintf("%s.%s", field, k)
		p.maskValue(childField, v)
		return true
	}
}
