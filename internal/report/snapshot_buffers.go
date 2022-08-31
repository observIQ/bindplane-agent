package report

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type logBuffer struct {
	mutex     sync.Mutex
	buffer    []plog.Logs
	idealSize int
}

func newLogBuffer(idealSize int) *logBuffer {
	return &logBuffer{
		buffer:    make([]plog.Logs, 0),
		idealSize: idealSize,
	}
}

func (l *logBuffer) Len() int {
	size := 0
	for _, ld := range l.buffer {
		size += ld.LogRecordCount()
	}

	return size
}

func (l *logBuffer) Add(ld plog.Logs) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	logSize := ld.LogRecordCount()
	bufferSize := l.Len()
	switch {
	// The number of logs is more than idealSize so reset this to just this log set
	case logSize > l.idealSize:
		l.buffer = []plog.Logs{ld}

	// Haven't reached idealSize yet so add this
	case logSize+bufferSize < l.idealSize:
		l.buffer = append(l.buffer, ld)

	// Adding this will put us over idealSize so and add the new logs.
	// Only remove the oldest if it does not bring buffer under idealSize
	case logSize+bufferSize >= l.idealSize:
		l.buffer = append(l.buffer, ld)

		// If removing the oldest item does not put us under the ideal size then it's ok to do so
		oldest := l.buffer[0]
		newBufferSize := logSize + bufferSize
		if newBufferSize-oldest.LogRecordCount() >= l.idealSize {
			l.buffer = l.buffer[1:]
		}
	}
}

func (l *logBuffer) ConstructPayload() ([]byte, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	logsMarshaler := plog.NewProtoMarshaler()

	payloadLogs := plog.NewLogs()
	for _, ld := range l.buffer {
		ld.ResourceLogs().MoveAndAppendTo(payloadLogs.ResourceLogs())
	}

	// update the buffer to retain the current logs which were moved to the new payload
	l.buffer = []plog.Logs{payloadLogs}

	payload, err := logsMarshaler.MarshalLogs(payloadLogs)
	if err != nil {
		return nil, fmt.Errorf("failed to construct payload: %w", err)
	}

	return payload, nil
}

type metricBuffer struct {
	mutex     sync.Mutex
	buffer    []pmetric.Metrics
	idealSize int
}

func newMetricBuffer(idealSize int) *metricBuffer {
	return &metricBuffer{
		buffer:    make([]pmetric.Metrics, 0),
		idealSize: idealSize,
	}
}

func (l *metricBuffer) Len() int {
	size := 0
	for _, md := range l.buffer {
		size += md.DataPointCount()
	}

	return size
}

func (l *metricBuffer) Add(md pmetric.Metrics) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	metricSize := md.DataPointCount()
	bufferSize := l.Len()
	switch {
	// The number of metrics is more than idealSize so reset this to just this metric set
	case metricSize > l.idealSize:
		l.buffer = []pmetric.Metrics{md}

	// Haven't reached idealSize yet so add this
	case metricSize+bufferSize < l.idealSize:
		l.buffer = append(l.buffer, md)

	// Adding this will put us over idealSize so and add the new metrics.
	// Only remove the oldest if it does not bring buffer under idealSize
	case metricSize+bufferSize >= l.idealSize:
		l.buffer = append(l.buffer, md)

		// If removing the oldest item does not put us under the ideal size then it's ok to do so
		oldest := l.buffer[0]
		newBufferSize := metricSize + bufferSize
		if newBufferSize-oldest.DataPointCount() >= l.idealSize {
			l.buffer = l.buffer[1:]
		}
	}
}

func (l *metricBuffer) ConstructPayload() ([]byte, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	metricMarshaler := pmetric.NewProtoMarshaler()

	payloadMetrics := pmetric.NewMetrics()
	for _, md := range l.buffer {
		md.ResourceMetrics().MoveAndAppendTo(payloadMetrics.ResourceMetrics())
	}

	// update the buffer to retain the current metrics which were moved to the new payload
	l.buffer = []pmetric.Metrics{payloadMetrics}

	payload, err := metricMarshaler.MarshalMetrics(payloadMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to construct payload: %w", err)
	}

	return payload, nil
}

type traceBuffer struct {
	mutex     sync.Mutex
	buffer    []ptrace.Traces
	idealSize int
}

func newTraceBuffer(idealSize int) *traceBuffer {
	return &traceBuffer{
		buffer:    make([]ptrace.Traces, 0),
		idealSize: idealSize,
	}
}

func (l *traceBuffer) Len() int {
	size := 0
	for _, td := range l.buffer {
		size += td.SpanCount()
	}

	return size
}

func (l *traceBuffer) Add(td ptrace.Traces) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	traceSize := td.SpanCount()
	bufferSize := l.Len()
	switch {
	// The number of traces is more than idealSize so reset this to just this trace set
	case traceSize > l.idealSize:
		l.buffer = []ptrace.Traces{td}

	// Haven't reached idealSize yet so add this
	case traceSize+bufferSize < l.idealSize:
		l.buffer = append(l.buffer, td)

	// Adding this will put us over idealSize so and add the new traces.
	// Only remove the oldest if it does not bring buffer under idealSize
	case traceSize+bufferSize >= l.idealSize:
		l.buffer = append(l.buffer, td)

		// If removing the oldest item does not put us under the ideal size then it's ok to do so
		oldest := l.buffer[0]
		newBufferSize := traceSize + bufferSize
		if newBufferSize-oldest.SpanCount() >= l.idealSize {
			l.buffer = l.buffer[1:]
		}
	}
}

func (l *traceBuffer) ConstructPayload() ([]byte, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	traceMarshaler := ptrace.NewProtoMarshaler()

	payloadTraces := ptrace.NewTraces()
	for _, md := range l.buffer {
		md.ResourceSpans().MoveAndAppendTo(payloadTraces.ResourceSpans())
	}

	// update the buffer to retain the current traces which were moved to the new payload
	l.buffer = []ptrace.Traces{payloadTraces}

	payload, err := traceMarshaler.MarshalTraces(payloadTraces)
	if err != nil {
		return nil, fmt.Errorf("failed to construct payload: %w", err)
	}

	return payload, nil
}
