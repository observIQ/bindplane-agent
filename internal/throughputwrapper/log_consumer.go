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
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

var _ consumer.Logs = (*logConsumer)(nil)

type logConsumer struct {
	logger       *zap.Logger
	mutators     []tag.Mutator
	logsSizer    plog.MarshalSizer
	baseConsumer consumer.Logs
}

func newLogConsumer(logger *zap.Logger, componentID string, baseConsumer consumer.Logs) *logConsumer {
	return &logConsumer{
		logger:       logger,
		mutators:     []tag.Mutator{tag.Upsert(componentTagKey, componentID, tag.WithTTL(tag.TTLNoPropagation))},
		logsSizer:    plog.NewProtoMarshaler(),
		baseConsumer: baseConsumer,
	}
}

func (l *logConsumer) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	if err := stats.RecordWithTags(
		ctx,
		l.mutators,
		logThroughputSize.M(int64(l.logsSizer.LogsSize(ld))),
	); err != nil {
		l.logger.Warn("Error while measuring receiver log throughput", zap.Error(err))
	}
	return l.baseConsumer.ConsumeLogs(ctx, ld)
}

func (l *logConsumer) Capabilities() consumer.Capabilities {
	return l.baseConsumer.Capabilities()
}
