// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package observiq

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/snappy"
	"github.com/observiq/bindplane-agent/internal/measurements"
	"github.com/observiq/bindplane-agent/opamp/mocks"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
)

var encodedSizeMetric = []byte{0x6f, 0xf0, 0x6e, 0xa, 0x6d, 0xa, 0x0, 0x12, 0x69, 0xa, 0x0, 0x12, 0x0, 0x12, 0x61, 0xa, 0x38, 0x6f, 0x74, 0x65, 0x6c, 0x63, 0x6f, 0x6c, 0x5f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x5f, 0x74, 0x68, 0x72, 0x6f, 0x75, 0x67, 0x68, 0x70, 0x75, 0x74, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x3a, 0x25, 0xa, 0x23, 0x19, 0x18, 0x59, 0x22, 0xbc, 0x5c, 0x12, 0xde, 0x17, 0x31, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3a, 0xf, 0xa, 0x9, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x2, 0xa, 0x0, 0x12, 0x0}

func TestMeasurementsSender(t *testing.T) {
	t.Run("Test emits metrics", func(t *testing.T) {
		dataChan := make(chan []byte, 1)
		client := mocks.NewMockOpAMPClient(t)
		client.On("SendCustomMessage", mock.Anything).Run(func(args mock.Arguments) {
			cm := args.Get(0).(*protobufs.CustomMessage)
			dataChan <- cm.Data
		}).Return(make(chan struct{}), nil)

		mp := metric.NewMeterProvider()
		defer mp.Shutdown(context.Background())

		processorID := "throughputmeasurement/1"

		tm, err := measurements.NewThroughputMeasurements(mp, "throughputmeasurement", processorID, map[string]string{})
		require.NoError(t, err)

		m, err := golden.ReadMetrics(filepath.Join("testdata", "metrics", "host-metrics.yaml"))
		require.NoError(t, err)

		tm.AddMetrics(context.Background(), m)

		reg := measurements.NewResettableThroughputMeasurementsRegistry(false)
		reg.RegisterThroughputMeasurements(processorID, tm)

		ms := newMeasurementsSender(zap.NewNop(), reg, client, 1*time.Millisecond)
		ms.Start()

		select {
		case <-time.After(1 * time.Second):
			require.FailNow(t, "timed out waiting for metrics payload")
		case d := <-dataChan:
			decoded, err := snappy.Decode(nil, d)
			require.NoError(t, err)

			um := &pmetric.ProtoUnmarshaler{}
			actualMetrics, err := um.UnmarshalMetrics(decoded)
			require.NoError(t, err)

			expectedMetrics, err := golden.ReadMetrics(filepath.Join("testdata", "metrics", "expected-throughput.yaml"))
			require.NoError(t, err)

			require.NoError(t, pmetrictest.CompareMetrics(expectedMetrics, actualMetrics, pmetrictest.IgnoreTimestamp()))
		}

		ms.Stop()
	})

	t.Run("Test set interval", func(t *testing.T) {
		dataChan := make(chan []byte, 1)
		client := mocks.NewMockOpAMPClient(t)
		client.On("SendCustomMessage", mock.Anything).Run(func(args mock.Arguments) {
			cm := args.Get(0).(*protobufs.CustomMessage)
			dataChan <- cm.Data
		}).Return(make(chan struct{}), nil)

		mp := metric.NewMeterProvider()
		defer mp.Shutdown(context.Background())

		processorID := "throughputmeasurement/1"

		tm, err := measurements.NewThroughputMeasurements(mp, "throughputmeasurement", processorID, map[string]string{})
		require.NoError(t, err)

		m, err := golden.ReadMetrics(filepath.Join("testdata", "metrics", "host-metrics.yaml"))
		require.NoError(t, err)

		tm.AddMetrics(context.Background(), m)

		reg := measurements.NewResettableThroughputMeasurementsRegistry(false)
		reg.RegisterThroughputMeasurements(processorID, tm)

		ms := newMeasurementsSender(zap.NewNop(), reg, client, 5*time.Hour)
		ms.Start()

		// Wait 200 ms and ensure no data emitted
		time.Sleep(200 * time.Millisecond)

		require.Len(t, dataChan, 0)

		// Set time to 1ms. We should see data emit quickly after.
		ms.SetInterval(1 * time.Millisecond)

		select {
		case <-time.After(1 * time.Second):
			require.FailNow(t, "timed out waiting for metrics payload")
		case d := <-dataChan:
			decoded, err := snappy.Decode(nil, d)
			require.NoError(t, err)

			um := &pmetric.ProtoUnmarshaler{}
			actualMetrics, err := um.UnmarshalMetrics(decoded)
			require.NoError(t, err)

			expectedMetrics, err := golden.ReadMetrics(filepath.Join("testdata", "metrics", "expected-throughput.yaml"))
			require.NoError(t, err)

			require.NoError(t, pmetrictest.CompareMetrics(expectedMetrics, actualMetrics, pmetrictest.IgnoreTimestamp()))
		}

		ms.Stop()
	})
}
