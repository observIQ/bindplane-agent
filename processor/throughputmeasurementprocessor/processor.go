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

package throughputmeasurementprocessor

import (
	"context"
	"crypto/rand"
	"math/big"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// upperBound is the upper bound to the sampling generator
var upperBound = big.NewInt(1000)

type throughputMeasurementProcessor struct {
	logger        *zap.Logger
	enabled       bool
	samplingRatio *big.Int
	mutators      []tag.Mutator
	tracesSizer   ptrace.Sizer
	metricsSizer  pmetric.Sizer
	logsSizer     plog.Sizer
}

func newThroughputMeasurementProcessor(logger *zap.Logger, cfg *Config, processorID string) *throughputMeasurementProcessor {
	return &throughputMeasurementProcessor{
		logger:        logger,
		enabled:       cfg.Enabled,
		samplingRatio: big.NewInt(int64(cfg.SamplingRatio * 1000)),
		mutators:      []tag.Mutator{tag.Upsert(processorTagKey, processorID, tag.WithTTL(tag.TTLNoPropagation))},
		tracesSizer:   ptrace.NewProtoMarshaler().(ptrace.Sizer),
		metricsSizer:  pmetric.NewProtoMarshaler().(pmetric.Sizer),
		logsSizer:     plog.NewProtoMarshaler().(plog.Sizer),
	}
}

func (tmp *throughputMeasurementProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	if tmp.enabled {
		i, err := rand.Int(rand.Reader, upperBound)
		if err != nil {
			return td, err
		}

		if i.Cmp(tmp.samplingRatio) <= 0 {
			err := stats.RecordWithTags(
				ctx,
				tmp.mutators,
				traceDataSize.M(int64(tmp.tracesSizer.TracesSize(td))),
			)

			if err != nil {
				return td, err
			}
		}
	}

	return td, nil
}

func (tmp *throughputMeasurementProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	if tmp.enabled {
		i, err := rand.Int(rand.Reader, upperBound)
		if err != nil {
			return ld, err
		}

		if i.Cmp(tmp.samplingRatio) <= 0 {
			err := stats.RecordWithTags(
				ctx,
				tmp.mutators,
				logDataSize.M(int64(tmp.logsSizer.LogsSize(ld))),
			)

			if err != nil {
				return ld, err
			}
		}
	}

	return ld, nil
}

func (tmp *throughputMeasurementProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	if tmp.enabled {
		i, err := rand.Int(rand.Reader, upperBound)
		if err != nil {
			return md, err
		}

		if i.Cmp(tmp.samplingRatio) <= 0 {
			err := stats.RecordWithTags(
				ctx,
				tmp.mutators,
				metricDataSize.M(int64(tmp.metricsSizer.MetricsSize(md))),
			)

			if err != nil {
				return md, err
			}
		}
	}

	return md, nil
}
