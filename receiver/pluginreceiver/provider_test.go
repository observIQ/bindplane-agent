package pluginreceiver

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/extension/extensionhelper"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

func TestConfigProviderWatch(t *testing.T) {
	provider := createConfigProvider(&config.Map{})
	err := errors.New("config err")
	go func() {
		provider.errChan <- err
	}()

	receivedErr := <-provider.Watch()
	require.Equal(t, err, receivedErr)
}

func TestConfigProviderShutdown(t *testing.T) {
	provider := createConfigProvider(&config.Map{})
	err := provider.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestConfigProviderGet(t *testing.T) {
	ctx := context.Background()
	configMap := &config.Map{}
	provider := createConfigProvider(configMap)
	factories := component.Factories{}

	unmarshaller := &MockUnmarshaller{}
	unmarshaller.On("Unmarshal", configMap, mock.Anything).Return(nil, errors.New("failure")).Once()
	unmarshaller.On("Unmarshal", configMap, mock.Anything).Return(&config.Config{}, nil).Once()
	provider.unmarshaller = unmarshaller

	cfg, err := provider.Get(ctx, factories)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal config")
	require.Nil(t, cfg)

	cfg, err = provider.Get(ctx, factories)
	require.NoError(t, err)
	require.Equal(t, &config.Config{}, cfg)
}

func TestGetFactories(t *testing.T) {
	testType := config.Type("test")
	testReceiverFactory := receiverhelper.NewFactory(testType, nil)
	testProcessorFactory := processorhelper.NewFactory(testType, nil)
	testExporterFactory := exporterhelper.NewFactory(testType, nil)
	testExtensionFactory := extensionhelper.NewFactory(testType, nil, nil)

	testCases := []struct {
		name              string
		config            map[string]interface{}
		providerFactories component.Factories
		hostFactories     component.Factories
		expectedResult    *component.Factories
		expectedErr       error
	}{
		{
			name: "invalid config",
			config: map[string]interface{}{
				"receivers": 5,
			},
			expectedErr: errors.New("failed to unmarshal config"),
		},
		{
			name: "missing receiver factory",
			config: map[string]interface{}{
				"receivers": map[string]interface{}{
					"test": nil,
				},
			},
			expectedErr: errors.New("failed to get receiver factories"),
		},
		{
			name: "receiver factory exists on provider",
			config: map[string]interface{}{
				"receivers": map[string]interface{}{
					"test": nil,
				},
			},
			providerFactories: component.Factories{
				Receivers: map[config.Type]component.ReceiverFactory{
					testType: testReceiverFactory,
				},
			},
			expectedResult: &component.Factories{
				Receivers: map[config.Type]component.ReceiverFactory{
					testType: testReceiverFactory,
				},
				Processors: map[config.Type]component.ProcessorFactory{},
				Extensions: map[config.Type]component.ExtensionFactory{},
				Exporters:  map[config.Type]component.ExporterFactory{},
			},
		},
		{
			name: "receiver factory exists on host",
			config: map[string]interface{}{
				"receivers": map[string]interface{}{
					"test": nil,
				},
			},
			hostFactories: component.Factories{
				Receivers: map[config.Type]component.ReceiverFactory{
					testType: testReceiverFactory,
				},
			},
			expectedResult: &component.Factories{
				Receivers: map[config.Type]component.ReceiverFactory{
					testType: testReceiverFactory,
				},
				Processors: map[config.Type]component.ProcessorFactory{},
				Extensions: map[config.Type]component.ExtensionFactory{},
				Exporters:  map[config.Type]component.ExporterFactory{},
			},
		},
		{
			name: "missing processor factory",
			config: map[string]interface{}{
				"processors": map[string]interface{}{
					"test": nil,
				},
			},
			expectedErr: errors.New("failed to get processor factories"),
		},
		{
			name: "processor factory exists on provider",
			config: map[string]interface{}{
				"processors": map[string]interface{}{
					"test": nil,
				},
			},
			providerFactories: component.Factories{
				Processors: map[config.Type]component.ProcessorFactory{
					testType: testProcessorFactory,
				},
			},
			expectedResult: &component.Factories{
				Processors: map[config.Type]component.ProcessorFactory{
					testType: testProcessorFactory,
				},
				Receivers:  map[config.Type]component.ReceiverFactory{},
				Extensions: map[config.Type]component.ExtensionFactory{},
				Exporters:  map[config.Type]component.ExporterFactory{},
			},
		},
		{
			name: "processor factory exists on host",
			config: map[string]interface{}{
				"processors": map[string]interface{}{
					"test": nil,
				},
			},
			hostFactories: component.Factories{
				Processors: map[config.Type]component.ProcessorFactory{
					testType: testProcessorFactory,
				},
			},
			expectedResult: &component.Factories{
				Processors: map[config.Type]component.ProcessorFactory{
					testType: testProcessorFactory,
				},
				Receivers:  map[config.Type]component.ReceiverFactory{},
				Extensions: map[config.Type]component.ExtensionFactory{},
				Exporters:  map[config.Type]component.ExporterFactory{},
			},
		},
		{
			name: "missing exporter factory",
			config: map[string]interface{}{
				"exporters": map[string]interface{}{
					"test": nil,
				},
			},
			expectedErr: errors.New("failed to get exporter factories"),
		},
		{
			name: "exporter factory exists on provider",
			config: map[string]interface{}{
				"exporters": map[string]interface{}{
					"test": nil,
				},
			},
			providerFactories: component.Factories{
				Exporters: map[config.Type]component.ExporterFactory{
					testType: testExporterFactory,
				},
			},
			expectedResult: &component.Factories{
				Exporters: map[config.Type]component.ExporterFactory{
					testType: testExporterFactory,
				},
				Processors: map[config.Type]component.ProcessorFactory{},
				Extensions: map[config.Type]component.ExtensionFactory{},
				Receivers:  map[config.Type]component.ReceiverFactory{},
			},
		},
		{
			name: "exporter factory exists on host",
			config: map[string]interface{}{
				"exporters": map[string]interface{}{
					"test": nil,
				},
			},
			hostFactories: component.Factories{
				Exporters: map[config.Type]component.ExporterFactory{
					testType: testExporterFactory,
				},
			},
			expectedResult: &component.Factories{
				Exporters: map[config.Type]component.ExporterFactory{
					testType: testExporterFactory,
				},
				Processors: map[config.Type]component.ProcessorFactory{},
				Extensions: map[config.Type]component.ExtensionFactory{},
				Receivers:  map[config.Type]component.ReceiverFactory{},
			},
		},
		{
			name: "missing extension factory",
			config: map[string]interface{}{
				"extensions": map[string]interface{}{
					"test": nil,
				},
			},
			expectedErr: errors.New("failed to get extension factories"),
		},
		{
			name: "extension factory exists on provider",
			config: map[string]interface{}{
				"extensions": map[string]interface{}{
					"test": nil,
				},
			},
			providerFactories: component.Factories{
				Extensions: map[config.Type]component.ExtensionFactory{
					testType: testExtensionFactory,
				},
			},
			expectedResult: &component.Factories{
				Extensions: map[config.Type]component.ExtensionFactory{
					testType: testExtensionFactory,
				},
				Processors: map[config.Type]component.ProcessorFactory{},
				Receivers:  map[config.Type]component.ReceiverFactory{},
				Exporters:  map[config.Type]component.ExporterFactory{},
			},
		},
		{
			name: "extension factory exists on host",
			config: map[string]interface{}{
				"extensions": map[string]interface{}{
					"test": nil,
				},
			},
			hostFactories: component.Factories{
				Extensions: map[config.Type]component.ExtensionFactory{
					testType: testExtensionFactory,
				},
			},
			expectedResult: &component.Factories{
				Extensions: map[config.Type]component.ExtensionFactory{
					testType: testExtensionFactory,
				},
				Processors: map[config.Type]component.ProcessorFactory{},
				Receivers:  map[config.Type]component.ReceiverFactory{},
				Exporters:  map[config.Type]component.ExporterFactory{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			host := &MockHost{}
			for key, factory := range tc.hostFactories.Receivers {
				host.On("GetFactory", component.KindReceiver, key).Return(factory)
			}

			for key, factory := range tc.hostFactories.Processors {
				host.On("GetFactory", component.KindProcessor, key).Return(factory)
			}

			for key, factory := range tc.hostFactories.Exporters {
				host.On("GetFactory", component.KindExporter, key).Return(factory)
			}

			for key, factory := range tc.hostFactories.Extensions {
				host.On("GetFactory", component.KindExtension, key).Return(factory)
			}
			host.On("GetFactory", mock.Anything, mock.Anything).Return(nil)
			configMap := config.NewMapFromStringMap(tc.config)
			provider := FactoryProvider{factories: tc.providerFactories}

			factories, err := provider.GetFactories(host, configMap)
			switch tc.expectedErr {
			case nil:
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, factories)
			default:
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

// MockUnmarshaller is a mock type for the configunmarshaler.Unmarshaller type
type MockUnmarshaller struct {
	mock.Mock
}

// Unmarshal provides a mock function with given fields: v, factories
func (_m *MockUnmarshaller) Unmarshal(v *config.Map, factories component.Factories) (*config.Config, error) {
	ret := _m.Called(v, factories)

	var r0 *config.Config
	if rf, ok := ret.Get(0).(func(*config.Map, component.Factories) *config.Config); ok {
		r0 = rf(v, factories)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Config)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*config.Map, component.Factories) error); ok {
		r1 = rf(v, factories)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHost is a mock type for the component.Host type
type MockHost struct {
	mock.Mock
}

// GetExporters provides a mock function with given fields:
func (_m *MockHost) GetExporters() map[config.Type]map[config.ComponentID]component.Exporter {
	ret := _m.Called()

	var r0 map[config.Type]map[config.ComponentID]component.Exporter
	if rf, ok := ret.Get(0).(func() map[config.Type]map[config.ComponentID]component.Exporter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[config.Type]map[config.ComponentID]component.Exporter)
		}
	}

	return r0
}

// GetExtensions provides a mock function with given fields:
func (_m *MockHost) GetExtensions() map[config.ComponentID]component.Extension {
	ret := _m.Called()

	var r0 map[config.ComponentID]component.Extension
	if rf, ok := ret.Get(0).(func() map[config.ComponentID]component.Extension); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[config.ComponentID]component.Extension)
		}
	}

	return r0
}

// GetFactory provides a mock function with given fields: kind, componentType
func (_m *MockHost) GetFactory(kind component.Kind, componentType config.Type) component.Factory {
	ret := _m.Called(kind, componentType)

	var r0 component.Factory
	if rf, ok := ret.Get(0).(func(component.Kind, config.Type) component.Factory); ok {
		r0 = rf(kind, componentType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(component.Factory)
		}
	}

	return r0
}

// ReportFatalError provides a mock function with given fields: err
func (_m *MockHost) ReportFatalError(err error) {
	_m.Called(err)
}
