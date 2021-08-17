package startup

import (
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/observiq/observiq-collector/collector"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestStartupCreation(t *testing.T) {
	templateID := uuid.New().String()
	name := "test name"
	col := collector.New(path.Join(".", "testdata", "valid_config.yaml"), []zap.Option{})
	statusReport := New(templateID, name, col)
	require.Equal(t, statusReport.AgentName, name)
	require.Equal(t, statusReport.TemplateID, templateID)

	// unmarshall test
	msg := statusReport.ToMessage()
	require.Equal(t, msg.Type, "onStartup")
	require.NotEmpty(t, msg.Content["osDetails"])
	require.Equal(t, msg.Content["agentName"], name)
	require.Equal(t, msg.Content["configuration_id"], templateID)
}

func TestStartupCreation_NoTemplateID(t *testing.T) {
	name := "no template id"
	col := collector.New(path.Join(".", "testdata", "valid_config.yaml"), []zap.Option{})
	statusReport := New("", name, col)
	require.Equal(t, statusReport.TemplateID, "")

	content := statusReport.ToMessage().Content
	require.Equal(t, content["configuration_id"], "")
}
