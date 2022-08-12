package throughputmeasurementprocessor

import (
	"context"
	"crypto/rand"
	"math/big"

	"go.opencensus.io/stats"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

var upperBound = big.NewInt(1000)

type throughputMeasurementProcessor struct {
	logger        *zap.Logger
	enabled       bool
	samplingRatio *big.Int
}

func newThroughputMeasurementProcessor(logger *zap.Logger, cfg *Config) *throughputMeasurementProcessor {
	return &throughputMeasurementProcessor{
		logger:        logger,
		enabled:       cfg.Enabled,
		samplingRatio: big.NewInt(int64(cfg.SamplingRatio * 1000)),
	}
}

func (tmp *throughputMeasurementProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	if tmp.enabled {
		i, err := rand.Int(rand.Reader, upperBound)
		if err != nil {
			return td, err
		}

		if i.Cmp(tmp.samplingRatio) <= 0 {
			stats.Record(ctx, traceDataSize.M(int64(td.Size())))
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
			stats.Record(ctx, logDataSize.M(int64(ld.Size())))
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
			stats.Record(ctx, metricDataSize.M(int64(md.Size())))
		}
	}

	return md, nil
}
