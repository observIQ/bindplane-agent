package generatereceiver

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

// Receiver is a receiver that generates telemetry.
type Receiver struct {
	config         *Config
	logConsumer    consumer.Logs
	metricConsumer consumer.Metrics
	traceConsumer  consumer.Traces
	logger         *zap.Logger
	cancel         context.CancelFunc
}

// Start starts the receiver.
func (r *Receiver) Start(_ context.Context, _ component.Host) error {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	if r.logConsumer != nil {
		go r.GenerateLogs(ctx)
	}

	return nil
}

// Shutdown stops the receiver.
func (r *Receiver) Shutdown(_ context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}

	return nil
}

// GenerateLogs generates logs.
func (r *Receiver) GenerateLogs(ctx context.Context) {
	ticker := time.NewTicker(r.config.Logs.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logs, err := r.CreateLogs()
			if err != nil {
				r.logger.Error("failed to create logs", zap.Error(err))
				continue
			}

			if err := r.logConsumer.ConsumeLogs(ctx, logs); err != nil {
				r.logger.Error("failed to consume logs", zap.Error(err))
			}
		}
	}
}

// CreateLogs creates synthetic logs.
func (r *Receiver) CreateLogs() (plog.Logs, error) {
	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	resourceLogs.Resource().Attributes().FromRaw(r.config.Logs.Resource)
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
	logRecords := scopeLogs.LogRecords().AppendEmpty()

	logRecords.Attributes().FromRaw(r.config.Logs.Attributes)

	switch body := r.config.Logs.Body.(type) {
	case string:
		logRecords.Body().SetStr(body)
	case map[string]any:
		logRecords.Body().SetEmptyMap().FromRaw(body)
	default:
		return logs, fmt.Errorf("invalid body type: %T", body)
	}

	return logs, nil
}
