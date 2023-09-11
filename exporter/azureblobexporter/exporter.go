package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type azureBlobExporter struct {
	blobClient *azblob.Client
	logger     *zap.Logger
}

func newExporter(config *Config, params exporter.CreateSettings) (*azureBlobExporter, error) {
	blobClient, err := azblob.NewClientFromConnectionString(config.ConnectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client: %w", err)
	}

	return &azureBlobExporter{
		blobClient: blobClient,
		logger:     params.Logger,
	}, nil
}

func (a *azureBlobExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (a *azureBlobExporter) metricsDataPusher(_ context.Context, md pmetric.Metrics) error {
	return nil
}

func (a *azureBlobExporter) logsDataPusher(_ context.Context, ld plog.Logs) error {
	return nil
}

func (a *azureBlobExporter) traceDataPusher(_ context.Context, td ptrace.Traces) error {
	return nil
}
