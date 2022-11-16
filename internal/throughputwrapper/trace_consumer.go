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
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

var _ consumer.Traces = (*traceConsumer)(nil)

type traceConsumer struct {
	logger       *zap.Logger
	mutators     []tag.Mutator
	tracesSizer  ptrace.MarshalSizer
	baseConsumer consumer.Traces
}

func newTraceConsumer(logger *zap.Logger, componentID string, baseConsumer consumer.Traces) *traceConsumer {
	return &traceConsumer{
		logger:       logger,
		mutators:     []tag.Mutator{tag.Upsert(componentTagKey, componentID, tag.WithTTL(tag.TTLNoPropagation))},
		tracesSizer:  &ptrace.ProtoMarshaler{},
		baseConsumer: baseConsumer,
	}
}

// ConsumeTraces measures the ptrace.Traces size before passing it onto the baseConsumer
func (t *traceConsumer) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	if err := stats.RecordWithTags(
		ctx,
		t.mutators,
		traceThroughputSize.M(int64(t.tracesSizer.TracesSize(td))),
	); err != nil {
		t.logger.Warn("Error while measuring receiver trace throughput", zap.Error(err))
	}
	return t.baseConsumer.ConsumeTraces(ctx, td)
}

// Capabilities returns the baseConsumer's capabilities
func (t *traceConsumer) Capabilities() consumer.Capabilities {
	return t.baseConsumer.Capabilities()
}
