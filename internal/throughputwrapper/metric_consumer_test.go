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
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func Test_newMetricConsumer(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := consumertest.NewNop()
	mConsumer := newMetricConsumer(nopLogger, componentID, baseConsumer)

	require.Equal(t, nopLogger, mConsumer.logger)
	require.Equal(t, baseConsumer, mConsumer.baseConsumer)
	require.Len(t, mConsumer.mutators, 1)
	require.Equal(t, &pmetric.ProtoMarshaler{}, mConsumer.metricsSizer)
}

func Test_metricConsumer_ConsumeMetrics(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := new(consumertest.MetricsSink)
	mConsumer := newMetricConsumer(nopLogger, componentID, baseConsumer)

	md := pmetric.NewMetrics()
	metric := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
	metric.SetEmptyGauge()
	metric.Gauge().DataPoints().AppendEmpty()

	err := mConsumer.ConsumeMetrics(context.Background(), md)
	require.NoError(t, err)

	require.Equal(t, 1, baseConsumer.DataPointCount())
}

func Test_metricConsumer_Capabilities(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := consumertest.NewNop()
	mConsumer := newMetricConsumer(nopLogger, componentID, baseConsumer)

	require.Equal(t, baseConsumer.Capabilities(), mConsumer.Capabilities())
}
