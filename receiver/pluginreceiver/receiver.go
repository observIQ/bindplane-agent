package pluginreceiver

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

// Receiver is a receiver that runs an embedded open telemetry config
// as an internal service.
type Receiver struct {
	plugin          Plugin
	configProvider  *ConfigProvider
	factoryProvider *FactoryProvider
	buildInfo       component.BuildInfo
	logger          *zap.Logger
	svc             Service
	createSvc       func(set service.CollectorSettings) (Service, error)
}

// Start starts the receiver's internal service
func (r *Receiver) Start(ctx context.Context, host component.Host) error {
	r.logger.Info("starting plugin", zap.String("plugin", r.plugin.Title), zap.String("plugin-version", r.plugin.Version))

	factories, err := r.factoryProvider.GetFactories(host, r.configProvider.configMap)
	if err != nil {
		return fmt.Errorf("failed to get factories: %w", err)
	}

	settings := service.CollectorSettings{
		Factories:               *factories,
		BuildInfo:               r.buildInfo,
		DisableGracefulShutdown: true,
		ConfigProvider:          r.configProvider,
	}

	svc, err := r.createSvc(settings)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	r.svc = svc

	if err := startService(ctx, svc); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
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

// Shutdown stops the receiver's internal service
func (r *Receiver) Shutdown(ctx context.Context) error {
	if r.svc != nil {
		r.svc.Shutdown()
	}

	return nil
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
