package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_baseMarshaler(t *testing.T) {
	m := newMarshaler(noCompression)

	require.IsType(t, &baseMarshaler{}, m)
	require.Equal(t, "json", m.Format())

	t.Run("Metrics Marshal", func(t *testing.T) {
		t.Parallel()
		md, expectedBytes := generateTestMetrics(t)
		actualBytes, err := m.MarshalMetrics(md)
		require.NoError(t, err)

		require.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("Logs Marshal", func(t *testing.T) {
		t.Parallel()
		ld, expectedBytes := generateTestLogs(t)
		actualBytes, err := m.MarshalLogs(ld)
		require.NoError(t, err)

		require.Equal(t, expectedBytes, actualBytes)
	})

	t.Run("Trace Marshal", func(t *testing.T) {
		t.Parallel()
		td, expectedBytes := generateTestTraces(t)
		actualBytes, err := m.MarshalTraces(td)
		require.NoError(t, err)

		require.Equal(t, expectedBytes, actualBytes)
	})
}

func Test_gzipMarshaler(t *testing.T) {
	m := newMarshaler(gzipCompression)

	require.IsType(t, &gzipMarshaler{}, m)
	require.Equal(t, "json.gz", m.Format())

	t.Run("Metrics Marshal", func(t *testing.T) {
		t.Parallel()
		md, inputBytes := generateTestMetrics(t)
		actualBytes, err := m.MarshalMetrics(md)
		require.NoError(t, err)

		verifyGZipCompress(t, inputBytes, actualBytes)
	})

	t.Run("Logs Marshal", func(t *testing.T) {
		t.Parallel()
		ld, inputBytes := generateTestLogs(t)
		actualBytes, err := m.MarshalLogs(ld)
		require.NoError(t, err)

		verifyGZipCompress(t, inputBytes, actualBytes)
	})

	t.Run("Trace Marshal", func(t *testing.T) {
		t.Parallel()
		td, inputBytes := generateTestTraces(t)
		actualBytes, err := m.MarshalTraces(td)
		require.NoError(t, err)

		verifyGZipCompress(t, inputBytes, actualBytes)
	})
}

func verifyGZipCompress(t *testing.T, input, actual []byte) {
	t.Helper()
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(input)
	require.NoError(t, err)

	require.Equal(t, buf.Bytes(), actual)
}
