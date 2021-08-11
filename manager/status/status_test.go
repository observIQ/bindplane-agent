package status

import (
	"errors"
	"testing"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/manager/message"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	status := collector.Status{Err: nil, Running: true}
	report := Get("41794e4d-9564-4d98-9096-698302577c98", status)
	require.Equal(t, report.ComponentID, "41794e4d-9564-4d98-9096-698302577c98")
}

func TestReportToMessage(t *testing.T) {
	status := collector.Status{Err: nil}
	report := Get("41794e4d-9564-4d98-9096-698302577c98", status)

	msg := report.ToMessage()
	require.Equal(t, msg.Type, message.StatusReport)
}

func TestErrorToMessage(t *testing.T) {
	status := collector.Status{Err: errors.New("Error for testing")}
	report := Get("41794e4d-9564-4d98-9096-698302577c98", status)

	msg := report.ToMessage()
	require.Equal(t, msg.Content["status"], ERROR)
}
