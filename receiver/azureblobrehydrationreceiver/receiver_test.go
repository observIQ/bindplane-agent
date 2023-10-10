package azureblobrehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver"

import (
	"testing"

	"github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver/internal/azureblob"
	blobmocks "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver/internal/azureblob/mocks"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.uber.org/zap"
)

func Test_newMetricsReceiver(t *testing.T) {
	mockClient := setNewAzureBlobClient(t)
	id := component.NewID("test")
	testLogger := zap.NewNop()
	cfg := &Config{
		StartingTime: "2023-10-02T17:00",
		EndingTime:   "2023-10-02T17:01",
	}
	co := consumertest.NewNop()
	r, err := newMetricsReceiver(id, testLogger, cfg, co)
	require.NoError(t, err)

	require.Equal(t, testLogger, r.logger)
	require.Equal(t, id, r.id)
	require.Equal(t, mockClient, r.azureClient)
	require.Equal(t, component.DataTypeMetrics, r.supportedTelemetry)
	require.IsType(t, &metricsConsumer{}, r.consumer)
}

func Test_newLogsReceiver(t *testing.T) {
	mockClient := setNewAzureBlobClient(t)
	id := component.NewID("test")
	testLogger := zap.NewNop()
	cfg := &Config{
		StartingTime: "2023-10-02T17:00",
		EndingTime:   "2023-10-02T17:01",
	}
	co := consumertest.NewNop()
	r, err := newLogsReceiver(id, testLogger, cfg, co)
	require.NoError(t, err)

	require.Equal(t, testLogger, r.logger)
	require.Equal(t, id, r.id)
	require.Equal(t, mockClient, r.azureClient)
	require.Equal(t, component.DataTypeLogs, r.supportedTelemetry)
	require.IsType(t, &logsConsumer{}, r.consumer)
}

func Test_newTracesReceiver(t *testing.T) {
	mockClient := setNewAzureBlobClient(t)
	id := component.NewID("test")
	testLogger := zap.NewNop()
	cfg := &Config{
		StartingTime: "2023-10-02T17:00",
		EndingTime:   "2023-10-02T17:01",
	}
	co := consumertest.NewNop()
	r, err := newTracesReceiver(id, testLogger, cfg, co)
	require.NoError(t, err)

	require.Equal(t, testLogger, r.logger)
	require.Equal(t, id, r.id)
	require.Equal(t, mockClient, r.azureClient)
	require.Equal(t, component.DataTypeTraces, r.supportedTelemetry)
	require.IsType(t, &tracesConsumer{}, r.consumer)
}

func setNewAzureBlobClient(t *testing.T) *blobmocks.MockBlobClient {
	t.Helper()
	oldfunc := newAzureBlobClient

	mockClient := blobmocks.NewMockBlobClient(t)

	newAzureBlobClient = func(_ string) (azureblob.BlobClient, error) {
		return mockClient, nil
	}

	t.Cleanup(func() {
		newAzureBlobClient = oldfunc
	})

	return mockClient
}
