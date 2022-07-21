package bindplaneexporter

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// Exporter is the bindplane exporter
type Exporter struct {
	mux      sync.RWMutex
	sessions []Session
	client   *websocket.Conn
}

// NewExporter creates a new Exporter with a connection to the given endpoint
func NewExporter(ctx context.Context, cfg Config) (*Exporter, error) {
	wsClient, _, err := websocket.DefaultDialer.DialContext(ctx, cfg.Endpoint, http.Header{"Agent-ID": []string{"TODO"}})
	if err != nil {
		return nil, err
	}
	return &Exporter{
		client: wsClient,
	}, nil
}

func (e *Exporter) Stop() error {
	_ = e.client.WriteControl(websocket.CloseMessage, []byte(""), time.Now().Add(10*time.Second))
	return e.client.Close()
}

// ConsumeMetrics will consume metrics
func (e *Exporter) ConsumeMetrics(_ context.Context, metrics pmetric.Metrics) error {
	e.mux.RLock()
	defer e.mux.RUnlock()

	if len(e.sessions) == 0 {
		return nil
	}

	records := getRecordsFromMetrics(metrics)
	for _, record := range records {
		sessionIDs := []string{}
		for _, session := range e.sessions {
			if session.Matches(record) {
				sessionIDs = append(sessionIDs, session.ID)
			}
		}

		if len(sessionIDs) == 0 {
			continue
		}

		msg := NewMetricsMessage(record, sessionIDs)
		e.client.WriteJSON(msg)
	}

	return nil
}

// ConsumeLogs will consume logs
func (e *Exporter) ConsumeLogs(_ context.Context, logs plog.Logs) error {
	e.mux.RLock()
	defer e.mux.RUnlock()

	if len(e.sessions) == 0 {
		return nil
	}

	records := getRecordsFromLogs(logs)
	for _, record := range records {
		sessionIDs := []string{}
		for _, session := range e.sessions {
			if session.Matches(record) {
				sessionIDs = append(sessionIDs, session.ID)
			}
		}

		if len(sessionIDs) == 0 {
			continue
		}

		msg := NewLogsMessage(record, sessionIDs)
		e.client.WriteJSON(msg)
	}

	return nil
}

// ConsumeTraces will consume traces
func (e *Exporter) ConsumeTraces(_ context.Context, traces ptrace.Traces) error {
	e.mux.RLock()
	defer e.mux.RUnlock()

	if len(e.sessions) == 0 {
		return nil
	}

	records := getRecordsFromTraces(traces)
	for _, record := range records {
		sessionIDs := []string{}
		for _, session := range e.sessions {
			if session.Matches(record) {
				sessionIDs = append(sessionIDs, session.ID)
			}
		}

		if len(sessionIDs) == 0 {
			continue
		}

		msg := NewTracesMessage(record, sessionIDs)
		e.client.WriteJSON(msg)
	}

	return nil
}
