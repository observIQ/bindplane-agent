// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

const endpoint = "https://malachiteingestion-pa.googleapis.com/v2/unstructuredlogentries:batchCreate"

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
			return nil, fmt.Errorf("obtain credentials from JSON: %w", err)
		}
	case cfg.CredsFilePath != "":
		credsData, err := os.ReadFile(cfg.CredsFilePath)
		if err != nil {
			return nil, fmt.Errorf("read credentials file: %w", err)
		}

		scopes := []string{scope}
		if cfg.Region != "" {
			scopes = append(scopes, regions[cfg.Region])
		}

		creds, err = google.CredentialsFromJSON(context.Background(), credsData, scope)
		if err != nil {
			return nil, fmt.Errorf("obtain credentials from JSON: %w", err)
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
		return fmt.Errorf("marshal logs: %w", err)
	}

	return ce.uploadToChronicle(ctx, udmData)
}

func (ce *chronicleExporter) uploadToChronicle(ctx context.Context, data []byte) error {
	request, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := ce.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send request to Chronicle: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody := []byte{}
		_, err := resp.Body.Read(respBody)
		if err != nil {
			ce.logger.Warn("Failed to read response body", zap.Error(err))
		} else {
			ce.logger.Warn("Received non-OK response from Chronicle", zap.String("status", resp.Status), zap.ByteString("body", respBody))
		}

		return fmt.Errorf("received non-OK response from Chronicle: %s", resp.Status)
	}

	return nil
}
