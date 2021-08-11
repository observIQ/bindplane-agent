package manager

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/jpillora/backoff"
	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/version"
	"github.com/observiq/observiq-collector/manager/message"
	"github.com/observiq/observiq-collector/manager/status"
	"github.com/observiq/observiq-collector/manager/task"
	"github.com/observiq/observiq-collector/manager/websocket"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Manager is the manager for the observiq control plane.
type Manager struct {
	config    *Config
	collector *collector.Collector
	logger    *zap.Logger
}

// New returns a new manager with the supplied config.
func New(config *Config, collector *collector.Collector, logger *zap.Logger) *Manager {
	return &Manager{
		config:    config,
		collector: collector,
		logger:    logger,
	}
}

// Start will start the observiq extension.
func (m *Manager) Run(ctx context.Context) error {
	pipeline := message.NewPipeline(m.config.BufferSize)
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := m.collector.Run()
	if err != nil {
		m.logger.Error("failed to start collector: %s", zap.Error(err))
	}

	group, groupCtx := errgroup.WithContext(cancelCtx)
	group.Go(func() error { return m.handleConnection(groupCtx, pipeline) })
	group.Go(func() error { return m.handleStatus(groupCtx, pipeline) })
	group.Go(func() error { return m.handleTasks(groupCtx, pipeline) })

	return group.Wait()
}

// handleConnection will handle the connection to observiq cloud.
func (m *Manager) handleConnection(ctx context.Context, pipeline *message.Pipeline) error {
	headers := m.headers()
	backoff := backoff.Backoff{Max: m.config.MaxConnectBackoff}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff.Duration()):
			conn, err := websocket.Open(ctx, m.config.Endpoint, headers)
			if err != nil {
				m.logger.Error("Failed to open connection", zap.Error(err), zap.Float64("attempt", backoff.Attempt()))
				continue
			}

			m.logger.Info("Connected to observiq cloud")
			backoff.Reset()

			err = websocket.PumpWithTimeout(ctx, conn, pipeline, m.config.ReconnectInterval)
			switch err {
			case nil, context.DeadlineExceeded, context.Canceled:
			default:
				m.logger.Error("Unexpected connection error", zap.Error(err))
			}
		}
	}
}

// handleStatus will handle reporting status at the specified interval.
func (m *Manager) handleStatus(ctx context.Context, pipeline *message.Pipeline) error {
	ticker := time.NewTicker(m.config.StatusInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			collectorStatus := m.collector.Status()
			m.logger.Info("Collector status", zap.Bool("Running", collectorStatus.Running), zap.Error(collectorStatus.Err))
			report := status.Get()
			err := report.AddPerformanceMetrics()
			if err != nil {
				m.logger.Error("Failed to add performance metrics to status report", zap.Error(err))
				continue
			}

			pipeline.Outbound() <- report.ToMessage()
		}
	}
}

// handleTasks will handle executing tasks from the pipeline.
func (m *Manager) handleTasks(ctx context.Context, pipeline *message.Pipeline) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-pipeline.Inbound():
			t, err := task.FromMessage(msg)
			if err != nil {
				m.logger.Error("Failed to read task", zap.Error(err))
				continue
			}

			m.logger.Info("Executing task", zap.Any("task", msg.Content))
			response, err := task.Execute(t)
			if err != nil {
				m.logger.Error("Failed to execute task", zap.Error(err))
				continue
			}

			pipeline.Outbound() <- response.ToMessage()
		}
	}
}

// headers returns the headers used to connect to observiq cloud.
func (e *Manager) headers() map[string][]string {
	hostname, _ := os.Hostname()
	return map[string][]string{
		"Secret-Key":       {e.config.SecretKey},
		"Agent-Id":         {e.config.AgentID},
		"Hostname":         {hostname},
		"Version":          {version.Version},
		"Operating-System": {runtime.GOOS},
		"Architecture":     {runtime.GOARCH},
	}
}
