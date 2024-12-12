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
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/observiq/bindplane-otel-collector/exporter/chronicleexporter/protos/api"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
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
	cfg                     *Config
	set                     component.TelemetrySettings
	marshaler               logMarshaler
	collectorID, exporterID string

	// fields used for gRPC
	grpcClient api.IngestionServiceV2Client
	grpcConn   *grpc.ClientConn
	wg         sync.WaitGroup
	cancel     context.CancelFunc
	metrics    *exporterMetrics

	// fields used for HTTP
	httpClient *http.Client
}

func newExporter(cfg *Config, params exporter.Settings, collectorID, exporterID string) (*chronicleExporter, error) {
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

func (ce *chronicleExporter) Start(ctx context.Context, _ component.Host) error {
	ts, err := tokenSource(ctx, ce.cfg)
	if err != nil {
		return fmt.Errorf("load Google credentials: %w", err)
	}

	return &chronicleExporter{
		cfg:         cfg,
		set:         params.TelemetrySettings,
		metrics:     newHostMetricsReporter(uuidCID[:], customerID[:], exporterID, cfg.Namespace),
		marshaler:   marshaller,
		collectorID: collectorID,
		exporterID:  exporterID,
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
		ctx, cancel := context.WithCancel(context.Background())
		ce.cancel = cancel
		ce.wg.Add(1)
		go ce.startHostMetricsCollection(ctx)
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
	if ce.cancel != nil {
		ce.cancel()
		ce.wg.Wait()
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
			return err
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

func (ce *chronicleExporter) startHostMetricsCollection(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	defer ce.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := ce.metrics.collectHostMetrics()
			if err != nil {
				ce.set.Logger.Error("Failed to collect host metrics", zap.Error(err))
			}
			request := ce.metrics.getAndReset()
			_, err = ce.grpcClient.BatchCreateEvents(ctx, request, ce.buildOptions()...)
			if err != nil {
				ce.set.Logger.Error("Failed to upload host metrics", zap.Error(err))
			}
		}
	}
}

func (ce *chronicleExporter) logsHTTPDataPusher(ctx context.Context, ld plog.Logs) error {
	payloads, err := ce.marshaler.MarshalRawLogsForHTTP(ctx, ld)
	if err != nil {
		return fmt.Errorf("marshal logs: %w", err)
	}

	for logType, logTypePayloads := range payloads {
		for _, payload := range logTypePayloads {
			if err := ce.uploadToChronicleHTTP(ctx, payload, logType); err != nil {
				return err
			}
		}
	}

	return nil
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

	request, err := http.NewRequestWithContext(ctx, "POST", httpEndpoint(ce.cfg, logType), body)
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
