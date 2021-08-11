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

// New returns a new manager with the supplied parameters.
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

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error { return m.handleCollector(groupCtx) })
	group.Go(func() error { return m.handleConnection(groupCtx, pipeline) })
	group.Go(func() error { return m.handleStatus(groupCtx, pipeline) })
	group.Go(func() error { return m.handleTasks(groupCtx, pipeline) })

	return group.Wait()
}

// handleCollector will handle running the collector.
func (m *Manager) handleCollector(ctx context.Context) error {
	logger := m.logger.Named("collector-handler")
	logger.Info("Starting collector handler")
	defer logger.Info("Exiting collector handler")

	err := m.collector.Run()
	if err != nil {
		logger.Error("failed to start collector: %s", zap.Error(err))
	}

	<-ctx.Done()
	m.collector.Stop()
	return ctx.Err()
}

// handleConnection will handle the connection to observiq cloud.
func (m *Manager) handleConnection(ctx context.Context, pipeline *message.Pipeline) error {
	logger := m.logger.Named("connection-handler")
	logger.Info("Starting connection handler")
	defer logger.Info("Exiting connection handler")

	headers := m.headers()
	backoff := backoff.Backoff{Max: m.config.MaxConnectBackoff}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff.Duration()):
			conn, err := websocket.Open(ctx, m.config.Endpoint, headers)
			if err != nil {
				logger.Error("Failed to open connection", zap.Error(err), zap.Float64("attempt", backoff.Attempt()))
				continue
			}

			logger.Info("Connected to observiq cloud")
			backoff.Reset()

			err = websocket.PumpWithTimeout(ctx, conn, pipeline, m.config.ReconnectInterval)
			switch err {
			case nil, context.DeadlineExceeded, context.Canceled:
			default:
				logger.Error("Unexpected connection error", zap.Error(err))
			}
		}
	}
}

// handleStatus will handle reporting status at the specified interval.
func (m *Manager) handleStatus(ctx context.Context, pipeline *message.Pipeline) error {
	logger := m.logger.Named("status-handler")
	logger.Info("Starting status handler")
	defer logger.Info("Exiting status handler")

	ticker := time.NewTicker(m.config.StatusInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			logger.Debug("Running status update")

			collectorStatus := m.collector.Status()
			logger.Info("Collector status", zap.Bool("Running", collectorStatus.Running), zap.Error(collectorStatus.Err))

			report, err := status.Get()
			if err != nil {
				logger.Error("Failed to report status", zap.Error(err))
				continue
			}

			pipeline.Outbound() <- report.ToMessage()
		}
	}
}

// handleTasks will handle executing tasks from the pipeline.
func (m *Manager) handleTasks(ctx context.Context, pipeline *message.Pipeline) error {
	logger := m.logger.Named("task-handler")
	logger.Info("Starting task handler")
	defer logger.Info("Exiting task handler")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-pipeline.Inbound():
			logger.Debug("Received message")

			t, err := task.FromMessage(msg)
			if err != nil {
				logger.Error("Failed to covert message to task", zap.Error(err))
				continue
			}

			logger.Info("Executing task", zap.Any("task", msg.Content))
			response := m.executeTask(t)

			logger.Info("Sending task response", zap.Any("response", response))
			pipeline.Outbound() <- response.ToMessage()
		}
	}
}

// executeTask will execute a task with the manager.
func (m *Manager) executeTask(t *task.Task) task.Response {
	switch t.Type {
	case task.Reconfigure:
		return task.ExecuteReconfigure(t, m.collector)
	default:
		return task.ExecuteUnknown(t)
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
