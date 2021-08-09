package status

import (
	"context"
	"testing"

	"github.com/observiq/observiq-collector/manager/message"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	ctx := context.TODO()
	report, err := Get(ctx)
	require.NoError(t, err)
	require.Equal(t, report.ComponentID, "bpagent")
}

func TestReportToMessage(t *testing.T) {
	ctx := context.TODO()
	report, err := Get(ctx)
	require.NoError(t, err)

	msg := report.ToMessage()
	require.Equal(t, msg.Type, message.StatusReport)
}
