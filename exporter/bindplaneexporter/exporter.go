package bindplaneexporter

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// Exporter is the bindplane exporter
type Exporter struct {
	endpoint       string
	client         *websocket.Conn
	clientMux      sync.RWMutex
	liveTailConfig *LiveTailConfig
	liveTailMux    sync.RWMutex
	viper          *viper.Viper
	logger         *zap.Logger
	cancel         context.CancelFunc
}

// NewExporter creates a new Exporter with a connection to the given endpoint
func NewExporter(_ context.Context, cfg *Config, set component.ExporterCreateSettings) (*Exporter, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(cfg.LiveTail)

	return &Exporter{
		endpoint: cfg.Endpoint,
		viper:    v,
		logger:   set.Logger.Named(cfg.ID().String()),
	}, nil
}

// Start starts the exporter
func (e *Exporter) Start(_ context.Context, _ component.Host) error {
	e.updateLiveTail()
	e.viper.OnConfigChange(func(_ fsnotify.Event) {
		e.updateLiveTail()
	})
	e.viper.WatchConfig()

	ctx, cancel := context.WithCancel(context.Background())
	go e.handleWebsocketOpen(ctx)
	e.cancel = cancel

	return nil
}

// Shutdown shutsdown the exporter
func (e *Exporter) Shutdown(_ context.Context) error {
	e.clientMux.RLock()
	defer e.clientMux.RUnlock()

	e.cancel()

	if e.client == nil {
		return nil
	}

	_ = e.client.WriteControl(websocket.CloseMessage, []byte(""), time.Now().Add(10*time.Second))
	return e.client.Close()
}

// updateLiveTail updates the live tail config
func (e *Exporter) updateLiveTail() {
	e.liveTailMux.Lock()
	defer e.liveTailMux.Unlock()

	e.logger.Info("Updating live tail config")
	if err := e.viper.ReadInConfig(); err != nil {
		e.logger.Error("Failed to read in live tail config", zap.Error(err))
		return
	}

	cfg := &LiveTailConfig{}
	if err := e.viper.Unmarshal(cfg); err != nil {
		e.logger.Error("Failed to unmarshal config", zap.Error(err))
		return
	}

	e.logger.Info("Updated live tail config")
	e.liveTailConfig = cfg
}

// handleWebsocketOpen handles opening the websocket
func (e *Exporter) handleWebsocketOpen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		client, _, err := websocket.DefaultDialer.DialContext(ctx, e.endpoint, http.Header{})
		if err != nil {
			e.logger.Error("Failed to create ws client", zap.Error(err))
			time.Sleep(time.Second * 5)
			continue
		}

		e.clientMux.Lock()
		e.client = client
		e.clientMux.Unlock()

		if _, _, err = client.ReadMessage(); err != nil {
			e.logger.Error("Websocket client received error", zap.Error(err))
		}

		e.clientMux.Lock()
		e.client = nil
		e.clientMux.Unlock()
	}
}

// sendMessage sends a websocket message
func (e *Exporter) sendMessage(msg Message) {
	e.clientMux.RLock()
	defer e.clientMux.RUnlock()

	if e.client == nil {
		e.logger.Info("Failed to send message. Client not open.")
		return
	}

	if err := e.client.WriteJSON(msg); err != nil {
		e.logger.Error("Failed to write message", zap.Error(err))
	}
}

// consumeMetrics will consume metrics
func (e *Exporter) consumeMetrics(_ context.Context, metrics pmetric.Metrics) error {
	e.liveTailMux.RLock()
	sessions := e.liveTailConfig.Sessions
	e.liveTailMux.RUnlock()

	if len(sessions) == 0 {
		e.logger.Info("Skipping consume metrics. No active sessions.")
		return nil
	}

	records := getRecordsFromMetrics(metrics)
	for _, record := range records {
		sessionIDs := []string{}
		for _, session := range sessions {
			if session.Matches(record) {
				sessionIDs = append(sessionIDs, session.ID)
			}
		}

		if len(sessionIDs) == 0 {
			e.logger.Info("No matching sessions. Skipping metric record.")
			continue
		}

		msg := NewMetricsMessage(record, sessionIDs)
		e.sendMessage(msg)
	}

	return nil
}

// consumeLogs will consume logs
func (e *Exporter) consumeLogs(_ context.Context, logs plog.Logs) error {
	e.liveTailMux.RLock()
	sessions := e.liveTailConfig.Sessions
	e.liveTailMux.RUnlock()

	if len(sessions) == 0 {
		e.logger.Info("Skipping consume logs. No active sessions.")
		return nil
	}

	records := getRecordsFromLogs(logs)
	for _, record := range records {
		sessionIDs := []string{}
		for _, session := range sessions {
			if session.Matches(record) {
				sessionIDs = append(sessionIDs, session.ID)
			}
		}

		if len(sessionIDs) == 0 {
			e.logger.Info("No matching sessions. Skipping log record.")
			continue
		}

		msg := NewLogsMessage(record, sessionIDs)
		e.sendMessage(msg)
	}

	return nil
}

// consumeTraces will consume traces
func (e *Exporter) consumeTraces(_ context.Context, traces ptrace.Traces) error {
	e.liveTailMux.RLock()
	sessions := e.liveTailConfig.Sessions
	e.liveTailMux.RUnlock()

	if len(sessions) == 0 {
		e.logger.Info("Skipping consume traces. No active sessions.")
		return nil
	}

	records := getRecordsFromTraces(traces)
	for _, record := range records {
		sessionIDs := []string{}
		for _, session := range sessions {
			if session.Matches(record) {
				sessionIDs = append(sessionIDs, session.ID)
			}
		}

		if len(sessionIDs) == 0 {
			e.logger.Info("No matching sessions. Skipping trace record.")
			continue
		}

		msg := NewTracesMessage(record, sessionIDs)
		e.sendMessage(msg)
	}

	return nil
}
