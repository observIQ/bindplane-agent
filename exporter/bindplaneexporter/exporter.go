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
	mux            sync.RWMutex
	endpoint       string
	client         *websocket.Conn
	liveTailConfig *LiveTailConfig
	viper          *viper.Viper
	logger         *zap.Logger
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
func (e *Exporter) Start(ctx context.Context, _ component.Host) error {
	e.updateLiveTail()
	e.viper.OnConfigChange(func(_ fsnotify.Event) {
		e.updateLiveTail()
	})
	e.viper.WatchConfig()

	client, _, err := websocket.DefaultDialer.DialContext(ctx, e.endpoint, http.Header{})
	if err != nil {
		return err
	}
	e.client = client

	return nil
}

// Shutdown shutsdown the exporter
func (e *Exporter) Shutdown(_ context.Context) error {
	_ = e.client.WriteControl(websocket.CloseMessage, []byte(""), time.Now().Add(10*time.Second))
	return e.client.Close()
}

// updateLiveTail updates the live tail config
func (e *Exporter) updateLiveTail() {
	e.mux.Lock()
	defer e.mux.Unlock()

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

// consumeMetrics will consume metrics
func (e *Exporter) consumeMetrics(_ context.Context, metrics pmetric.Metrics) error {
	e.mux.RLock()
	sessions := e.liveTailConfig.Sessions
	e.mux.RUnlock()

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
		if err := e.client.WriteJSON(msg); err != nil {
			e.logger.Error("Failed to write metrics message", zap.Error(err))
		}
	}

	return nil
}

// consumeLogs will consume logs
func (e *Exporter) consumeLogs(_ context.Context, logs plog.Logs) error {
	e.mux.RLock()
	sessions := e.liveTailConfig.Sessions
	e.mux.RUnlock()

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
		if err := e.client.WriteJSON(msg); err != nil {
			e.logger.Error("Failed to write logs message", zap.Error(err))
		}
	}

	return nil
}

// consumeTraces will consume traces
func (e *Exporter) consumeTraces(_ context.Context, traces ptrace.Traces) error {
	e.mux.RLock()
	sessions := e.liveTailConfig.Sessions
	e.mux.RUnlock()

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
		if err := e.client.WriteJSON(msg); err != nil {
			e.logger.Error("Failed to write traces message", zap.Error(err))
		}
	}

	return nil
}
