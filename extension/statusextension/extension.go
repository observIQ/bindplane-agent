package statusextension

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

type statusExtension struct {
	logger *zap.Logger
}

func newStatusExtension(params component.TelemetrySettings) *statusExtension {
	return &statusExtension{
		logger: params.Logger.Named("statusExtension"),
	}
}

func (s *statusExtension) Start(_ context.Context, _ component.Host) error {
	return nil
}

func (s *statusExtension) Shutdown(_ context.Context) error {
	return nil
}

func (s *statusExtension) ComponentStatusChanged(source *component.InstanceID, event *component.StatusEvent) {
	s.logger.Info("Status Changed",
		zap.String("sourceID", source.ID.String()),
		zap.String("event_status", event.Status().String()),
		zap.String("event_timestamp", event.Timestamp().String()),
		zap.Error(event.Err()),
	)
}
