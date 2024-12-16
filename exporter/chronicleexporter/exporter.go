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
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/observiq/bindplane-agent/exporter/chronicleexporter/protos/api"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	grpcgzip "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	grpcScope = "https://www.googleapis.com/auth/malachite-ingestion"
	httpScope = "https://www.googleapis.com/auth/cloud-platform"
)

type chronicleExporter struct {
	cfg        *Config
	set        component.TelemetrySettings
	marshaler  logMarshaler
	exporterID string

	// fields used for gRPC
	grpcClient api.IngestionServiceV2Client
	grpcConn   *grpc.ClientConn
	metrics    *hostMetricsReporter

	// fields used for HTTP
	httpClient *http.Client
}

func newExporter(cfg *Config, params exporter.Settings, exporterID string) (*chronicleExporter, error) {
	customerID, err := uuid.Parse(cfg.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("parse customer ID: %w", err)
	}

	marshaller, err := newProtoMarshaler(*cfg, params.TelemetrySettings, customerID[:])
	if err != nil {
		return nil, fmt.Errorf("create proto marshaller: %w", err)
	}

	return &chronicleExporter{
		cfg:        cfg,
		set:        params.TelemetrySettings,
		marshaler:  marshaller,
		exporterID: exporterID,
	}, nil
}

func (ce *chronicleExporter) Start(_ context.Context, _ component.Host) error {
	creds, err := loadGoogleCredentials(ce.cfg)
	if err != nil {
		return fmt.Errorf("load Google credentials: %w", err)
	}

	if ce.cfg.Protocol == protocolHTTPS {
		ce.httpClient = oauth2.NewClient(context.Background(), creds.TokenSource)
		return nil
	}

	opts := []grpc.DialOption{
		// Apply OAuth tokens for each RPC call
		grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: creds.TokenSource}),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	}
	conn, err := grpc.NewClient(ce.cfg.Endpoint+":443", opts...)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	ce.grpcConn = conn
	ce.grpcClient = api.NewIngestionServiceV2Client(conn)

	if ce.cfg.CollectAgentMetrics {
		f := func(ctx context.Context, request *api.BatchCreateEventsRequest) error {
			_, err := ce.grpcClient.BatchCreateEvents(ctx, request)
			return err
		}
		metrics, err := newHostMetricsReporter(ce.cfg, ce.set, ce.exporterID, f)
		if err != nil {
			return fmt.Errorf("create metrics reporter: %w", err)
		}
		ce.metrics = metrics
		ce.metrics.start()
	}

	return nil
}

func (ce *chronicleExporter) Shutdown(context.Context) error {
	defer http.DefaultTransport.(*http.Transport).CloseIdleConnections()
	if ce.cfg.Protocol == protocolHTTPS {
		t := ce.httpClient.Transport.(*oauth2.Transport)
		if t.Base != nil {
			t.Base.(*http.Transport).CloseIdleConnections()
		}
		return nil
	}
	if ce.metrics != nil {
		ce.metrics.shutdown()
	}
	if ce.grpcConn != nil {
		if err := ce.grpcConn.Close(); err != nil {
			return fmt.Errorf("connection close: %s", err)
		}
	}
	return nil
}

func (ce *chronicleExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func loadGoogleCredentials(cfg *Config) (*google.Credentials, error) {
	scope := grpcScope
	if cfg.Protocol == protocolHTTPS {
		scope = httpScope
	}

	switch {
	case cfg.Creds != "":
		return google.CredentialsFromJSON(context.Background(), []byte(cfg.Creds), scope)
	case cfg.CredsFilePath != "":
		credsData, err := os.ReadFile(cfg.CredsFilePath)
		if err != nil {
			return nil, fmt.Errorf("read credentials file: %w", err)
		}

		if len(credsData) == 0 {
			return nil, errors.New("credentials file is empty")
		}

		return google.CredentialsFromJSON(context.Background(), credsData, scope)
	default:
		return google.FindDefaultCredentials(context.Background(), scope)
	}
}

func (ce *chronicleExporter) logsDataPusher(ctx context.Context, ld plog.Logs) error {
	payloads, err := ce.marshaler.MarshalRawLogs(ctx, ld)
	if err != nil {
		return fmt.Errorf("marshal logs: %w", err)
	}

	for _, payload := range payloads {
		if err := ce.uploadToChronicle(ctx, payload); err != nil {
			return fmt.Errorf("upload to chronicle: %w", err)
		}
	}

	return nil
}

func (ce *chronicleExporter) uploadToChronicle(ctx context.Context, request *api.BatchCreateLogsRequest) error {
	if ce.metrics != nil {
		totalLogs := int64(len(request.GetBatch().GetEntries()))
		defer ce.metrics.recordSent(totalLogs)
	}

	_, err := ce.grpcClient.BatchCreateLogs(ctx, request, ce.buildOptions()...)
	if err != nil {
		errCode := status.Code(err)
		switch errCode {
		// These errors are potentially transient
		case codes.Canceled,
			codes.Unavailable,
			codes.DeadlineExceeded,
			codes.ResourceExhausted,
			codes.Aborted:
			return fmt.Errorf("upload logs to chronicle: %w", err)
		default:
			return consumererror.NewPermanent(fmt.Errorf("upload logs to chronicle: %w", err))
		}
	}

	return nil
}

func (ce *chronicleExporter) buildOptions() []grpc.CallOption {
	opts := make([]grpc.CallOption, 0)

	if ce.cfg.Compression == grpcgzip.Name {
		opts = append(opts, grpc.UseCompressor(grpcgzip.Name))
	}

	return opts
}

func (ce *chronicleExporter) logsHTTPDataPusher(ctx context.Context, ld plog.Logs) error {
	payloads, err := ce.marshaler.MarshalRawLogsForHTTP(ctx, ld)
	if err != nil {
		return fmt.Errorf("marshal logs: %w", err)
	}

	for logType, logTypePayloads := range payloads {
		for _, payload := range logTypePayloads {
			if err := ce.uploadToChronicleHTTP(ctx, payload, logType); err != nil {
				return fmt.Errorf("upload to chronicle: %w", err)
			}
		}
	}

	return nil
}

// This uses the DataPlane URL for the request
// URL for the request: https://{region}-chronicle.googleapis.com/{version}/projects/{project}/location/{region}/instances/{customerID}/logTypes/{logtype}/logs:import
func buildEndpoint(cfg *Config, logType string) string {
	//                Location Endpoint Version    Project      Location    Instance     LogType
	formatString := "https://%s-%s/%s/projects/%s/locations/%s/instances/%s/logTypes/%s/logs:import"
	return fmt.Sprintf(formatString, cfg.Location, cfg.Endpoint, "v1alpha", cfg.Project, cfg.Location, cfg.CustomerID, logType)
}

func (ce *chronicleExporter) uploadToChronicleHTTP(ctx context.Context, logs *api.ImportLogsRequest, logType string) error {

	data, err := protojson.Marshal(logs)
	if err != nil {
		return fmt.Errorf("marshal protobuf logs to JSON: %w", err)
	}

	var body io.Reader

	if ce.cfg.Compression == grpcgzip.Name {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write(data); err != nil {
			return fmt.Errorf("gzip write: %w", err)
		}
		if err := gz.Close(); err != nil {
			return fmt.Errorf("gzip close: %w", err)
		}
		body = &b
	} else {
		body = bytes.NewBuffer(data)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", buildEndpoint(ce.cfg, logType), body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if ce.cfg.Compression == grpcgzip.Name {
		request.Header.Set("Content-Encoding", "gzip")
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := ce.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send request to Chronicle: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			ce.set.Logger.Warn("Failed to read response body", zap.Error(err))
		} else {
			ce.set.Logger.Warn("Received non-OK response from Chronicle", zap.String("status", resp.Status), zap.ByteString("response", respBody))
		}
		return fmt.Errorf("received non-OK response from Chronicle: %s", resp.Status)
	}

	return nil
}
