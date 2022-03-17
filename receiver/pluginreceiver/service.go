package pluginreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Service is an interface for running an internal open telemetry service
type Service interface {
	Run(ctx context.Context) error
	Shutdown()
	GetState() service.State
}

// createService creates a default Service for running an open telemetry pipeline
func createService(factories component.Factories, configProvider service.ConfigProvider, logger *zap.Logger) (Service, error) {
	settings := service.CollectorSettings{
		Factories:               factories,
		DisableGracefulShutdown: true,
		ConfigProvider:          configProvider,
		LoggingOptions:          createServiceLoggerOpts(logger),
	}

	return service.New(settings)
}

// createServiceLoggerOpts creates the default logger opts for a Service
func createServiceLoggerOpts(baseLogger *zap.Logger) []zap.Option {
	levelOpt := zap.IncreaseLevel(zap.ErrorLevel)
	coreOpt := zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return baseLogger.Core()
	})
	return []zap.Option{coreOpt, levelOpt}
}

// createServiceFunc is a function used to create a service
type createServiceFunc = func(factories component.Factories, configProvider service.ConfigProvider, logger *zap.Logger) (Service, error)
