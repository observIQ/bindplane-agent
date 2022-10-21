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

package samplingprocessor

import (
	"context"
	"math/rand"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type throughputMeasurementProcessor struct {
	logger          *zap.Logger
	dropCutOffRatio float64
}

func newThroughputMeasurementProcessor(logger *zap.Logger, cfg *Config) *throughputMeasurementProcessor {
	return &throughputMeasurementProcessor{
		logger:          logger,
		dropCutOffRatio: cfg.DropRatio,
	}
}

func (tmp *throughputMeasurementProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	//#nosec G404 -- randomly generated number is not used for security purposes. It's ok if it's weak
	if rand.Float64() <= tmp.dropCutOffRatio {
		return ptrace.NewTraces(), nil
	}

	return td, nil
}

func (tmp *throughputMeasurementProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	//#nosec G404 -- randomly generated number is not used for security purposes. It's ok if it's weak
	if rand.Float64() <= tmp.dropCutOffRatio {
		return plog.NewLogs(), nil
	}

	return ld, nil
}

func (tmp *throughputMeasurementProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	//#nosec G404 -- randomly generated number is not used for security purposes. It's ok if it's weak
	if rand.Float64() <= tmp.dropCutOffRatio {
		return pmetric.NewMetrics(), nil
	}

	return md, nil
}
