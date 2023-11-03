package chronicleexporter

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const scope = "https://www.googleapis.com/auth/malachite-ingestion"

type chronicleExporter struct {
	cfg        *Config
	logger     *zap.Logger
	httpClient *http.Client
	marshaler  logMarshaler
	endpoint   string
}

func newExporter(cfg *Config, params exporter.CreateSettings) (*chronicleExporter, error) {
	var creds *google.Credentials
	var err error

	switch {
	case cfg.Creds != "":
		creds, err = google.CredentialsFromJSON(context.Background(), []byte(cfg.CredsFilePath), scope)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain credentials from JSON: %w", err)
		}
	case cfg.CredsFilePath != "":
		credsData, err := os.ReadFile(cfg.CredsFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file: %w", err)
		}

		creds, err = google.CredentialsFromJSON(context.Background(), credsData, scope)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain credentials from JSON: %w", err)
		}
	default:
		return nil, fmt.Errorf("no credentials provided")
	}

	// Use the credentials to create an HTTP client
	httpClient := oauth2.NewClient(context.Background(), creds.TokenSource)

	return &chronicleExporter{
		endpoint:   regions[cfg.Region],
		cfg:        cfg,
		logger:     params.Logger,
		httpClient: httpClient,
		marshaler:  newMarshaler(*cfg),
	}, nil
}

func (ce *chronicleExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (ce *chronicleExporter) logsDataPusher(ctx context.Context, ld plog.Logs) error {
	udmData, err := ce.marshaler.MarshalRawLogs(ld)
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}

	return ce.uploadToChronicle(ctx, udmData)
}

func (ce *chronicleExporter) uploadToChronicle(ctx context.Context, data []byte) error {
	request, err := http.NewRequestWithContext(ctx, "POST", ce.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := ce.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request to Chronicle: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Chronicle: %s", resp.Status)
	}

	return nil
}
