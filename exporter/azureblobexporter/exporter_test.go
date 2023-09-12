package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"context"
	"errors"
	"testing"

	"github.com/observiq/bindplane-agent/exporter/azureblobexporter/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

func Test_exporter_Capabilities(t *testing.T) {
	exp := &azureBlobExporter{}
	cap := exp.Capabilities()
	require.False(t, cap.MutatesData)
}

func Test_exporter_metricsDataPusher(t *testing.T) {
	cfg := &Config{
		Container:  "container",
		BlobPrefix: "blob",
		RootFolder: "root",
		Partition:  minutePartition,
	}

	testCases := []struct {
		desc        string
		mockGen     func(t *testing.T, input pmetric.Metrics, expectBuff []byte) (blobClient, marshaler)
		expectedErr error
	}{
		{
			desc: "marshal error",
			mockGen: func(t *testing.T, input pmetric.Metrics, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalMetrics(input).Return(nil, errors.New("marshal"))

				return mockBlobClient, mockMarshaler
			},
			expectedErr: errors.New("marshal"),
		},
		{
			desc: "Blob client error",
			mockGen: func(t *testing.T, input pmetric.Metrics, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalMetrics(input).Return(expectBuff, nil)
				mockMarshaler.EXPECT().Format().Return("json")

				mockBlobClient.EXPECT().UploadBuffer(mock.Anything, cfg.Container, mock.Anything, expectBuff).Return(errors.New("client"))

				return mockBlobClient, mockMarshaler
			},
			expectedErr: errors.New("client"),
		},
		{
			desc: "Successful push",
			mockGen: func(t *testing.T, input pmetric.Metrics, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalMetrics(input).Return(expectBuff, nil)
				mockMarshaler.EXPECT().Format().Return("json")

				mockBlobClient.EXPECT().UploadBuffer(mock.Anything, cfg.Container, mock.Anything, expectBuff).Return(nil)

				return mockBlobClient, mockMarshaler
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			md, expectedBytes := generateTestMetrics(t)
			mockBlobClient, mockMarshaler := tc.mockGen(t, md, expectedBytes)
			exp := &azureBlobExporter{
				cfg:        cfg,
				blobClient: mockBlobClient,
				logger:     zap.NewNop(),
				marshaler:  mockMarshaler,
			}

			err := exp.metricsDataPusher(context.Background(), md)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

		})
	}
}

func Test_exporter_logsDataPusher(t *testing.T) {
	cfg := &Config{
		Container:  "container",
		BlobPrefix: "blob",
		RootFolder: "root",
		Partition:  minutePartition,
	}

	testCases := []struct {
		desc        string
		mockGen     func(t *testing.T, input plog.Logs, expectBuff []byte) (blobClient, marshaler)
		expectedErr error
	}{
		{
			desc: "marshal error",
			mockGen: func(t *testing.T, input plog.Logs, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalLogs(input).Return(nil, errors.New("marshal"))

				return mockBlobClient, mockMarshaler
			},
			expectedErr: errors.New("marshal"),
		},
		{
			desc: "Blob client error",
			mockGen: func(t *testing.T, input plog.Logs, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalLogs(input).Return(expectBuff, nil)
				mockMarshaler.EXPECT().Format().Return("json")

				mockBlobClient.EXPECT().UploadBuffer(mock.Anything, cfg.Container, mock.Anything, expectBuff).Return(errors.New("client"))

				return mockBlobClient, mockMarshaler
			},
			expectedErr: errors.New("client"),
		},
		{
			desc: "Successful push",
			mockGen: func(t *testing.T, input plog.Logs, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalLogs(input).Return(expectBuff, nil)
				mockMarshaler.EXPECT().Format().Return("json")

				mockBlobClient.EXPECT().UploadBuffer(mock.Anything, cfg.Container, mock.Anything, expectBuff).Return(nil)

				return mockBlobClient, mockMarshaler
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ld, expectedBytes := generateTestLogs(t)
			mockBlobClient, mockMarshaler := tc.mockGen(t, ld, expectedBytes)
			exp := &azureBlobExporter{
				cfg:        cfg,
				blobClient: mockBlobClient,
				logger:     zap.NewNop(),
				marshaler:  mockMarshaler,
			}

			err := exp.logsDataPusher(context.Background(), ld)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

		})
	}
}

func Test_exporter_traceDataPusher(t *testing.T) {
	cfg := &Config{
		Container:  "container",
		BlobPrefix: "blob",
		RootFolder: "root",
		Partition:  minutePartition,
	}

	testCases := []struct {
		desc        string
		mockGen     func(t *testing.T, input ptrace.Traces, expectBuff []byte) (blobClient, marshaler)
		expectedErr error
	}{
		{
			desc: "marshal error",
			mockGen: func(t *testing.T, input ptrace.Traces, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalTraces(input).Return(nil, errors.New("marshal"))

				return mockBlobClient, mockMarshaler
			},
			expectedErr: errors.New("marshal"),
		},
		{
			desc: "Blob client error",
			mockGen: func(t *testing.T, input ptrace.Traces, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalTraces(input).Return(expectBuff, nil)
				mockMarshaler.EXPECT().Format().Return("json")

				mockBlobClient.EXPECT().UploadBuffer(mock.Anything, cfg.Container, mock.Anything, expectBuff).Return(errors.New("client"))

				return mockBlobClient, mockMarshaler
			},
			expectedErr: errors.New("client"),
		},
		{
			desc: "Successful push",
			mockGen: func(t *testing.T, input ptrace.Traces, expectBuff []byte) (blobClient, marshaler) {
				mockBlobClient := mocks.NewMockBlobClient(t)
				mockMarshaler := mocks.NewMockMarshaler(t)

				mockMarshaler.EXPECT().MarshalTraces(input).Return(expectBuff, nil)
				mockMarshaler.EXPECT().Format().Return("json")

				mockBlobClient.EXPECT().UploadBuffer(mock.Anything, cfg.Container, mock.Anything, expectBuff).Return(nil)

				return mockBlobClient, mockMarshaler
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			td, expectedBytes := generateTestTraces(t)
			mockBlobClient, mockMarshaler := tc.mockGen(t, td, expectedBytes)
			exp := &azureBlobExporter{
				cfg:        cfg,
				blobClient: mockBlobClient,
				logger:     zap.NewNop(),
				marshaler:  mockMarshaler,
			}

			err := exp.traceDataPusher(context.Background(), td)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

		})
	}
}
