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
	report := Get(status)
	require.Equal(t, report.ComponentID, "bpagent")
}

func TestReportToMessage(t *testing.T) {
	status := collector.Status{Err: nil}
	report := Get(status)

	msg := report.ToMessage()
	require.Equal(t, msg.Type, message.StatusReport)
}

func TestErrorToMessage(t *testing.T) {
	status := collector.Status{Err: errors.New("Error for testing")}
	report := Get(status)

	msg := report.ToMessage()
	require.Equal(t, msg.Content["status"], ERROR)
}
