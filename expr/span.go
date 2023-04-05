package expr

import (
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	SpanKindField          = "trace_kind"
	SpanStatusCodeField    = "trace_status_code"
	SpanStatusMessageField = "trace_status_message"
	SpanDurationField      = "span_duration_ms"
)

var spanKindToString = map[ptrace.SpanKind]string{
	ptrace.SpanKindUnspecified: "unspecified",
	ptrace.SpanKindInternal:    "internal",
	ptrace.SpanKindClient:      "client",
	ptrace.SpanKindServer:      "server",
	ptrace.SpanKindConsumer:    "consumer",
	ptrace.SpanKindProducer:    "producer",
}

var spanStatusCodeToString = map[ptrace.StatusCode]string{
	ptrace.StatusCodeError: "error",
	ptrace.StatusCodeOk:    "ok",
	ptrace.StatusCodeUnset: "unset",
}

type Span = map[string]any

func convertToSpan(span ptrace.Span, resource map[string]any) Span {
	return Span{
		ResourceField:          resource,
		AttributesField:        span.Attributes().AsRaw(),
		SpanDurationField:      span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()).Milliseconds(),
		SpanKindField:          spanKindToString[span.Kind()],
		SpanStatusCodeField:    spanStatusCodeToString[span.Status().Code()],
		SpanStatusMessageField: span.Status().Message(),
	}
}

type SpanResourceGroup struct {
	Resource map[string]any
	Spans    []Span
}

func ConvertToSpanResourceGroups(logs ptrace.Traces) []SpanResourceGroup {
	groups := make([]SpanResourceGroup, 0, logs.ResourceSpans().Len())

	for i := 0; i < logs.ResourceSpans().Len(); i++ {
		resourceLogs := logs.ResourceSpans().At(i)
		resource := resourceLogs.Resource().Attributes().AsRaw()
		group := SpanResourceGroup{
			Resource: resource,
			Spans:    make([]Span, 0, resourceLogs.ScopeSpans().Len()),
		}
		for j := 0; j < resourceLogs.ScopeSpans().Len(); j++ {
			logs := resourceLogs.ScopeSpans().At(j).Spans()
			for k := 0; k < logs.Len(); k++ {
				log := logs.At(k)
				group.Spans = append(group.Spans, convertToSpan(log, resource))
			}
		}
		groups = append(groups, group)
	}

	return groups
}
