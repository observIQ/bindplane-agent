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
			stats.RecordWithTags(
				ctx,
				tmp.mutators,
				traceDataSize.M(int64(tmp.tracesSizer.TracesSize(td))),
			)
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
			stats.RecordWithTags(
				ctx,
				tmp.mutators,
				logDataSize.M(int64(tmp.logsSizer.LogsSize(ld))),
			)
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
			stats.RecordWithTags(
				ctx,
				tmp.mutators,
				metricDataSize.M(int64(tmp.metricsSizer.MetricsSize(md))),
			)
		}
	}

	return md, nil
}
