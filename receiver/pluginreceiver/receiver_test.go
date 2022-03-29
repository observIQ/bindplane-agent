package pluginreceiver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

func TestReceiverGetFactoryFailure(t *testing.T) {
	ctx := context.Background()
	host := &MockHost{}
	host.On("GetFactory", mock.Anything, mock.Anything).Return(nil)

	components := ComponentMap{
		Receivers: map[config.ComponentID]map[string]interface{}{
			config.NewComponentID("missing"): nil,
		},
	}
	configProvider := createConfigProvider(&components)
	emitterFactory := createLogEmitterFactory(nil)

	receiver := Receiver{
		plugin:         &Plugin{},
		configProvider: configProvider,
		emitterFactory: emitterFactory,
		logger:         zap.NewNop(),
	}

	err := receiver.Start(ctx, host)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get factories")
}

func TestReceiverCreateServiceFailure(t *testing.T) {
	nopFactory := receiverhelper.NewFactory("nop", nil)
	ctx := context.Background()
	host := &MockHost{}
	host.On("GetFactory", mock.Anything, mock.Anything).Return(nopFactory)

	components := ComponentMap{
		Receivers: map[config.ComponentID]map[string]interface{}{
			config.NewComponentID("nop"): nil,
		},
	}
	configProvider := createConfigProvider(&components)
	emitterFactory := createLogEmitterFactory(nil)

	receiver := Receiver{
		plugin:         &Plugin{},
		configProvider: configProvider,
		emitterFactory: emitterFactory,
		logger:         zap.NewNop(),
		createService: func(factories component.Factories, configProvider service.ConfigProvider, logger *zap.Logger) (Service, error) {
			return nil, errors.New("failure")
		},
	}

	err := receiver.Start(ctx, host)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create internal service")
}

func TestReceiverStartServiceFailure(t *testing.T) {
	nopFactory := receiverhelper.NewFactory("nop", nil)
	ctx := context.Background()
	host := &MockHost{}
	host.On("GetFactory", mock.Anything, mock.Anything).Return(nopFactory)

	components := ComponentMap{
		Receivers: map[config.ComponentID]map[string]interface{}{
			config.NewComponentID("nop"): nil,
		},
	}
	configProvider := createConfigProvider(&components)
	emitterFactory := createLogEmitterFactory(nil)

	svc := &MockService{}
	svc.On("Run", mock.Anything).Return(errors.New("failure"))
	svc.On("GetState").Return(service.Starting)

	receiver := Receiver{
		plugin:         &Plugin{},
		configProvider: configProvider,
		emitterFactory: emitterFactory,
		logger:         zap.NewNop(),
		createService: func(factories component.Factories, configProvider service.ConfigProvider, logger *zap.Logger) (Service, error) {
			return svc, nil
		},
	}

	err := receiver.Start(ctx, host)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to start internal service")
}

func TestReceiverStartSuccess(t *testing.T) {
	nopFactory := receiverhelper.NewFactory("nop", nil)
	ctx := context.Background()
	host := &MockHost{}
	host.On("GetFactory", mock.Anything, mock.Anything).Return(nopFactory)

	components := ComponentMap{
		Receivers: map[config.ComponentID]map[string]interface{}{
			config.NewComponentID("nop"): nil,
		},
	}
	configProvider := createConfigProvider(&components)
	emitterFactory := createLogEmitterFactory(nil)

	svc := &MockService{}
	svc.On("Run", mock.Anything).WaitUntil(time.After(time.Second)).Return(errors.New("unexpected timeout"))
	svc.On("GetState").Return(service.Running)

	receiver := Receiver{
		plugin:         &Plugin{},
		configProvider: configProvider,
		emitterFactory: emitterFactory,
		logger:         zap.NewNop(),
		createService: func(factories component.Factories, configProvider service.ConfigProvider, logger *zap.Logger) (Service, error) {
			return svc, nil
		},
	}

	err := receiver.Start(ctx, host)
	require.NoError(t, err)
}

func TestReceiverShutdown(t *testing.T) {
	ctx := context.Background()
	receiver := Receiver{}
	err := receiver.Shutdown(ctx)
	require.NoError(t, err)

	service := &MockService{}
	service.On("Shutdown").Return().Once()
	receiver.service = service
	err = receiver.Shutdown(ctx)
	require.NoError(t, err)
	service.AssertExpectations(t)
}

// MockService is a mock type for the Service type
type MockService struct {
	mock.Mock
}

// GetState provides a mock function with given fields:
func (_m *MockService) GetState() service.State {
	ret := _m.Called()

	var r0 service.State
	if rf, ok := ret.Get(0).(func() service.State); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(service.State)
	}

	return r0
}

// Run provides a mock function with given fields: ctx
func (_m *MockService) Run(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields:
func (_m *MockService) Shutdown() {
	_m.Called()
}
