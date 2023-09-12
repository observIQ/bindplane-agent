package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"bytes"
	"compress/gzip"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func Test_baseMarshaler(t *testing.T) {
	m := newMarshaler(noCompression)

	require.IsType(t, &baseMarshaler{}, m)
	require.Equal(t, "json", m.format())

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
	require.Equal(t, "json.gz", m.format())

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

func generateTestMetrics(t *testing.T) (md pmetric.Metrics, expectedBytes []byte) {
	t.Helper()

	md = pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName("metric")
	gm := m.SetEmptyGauge()
	dp := gm.DataPoints().AppendEmpty()
	dp.SetIntValue(1)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	marshaler := &pmetric.JSONMarshaler{}

	var err error
	expectedBytes, err = marshaler.MarshalMetrics(md)
	require.NoError(t, err)

	return md, expectedBytes
}

func generateTestLogs(t *testing.T) (ld plog.Logs, expectedBytes []byte) {
	t.Helper()

	ld = plog.NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	l := sl.LogRecords().AppendEmpty()
	l.Body().SetStr("body")
	l.Attributes().PutBool("bool", true)
	l.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	marshaler := &plog.JSONMarshaler{}

	var err error
	expectedBytes, err = marshaler.MarshalLogs(ld)
	require.NoError(t, err)

	return ld, expectedBytes
}

func generateTestTraces(t *testing.T) (td ptrace.Traces, expectedBytes []byte) {
	t.Helper()

	td = ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	ss := rs.ScopeSpans().AppendEmpty()
	s := ss.Spans().AppendEmpty()
	s.Attributes().PutBool("bool", true)
	s.SetName("span")
	s.SetStartTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	marshaler := &ptrace.JSONMarshaler{}

	var err error
	expectedBytes, err = marshaler.MarshalTraces(td)
	require.NoError(t, err)

	return td, expectedBytes
}

func verifyGZipCompress(t *testing.T, input, actual []byte) {
	t.Helper()
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(input)
	require.NoError(t, err)

	require.Equal(t, buf.Bytes(), actual)
}
