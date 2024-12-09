// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chronicleexporter

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/observiq/bindplane-agent/exporter/chronicleexporter/protos/api"
	"github.com/observiq/bindplane-agent/exporter/chronicleexporter/protos/api/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestLogsDataPusher(t *testing.T) {

	// Set up configuration, logger, and context
	cfg := Config{Endpoint: baseEndpoint}
	ctx := context.Background()

	testCases := []struct {
		desc          string
		setupExporter func() *chronicleExporter
		setupMocks    func(*mocks.MockIngestionServiceV2Client)
		expectedErr   string
		permanentErr  bool
	}{
		{
			desc: "successful push to Chronicle",
			setupExporter: func() *chronicleExporter {
				mockClient := mocks.NewMockIngestionServiceV2Client(gomock.NewController(t))
				marshaller := NewMockMarshaler(t)
				marshaller.On("MarshalRawLogs", mock.Anything, mock.Anything).Return([]*api.BatchCreateLogsRequest{{}}, nil)
				return &chronicleExporter{
					cfg:        &cfg,
					metrics:    newExporterMetrics([]byte{}, []byte{}, "", cfg.Namespace),
					logger:     zap.NewNop(),
					grpcClient: mockClient,
					marshaler:  marshaller,
				}
			},
			setupMocks: func(mockClient *mocks.MockIngestionServiceV2Client) {
				mockClient.EXPECT().BatchCreateLogs(gomock.Any(), gomock.Any(), gomock.Any()).Return(&api.BatchCreateLogsResponse{}, nil)
			},
			expectedErr: "",
		},
		{
			desc: "upload to Chronicle fails (transient)",
			setupExporter: func() *chronicleExporter {
				mockClient := mocks.NewMockIngestionServiceV2Client(gomock.NewController(t))
				marshaller := NewMockMarshaler(t)
				marshaller.On("MarshalRawLogs", mock.Anything, mock.Anything).Return([]*api.BatchCreateLogsRequest{{}}, nil)
				return &chronicleExporter{
					cfg:        &cfg,
					metrics:    newExporterMetrics([]byte{}, []byte{}, "", cfg.Namespace),
					logger:     zap.NewNop(),
					grpcClient: mockClient,
					marshaler:  marshaller,
				}
			},
			setupMocks: func(mockClient *mocks.MockIngestionServiceV2Client) {
				// Simulate an error returned from the Chronicle service
				err := status.Error(codes.Unavailable, "service unavailable")
				mockClient.EXPECT().BatchCreateLogs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, err)
			},
			expectedErr: "service unavailable",
		},
		{
			desc: "upload to Chronicle fails (permanent)",
			setupExporter: func() *chronicleExporter {
				mockClient := mocks.NewMockIngestionServiceV2Client(gomock.NewController(t))
				marshaller := NewMockMarshaler(t)
				marshaller.On("MarshalRawLogs", mock.Anything, mock.Anything).Return([]*api.BatchCreateLogsRequest{{}}, nil)
				return &chronicleExporter{
					cfg:        &cfg,
					metrics:    newExporterMetrics([]byte{}, []byte{}, "", cfg.Namespace),
					logger:     zap.NewNop(),
					grpcClient: mockClient,
					marshaler:  marshaller,
				}
			},
			setupMocks: func(mockClient *mocks.MockIngestionServiceV2Client) {
				err := status.Error(codes.InvalidArgument, "Invalid argument detected.")
				mockClient.EXPECT().BatchCreateLogs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, err)
			},
			expectedErr:  "Invalid argument detected.",
			permanentErr: true,
		},
		{
			desc: "marshaler error",
			setupExporter: func() *chronicleExporter {
				mockClient := mocks.NewMockIngestionServiceV2Client(gomock.NewController(t))
				marshaller := NewMockMarshaler(t)
				// Simulate an error during log marshaling
				marshaller.On("MarshalRawLogs", mock.Anything, mock.Anything).Return(nil, errors.New("marshal error"))
				return &chronicleExporter{
					cfg:        &cfg,
					metrics:    newExporterMetrics([]byte{}, []byte{}, "", cfg.Namespace),
					logger:     zap.NewNop(),
					grpcClient: mockClient,
					marshaler:  marshaller,
				}
			},
			setupMocks: func(_ *mocks.MockIngestionServiceV2Client) {
				// No need to setup mocks for the client as the error occurs before the client is used
			},
			expectedErr: "marshal error",
		},
		{
			desc: "empty log records",
			setupExporter: func() *chronicleExporter {
				mockClient := mocks.NewMockIngestionServiceV2Client(gomock.NewController(t))
				marshaller := NewMockMarshaler(t)
				// Return an empty slice to simulate no logs to push
				marshaller.On("MarshalRawLogs", mock.Anything, mock.Anything).Return([]*api.BatchCreateLogsRequest{}, nil)
				return &chronicleExporter{
					cfg:        &cfg,
					metrics:    newExporterMetrics([]byte{}, []byte{}, "", cfg.Namespace),
					logger:     zap.NewNop(),
					grpcClient: mockClient,
					marshaler:  marshaller,
				}
			},
			setupMocks: func(_ *mocks.MockIngestionServiceV2Client) {
				// Expect no calls to BatchCreateLogs since there are no logs to push
			},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			exporter := tc.setupExporter()
			tc.setupMocks(exporter.grpcClient.(*mocks.MockIngestionServiceV2Client))

			// Create a dummy plog.Logs to pass to logsDataPusher
			logs := mockLogs(mockLogRecord("Test body", map[string]any{"key1": "value1"}))

			err := exporter.logsDataPusher(ctx, logs)

			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
				if tc.permanentErr {
					require.True(t, consumererror.IsPermanent(err), "Expected error to be permanent")
				} else {
					require.False(t, consumererror.IsPermanent(err), "Expected error to be transient")
				}
			}
		})
	}
}

// Conclusion: the HTTP endpoint may have a size limit, but will truncate any logs beyond the 1000th in the payload
// Secondary conclusion: the HTTP endpoint likely has a 14MB limit for our SecOps license
func TestReallyPushToChronicleHTTP(t *testing.T) {
	f := NewFactory()
	cfg := f.CreateDefaultConfig().(*Config)

	cfg.CustomerID = "b536658e-469e-44a5-b764-d5ab15b72ce0"
	cfg.Protocol = protocolHTTPS
	cfg.CredsFilePath = "/etc/otel/telemetry-sandbox-http.json"
	cfg.Project = "telemetry-sandbox-340915"
	cfg.Forwarder = "fbcb7d35-dd6f-4f20-aa2e-5ab2cceffa3c"
	cfg.Location = "us"
	cfg.Endpoint = "chronicle.googleapis.com"
	cfg.LogType = "OFFICE_365"
	cfg.TimeoutConfig.Timeout = 5 * time.Minute
	settings := exportertest.NewNopSettings()

	exp, err := f.CreateLogs(context.Background(), settings, cfg)
	require.NoError(t, err)

	require.NoError(t, exp.Start(context.Background(), componenttest.NewNopHost()))

	logs := generateFakeLogs(932)

	require.NoError(t, exp.ConsumeLogs(context.Background(), logs))

	// time.Sleep(5 * time.Second)

	require.NoError(t, exp.Shutdown(context.Background()))
}

// Conclusion: the GRPC endpoint may have a size limit, but will truncate any logs beyond the 1000th in the payload
func TestReallyPushToChronicleGRPC(t *testing.T) {
	f := NewFactory()
	cfg := f.CreateDefaultConfig().(*Config)

	cfg.CustomerID = "b536658e-469e-44a5-b764-d5ab15b72ce0"
	cfg.Protocol = protocolGRPC
	cfg.CredsFilePath = "/etc/otel/telemetry-sandbox-grpc.json"
	cfg.Project = "telemetry-sandbox-340915"
	cfg.Forwarder = "fbcb7d35-dd6f-4f20-aa2e-5ab2cceffa3c"
	cfg.Location = "us"
	cfg.Endpoint = "malachiteingestion-pa.googleapis.com"
	cfg.LogType = "OFFICE_365"
	cfg.TimeoutConfig.Timeout = 5 * time.Minute
	settings := exportertest.NewNopSettings()

	exp, err := f.CreateLogs(context.Background(), settings, cfg)
	require.NoError(t, err)

	require.NoError(t, exp.Start(context.Background(), componenttest.NewNopHost()))

	logs := generateFakeLogs(470)

	require.NoError(t, exp.ConsumeLogs(context.Background(), logs))

	// time.Sleep(10 * time.Second)

	require.NoError(t, exp.Shutdown(context.Background()))

	// time.Sleep(10 * time.Second)
}

func generateFakeLogs(numLogs int) plog.Logs {
	logs := plog.NewLogs()
	// Configure resource attributes
	res := logs.ResourceLogs().AppendEmpty()
	res.Resource().Attributes().PutStr("service.name", "google-secops-log-generator")
	res.Resource().Attributes().PutStr("log_type", "OFFICE_365")

	// this is allegedly a 16000 byte string
	body, err := os.ReadFile(filepath.Join("testdata", "longstring.txt"))
	if err != nil {
		panic(err)
	}

	// Add log entry
	scopeLogs := res.ScopeLogs().AppendEmpty()
	for i := 0; i < numLogs; i++ {
		log := scopeLogs.LogRecords().AppendEmpty()
		log.SetSeverityNumber(plog.SeverityNumberInfo)
		log.Body().SetStr(string(body))
		log.Attributes().PutStr("operation_name", "UserLogin")
		log.Attributes().PutStr("user_email", "user@example.com")
		log.Attributes().PutStr("status", "Success")
		log.Attributes().PutStr("ip_address", "203.0.113.42")
		log.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	}

	return logs
}

func TestRequestSizeHTTP(t *testing.T) {
	f := NewFactory()
	cfg := f.CreateDefaultConfig().(*Config)

	cfg.CustomerID = "b536658e-469e-44a5-b764-d5ab15b72ce0"
	cfg.Protocol = protocolHTTPS
	cfg.CredsFilePath = "/etc/otel/telemetry-sandbox-http.json"
	cfg.Project = "telemetry-sandbox-340915"
	cfg.Forwarder = "fbcb7d35-dd6f-4f20-aa2e-5ab2cceffa3c"
	cfg.Location = "us"
	cfg.Endpoint = "chronicle.googleapis.com"
	cfg.LogType = "OFFICE_365"
	settings := exportertest.NewNopSettings()

	marshaler, err := newProtoMarshaler(*cfg, settings.TelemetrySettings, []byte(cfg.CustomerID))
	require.NoError(t, err)

	testCases := []struct {
		desc     string
		entries  *api.ImportLogsRequest
		expected int
	}{
		{
			desc: "no entries",
			entries: func() *api.ImportLogsRequest {
				logs := generateFakeLogs(933)
				stuff, err := marshaler.MarshalRawLogsForHTTP(context.Background(), logs)
				require.NoError(t, err)
				return stuff["OFFICE_365"][0]
			}(),
			expected: 102,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// 10432025 - 932
			// 10443218 - 933
			size := proto.Size(tc.entries)

			// 7496680 - 500
			// 11244933 - 750
			// 13119052 - 875
			// 13643810 - 910
			// 13868708 - 925 S
			// 13928677 - 929 S
			// 13958666 - 931 S
			// 13973653 - 932 S
			// 13988643 - 933 F
			// 14003645 - 934 F
			// 14093594 - 940 F
			// 14993180 - 1000
			// body, err := protojson.Marshal(tc.entries)
			// require.NoError(t, err)
			// size := len(body)

			require.Equal(t, tc.expected, size)
		})
	}
}

func BenchmarkRequestSize(b *testing.B) {
	f := NewFactory()
	cfg := f.CreateDefaultConfig().(*Config)

	cfg.CustomerID = "b536658e-469e-44a5-b764-d5ab15b72ce0"
	cfg.Protocol = protocolHTTPS
	cfg.CredsFilePath = "/etc/otel/telemetry-sandbox-http.json"
	cfg.Project = "telemetry-sandbox-340915"
	cfg.Forwarder = "fbcb7d35-dd6f-4f20-aa2e-5ab2cceffa3c"
	cfg.Location = "us"
	cfg.Endpoint = "chronicle.googleapis.com"
	cfg.LogType = "OFFICE_365"
	settings := exportertest.NewNopSettings()

	marshaler, err := newProtoMarshaler(*cfg, settings.TelemetrySettings, []byte(cfg.CustomerID))
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		logs := generateFakeLogs(i)
		stuff, err := marshaler.MarshalRawLogsForHTTP(context.Background(), logs)
		require.NoError(b, err)

		x := proto.Size(stuff["OFFICE_365"][0])
		// body, err := protojson.Marshal(stuff["OFFICE_365"])
		// require.NoError(b, err)
		// x := len(body)
		_ = x
	}
}

func TestRequestSizeGRPC(t *testing.T) {
	f := NewFactory()
	cfg := f.CreateDefaultConfig().(*Config)

	cfg.CustomerID = "b536658e-469e-44a5-b764-d5ab15b72ce0"
	cfg.Protocol = protocolGRPC
	cfg.CredsFilePath = "/etc/otel/telemetry-sandbox-grpc.json"
	cfg.Project = "telemetry-sandbox-340915"
	cfg.Forwarder = "fbcb7d35-dd6f-4f20-aa2e-5ab2cceffa3c"
	cfg.Location = "us"
	cfg.Endpoint = "malachiteingestion-pa.googleapis.com"
	cfg.LogType = "OFFICE_365"
	settings := exportertest.NewNopSettings()

	marshaler, err := newProtoMarshaler(*cfg, settings.TelemetrySettings, []byte(cfg.CustomerID))
	require.NoError(t, err)

	testCases := []struct {
		desc     string
		entries  *api.BatchCreateLogsRequest
		expected int
	}{
		{
			desc: "no entries",
			entries: func() *api.BatchCreateLogsRequest {
				// 100 - 1119287 - S
				// 300 - 3357989 - S
				// 400 - 4477289 - S
				// 450 - 5036939 - S
				// 468 - 5237944 - F
				//       5242880
				// 469 - 5249136 - F
				// 470 - 5260328
				// 475 - 5316764
				// 500 - 5596589 - F
				logs := generateFakeLogs(469)
				a, err := marshaler.MarshalRawLogs(context.Background(), logs)
				require.NoError(t, err)
				return a[0]
			}(),
			expected: 102,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			size := proto.Size(tc.entries)
			require.Equal(t, tc.expected, size)
		})
	}
}
