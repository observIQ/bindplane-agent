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

package googlemanagedprometheusexporter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func TestExporterCapabilities(t *testing.T) {
	exporter := &googleManagedPrometheusExporter{}
	capabilities := exporter.Capabilities()
	assert.True(t, capabilities.MutatesData)
}

func TestExporterWithConsumers(t *testing.T) {
	consumer := &MockProcessor{}
	consumer.On("ConsumeMetrics", mock.Anything, mock.Anything).Return(nil).Once()
	exporter := &googleManagedPrometheusExporter{
		metricsConsumer: consumer,
	}

	ctx := context.Background()

	err := exporter.ConsumeMetrics(ctx, pmetric.NewMetrics())
	assert.Nil(t, err)

	consumer.AssertExpectations(t)
}

func TestExporterWithoutConsumers(t *testing.T) {
	exporter := &googleManagedPrometheusExporter{}

	ctx := context.Background()

	err := exporter.ConsumeMetrics(ctx, pmetric.NewMetrics())
	assert.Nil(t, err)
}

func TestExporterStart(t *testing.T) {
	testCases := []struct {
		name          string
		exporter      *googleManagedPrometheusExporter
		expectedError error
	}{
		{
			name: "Successful metrics",
			exporter: &googleManagedPrometheusExporter{
				metricsProcessors: []processor.Metrics{createValidProcessor()},
				metricsExporter:   createValidExporter(),
			},
		},
		{
			name: "Failing metrics processor",
			exporter: &googleManagedPrometheusExporter{
				metricsProcessors: []processor.Metrics{
					createValidProcessor(),
					createFailingProcessor(),
				},
				metricsExporter: createValidExporter(),
			},
			expectedError: errors.New("failed to start metrics processor"),
		},
		{
			name: "Failing metrics exporter",
			exporter: &googleManagedPrometheusExporter{
				metricsProcessors: []processor.Metrics{
					createValidProcessor(),
					createValidProcessor(),
				},
				metricsExporter: createFailingExporter(),
			},
			expectedError: errors.New("failed to start metrics exporter"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.exporter.Start(context.Background(), nil)

			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	testCases := []struct {
		name          string
		exporter      *googleManagedPrometheusExporter
		expectedError error
	}{
		{
			name: "Successful metrics",
			exporter: &googleManagedPrometheusExporter{
				metricsProcessors: []processor.Metrics{createValidProcessor()},
				metricsExporter:   createValidExporter(),
			},
		},
		{
			name: "Failing metrics processor",
			exporter: &googleManagedPrometheusExporter{
				metricsProcessors: []processor.Metrics{
					createValidProcessor(),
					createFailingProcessor(),
				},
				metricsExporter: createValidExporter(),
			},
			expectedError: errors.New("failed to shutdown metrics processor"),
		},
		{
			name: "Failing metrics exporter",
			exporter: &googleManagedPrometheusExporter{
				metricsProcessors: []processor.Metrics{
					createValidProcessor(),
					createValidProcessor(),
				},
				metricsExporter: createFailingExporter(),
			},
			expectedError: errors.New("failed to shutdown metrics exporter"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.exporter.Shutdown(context.Background())

			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestAppendMetricAttrs(t *testing.T) {
	metrics := pmetric.NewMetrics()
	metric1 := metrics.ResourceMetrics().AppendEmpty()
	metric1.Resource().Attributes().PutStr(string(semconv.HostNameKey), "test-hostname")
	metric2 := metrics.ResourceMetrics().AppendEmpty()

	e := googleManagedPrometheusExporter{}
	e.appendMetricHost(&metrics)

	metric1Host, ok := metric1.Resource().Attributes().Get(string(semconv.HostNameKey))
	require.True(t, ok)
	require.Equal(t, "test-hostname", metric1Host.AsString())

	metric2Host, ok := metric2.Resource().Attributes().Get(string(semconv.HostNameKey))
	require.True(t, ok)
	require.Equal(t, hostname, metric2Host.AsString())
}

func createValidProcessor() *MockProcessor {
	processor := &MockProcessor{}
	processor.On("Start", mock.Anything, mock.Anything).Return(nil)
	processor.On("Shutdown", mock.Anything).Return(nil)
	return processor
}

func createFailingProcessor() *MockProcessor {
	processor := &MockProcessor{}
	processor.On("Start", mock.Anything, mock.Anything).Return(errors.New("failure"))
	processor.On("Shutdown", mock.Anything).Return(errors.New("failure"))
	return processor
}

func createValidExporter() *MockExporter {
	exporter := &MockExporter{}
	exporter.On("Start", mock.Anything, mock.Anything).Return(nil)
	exporter.On("Shutdown", mock.Anything).Return(nil)
	return exporter
}

func createFailingExporter() *MockExporter {
	exporter := &MockExporter{}
	exporter.On("Start", mock.Anything, mock.Anything).Return(errors.New("failure"))
	exporter.On("Shutdown", mock.Anything).Return(errors.New("failure"))
	return exporter
}

// MockProcessor is an autogenerated mock type for the Processor type
type MockProcessor struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *MockProcessor) Capabilities() consumer.Capabilities {
	ret := _m.Called()

	var r0 consumer.Capabilities
	if rf, ok := ret.Get(0).(func() consumer.Capabilities); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(consumer.Capabilities)
	}

	return r0
}

// ConsumeLogs provides a mock function with given fields: ctx, ld
func (_m *MockProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	ret := _m.Called(ctx, ld)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, plog.Logs) error); ok {
		r0 = rf(ctx, ld)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumeMetrics provides a mock function with given fields: ctx, md
func (_m *MockProcessor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	ret := _m.Called(ctx, md)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pmetric.Metrics) error); ok {
		r0 = rf(ctx, md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumeTraces provides a mock function with given fields: ctx, td
func (_m *MockProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	ret := _m.Called(ctx, td)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ptrace.Traces) error); ok {
		r0 = rf(ctx, td)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields: ctx
func (_m *MockProcessor) Shutdown(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: ctx, host
func (_m *MockProcessor) Start(ctx context.Context, host component.Host) error {
	ret := _m.Called(ctx, host)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, component.Host) error); ok {
		r0 = rf(ctx, host)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockExporter is an autogenerated mock type for the Exporter type
type MockExporter struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *MockExporter) Capabilities() consumer.Capabilities {
	ret := _m.Called()

	var r0 consumer.Capabilities
	if rf, ok := ret.Get(0).(func() consumer.Capabilities); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(consumer.Capabilities)
	}

	return r0
}

// ConsumeMetrics provides a mock function with given fields: ctx, md
func (_m *MockExporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	ret := _m.Called(ctx, md)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pmetric.Metrics) error); ok {
		r0 = rf(ctx, md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields: ctx
func (_m *MockExporter) Shutdown(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: ctx, host
func (_m *MockExporter) Start(ctx context.Context, host component.Host) error {
	ret := _m.Called(ctx, host)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, component.Host) error); ok {
		r0 = rf(ctx, host)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
