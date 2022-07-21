package bindplaneexporter

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// client is a stand-in interface until the websocket implementation is complete
type client interface {
	send(message Message)
}

// Exporter is the bindplane exporter
type Exporter struct {
	mux      sync.RWMutex
	sessions []Session
	client   client
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
		e.client.send(msg)
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
		e.client.send(msg)
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
		e.client.send(msg)
	}

	return nil
}
