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

package throughputwrapper

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

var _ consumer.Metrics = (*metricConsumer)(nil)

type metricConsumer struct {
	logger       *zap.Logger
	mutators     []tag.Mutator
	metricsSizer pmetric.MarshalSizer
	baseConsumer consumer.Metrics
}

func newMetricConsumer(logger *zap.Logger, componentID string, baseConsumer consumer.Metrics) *metricConsumer {
	return &metricConsumer{
		logger:       logger,
		mutators:     []tag.Mutator{tag.Upsert(componentTagKey, componentID, tag.WithTTL(tag.TTLNoPropagation))},
		metricsSizer: pmetric.NewProtoMarshaler(),
		baseConsumer: baseConsumer,
	}
}

// ConsumeMetrics measures the pmetric.Metrics size before passing it onto the baseConsumer
func (m *metricConsumer) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	if err := stats.RecordWithTags(
		ctx,
		m.mutators,
		metricThroughputSize.M(int64(m.metricsSizer.MetricsSize(md))),
	); err != nil {
		m.logger.Warn("Error while measuring receiver metric throughput", zap.Error(err))
	}

	return m.baseConsumer.ConsumeMetrics(ctx, md)
}

// Capabilities returns the baseConsumer's capabilities
func (m *metricConsumer) Capabilities() consumer.Capabilities {
	return m.baseConsumer.Capabilities()
}
