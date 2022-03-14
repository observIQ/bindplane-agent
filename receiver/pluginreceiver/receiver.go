package pluginreceiver

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Receiver is a receiver that runs an embedded open telemetry config
// as an internal service.
type Receiver struct {
	plugin          *Plugin
	configProvider  *ConfigProvider
	factoryProvider *FactoryProvider
	buildInfo       component.BuildInfo
	logger          *zap.Logger
	service         Service
	createService   func(set service.CollectorSettings) (Service, error)
}

// Start starts the receiver's internal service
func (r *Receiver) Start(ctx context.Context, host component.Host) error {
	r.logger.Info("Starting plugin...", zap.String("plugin", r.plugin.Title), zap.String("plugin-version", r.plugin.Version))

	settings, err := r.createServiceSettings(host)
	if err != nil {
		return fmt.Errorf("failed to create internal service settings: %w", err)
	}

	service, err := r.createService(*settings)
	if err != nil {
		return fmt.Errorf("failed to create internal service: %w", err)
	}
	r.service = service

	if err := startService(ctx, service); err != nil {
		return fmt.Errorf("failed to start internal service: %w", err)
	}
	r.logger.Info("Started plugin")

	return nil
}

// Shutdown stops the receiver's internal service
func (r *Receiver) Shutdown(_ context.Context) error {
	if r.service != nil {
		r.service.Shutdown()
	}

	return nil
}

// createServiceSettings creates the settings for the internal service
func (r *Receiver) createServiceSettings(host component.Host) (*service.CollectorSettings, error) {
	factories, err := r.factoryProvider.GetFactories(host, r.configProvider.configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get factories from factory provider: %w", err)
	}

	coreOpt := wrapLogger(r.logger)
	levelOpt := zap.IncreaseLevel(zap.ErrorLevel)
	loggingOpts := []zap.Option{coreOpt, levelOpt}

	return &service.CollectorSettings{
		Factories:               *factories,
		BuildInfo:               r.buildInfo,
		DisableGracefulShutdown: true,
		ConfigProvider:          r.configProvider,
		LoggingOptions:          loggingOpts,
	}, nil
}

// wrapLogger wraps a logger's core
func wrapLogger(logger *zap.Logger) zap.Option {
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return logger.Core()
	})
}

// startService starts the provided service
func startService(ctx context.Context, svc Service) error {
	errChan := make(chan error)
	go func() {
		if err := svc.Run(ctx); err != nil {
			errChan <- err
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 250)
	defer ticker.Stop()

	for {
		select {
		case err := <-errChan:
			return err
		case <-ticker.C:
			if svc.GetState() == service.Running {
				return nil
			}
		}
	}
}

// Service is an interface for running an internal open telemetry config
type Service interface {
	Run(ctx context.Context) error
	Shutdown()
	GetState() service.State
}

// createDefaultService creates the default service for running an internal open telemetry config
func createDefaultService(set service.CollectorSettings) (Service, error) {
	return service.New(set)
}
