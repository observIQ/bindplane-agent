package manager

import (
	"context"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/manager/message"
	"github.com/observiq/observiq-collector/manager/startup"
	"github.com/observiq/observiq-collector/manager/status"
	"github.com/observiq/observiq-collector/manager/task"
	"github.com/observiq/observiq-collector/manager/websocket"
	"go.uber.org/zap"
)

// Manager is the manager for the observiq control plane.
type Manager struct {
	config    *Config
	collector *collector.Collector
	logger    *zap.Logger
	in        chan *message.Message
	out       chan *message.Message
	exit      chan int
}

// New returns a new manager with the supplied parameters.
func New(config *Config, collector *collector.Collector, logger *zap.Logger) *Manager {
	return &Manager{
		config:    config,
		collector: collector,
		logger:    logger,
		in:        make(chan *message.Message, config.BufferSize),
		out:       make(chan *message.Message, config.BufferSize),
		exit:      make(chan int, 1),
	}
}

// Run will run the manager until the supplied context is cancelled
// or an exit code is received.
func (m *Manager) Run(ctx context.Context) (exitCode int) {
	groupCtx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	wg.Add(4)
	go m.handleCollector(groupCtx, wg)
	go m.handleConnection(groupCtx, wg)
	go m.handleStatus(groupCtx, wg)
	go m.handleTasks(groupCtx, wg)

	select {
	case <-ctx.Done():
	case exitCode = <-m.exit:
	}

	cancel()
	wg.Wait()
	m.drainMessages()

	return
}

// handleCollector will handle starting the collector.
// When the supplied context is cancelled, the collector handler will stop.
func (m *Manager) handleCollector(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := m.logger.Named("collector-handler")
	logger.Info("Starting collector handler")
	defer logger.Info("Exiting collector handler")

	err := m.collector.Run()
	if err != nil {
		logger.Error("failed to start collector: %s", zap.Error(err))
	}

	startupMsg := startup.New(m.config.TemplateID, m.config.AgentName, m.collector).ToMessage()
	m.out <- startupMsg

	<-ctx.Done()
	m.collector.Stop()
}

// handleConnection will handle the connection to the observiq control plane.
// If an error occurs while connecting, the connection is retried with backoff.
// This connection is periodically refreshed based on the configured reconnect interval.
// When the supplied context is cancelled, the connection handler will stop.
func (m *Manager) handleConnection(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := m.logger.Named("connection-handler")
	logger.Info("Starting connection handler")
	defer logger.Info("Exiting connection handler")

	backoff := backoff.Backoff{Max: m.config.MaxConnectBackoff}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff.Duration()):
			conn, err := websocket.Open(ctx, m.config.Endpoint, m.config.Headers())
			if err != nil {
				logger.Error("Failed to open connection", zap.Error(err), zap.Float64("attempt", backoff.Attempt()))
				continue
			}

			logger.Info("Connection started")
			backoff.Reset()

			timedCtx, cancel := context.WithTimeout(ctx, m.config.ReconnectInterval)
			err = websocket.HandleTraffic(timedCtx, conn, m.in, m.out)
			cancel()
			logger.Info("Connection stopped")

			switch err {
			case nil, context.DeadlineExceeded, context.Canceled:
			default:
				logger.Error("Received connection error", zap.Error(err))
			}
		}
	}
}

// handleStatus will handle reporting the agent's status to the outbound channel.
// When the supplied context is cancelled, the status handler will stop.
func (m *Manager) handleStatus(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := m.logger.Named("status-handler")
	logger.Info("Starting status handler")
	defer logger.Info("Exiting status handler")

	ticker := time.NewTicker(m.config.StatusInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Debug("Running status update")

			collectorStatus := m.collector.Status()
			if collectorStatus.Err != nil {
				logger.Info("Collector status", zap.Bool("Running", collectorStatus.Running), zap.Error(collectorStatus.Err))
			}
			report := status.Get(m.config.AgentID, collectorStatus)
			report.AddPerformanceMetrics(logger)

			m.out <- report.ToMessage()
		}
	}
}

// handleTasks will handle executing tasks received from the inbound channel
// and sending a response to the outbound channel. When the supplied context
// is cancelled, the task handler will stop.
func (m *Manager) handleTasks(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := m.logger.Named("task-handler")
	logger.Info("Starting task handler")
	defer logger.Info("Exiting task handler")

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-m.in:
			logger.Debug("Received message")

			t, err := task.FromMessage(msg)
			if err != nil {
				logger.Error("Failed to convert message to task", zap.Error(err))
				continue
			}

			logger.Info("Executing task", zap.Any("task", msg.Content))
			response := m.executeTask(t)

			logger.Info("Sending task response", zap.Any("response", response))
			m.out <- response.ToMessage()
		}
	}
}

// executeTask will execute a task with the manager.
func (m *Manager) executeTask(t *task.Task) task.Response {
	switch t.Type {
	case task.Reconfigure:
		return task.ExecuteReconfigure(t, m.collector)
	case task.Restart:
		return task.ExecuteRestart(t, m.collector)
	case task.Shutdown:
		return task.ExecuteShutdown(t, m.exit)
	default:
		return task.ExecuteUnknown(t)
	}
}

// drainMessages will attempt to drain outbound messages during a shutdown.
// This operation will timeout after 10 seconds if not completed.
func (m *Manager) drainMessages() {
	logger := m.logger.Named("message-drainer")
	logger.Info("Starting message drainer")
	defer logger.Info("Exiting message drainer")

	close(m.out)
	if len(m.out) == 0 {
		logger.Info("No messages to drain")
		return
	}

	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := websocket.Open(ctx, m.config.Endpoint, m.config.Headers())
	if err != nil {
		logger.Error("Failed to open connection", zap.Error(err))
		return
	}
	defer websocket.Close(conn)

	if err := websocket.HandleSend(ctx, conn, m.out); err != nil {
		logger.Error("Received error while draining", zap.Error(err))
	}
}
