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
	plugin         *Plugin
	configProvider *ConfigProvider
	emitterFactory component.ExporterFactory
	logger         *zap.Logger
	createService  createServiceFunc
	service        Service
}

// Start starts the receiver's internal service
func (r *Receiver) Start(ctx context.Context, host component.Host) error {
	r.logger.Info("Starting plugin...", zap.String("plugin", r.plugin.Title), zap.String("plugin-version", r.plugin.Version))

	factories, err := r.configProvider.GetRequiredFactories(host, r.emitterFactory)
	if err != nil {
		return fmt.Errorf("failed to get factories from factory provider: %w", err)
	}

	service, err := r.createService(*factories, r.configProvider, r.logger)
	if err != nil {
		return fmt.Errorf("failed to create internal service: %w", err)
	}
	r.service = service

	if err := startService(ctx, service); err != nil {
		return fmt.Errorf("failed to start internal service: %w", err)
	}

	return nil
}

// Shutdown stops the receiver's internal service
func (r *Receiver) Shutdown(_ context.Context) error {
	if r.service != nil {
		r.service.Shutdown()
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
