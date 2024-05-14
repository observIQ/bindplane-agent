package snapshotprocessor

import (
	"encoding/json"
)

type snapshotReport struct {
	SessionID     string `json:"session_id"`
	TelemetryType string `json:"telemetry_type"`
	// TelemetryPayload is the logs/metrics/traces in OTLP/JSON format
	TelemetryPayload json.RawMessage `json:"telemetry_payload"`
}

func logsReport(sessionID string, logsPayload []byte) snapshotReport {
	return snapshotReport{
		SessionID:        sessionID,
		TelemetryType:    "logs",
		TelemetryPayload: logsPayload,
	}
}

func metricsReport(sessionID string, metricsPayload []byte) snapshotReport {
	return snapshotReport{
		SessionID:        sessionID,
		TelemetryType:    "metrics",
		TelemetryPayload: metricsPayload,
	}
}

func tracesReport(sessionID string, tracesPayload []byte) snapshotReport {
	return snapshotReport{
		SessionID:        sessionID,
		TelemetryType:    "traces",
		TelemetryPayload: tracesPayload,
	}
}
