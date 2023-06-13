package expr

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"go.opentelemetry.io/collector/component"
)

// OTTLCondition evaluates an OTTL expression as a boolean value.
type OTTLCondition[T any] struct {
	statement *ottl.Statement[T]
}

// Match returns true if the expression is true for the given transform context.
func (e OTTLCondition[T]) Match(ctx context.Context, tCtx T) (bool, error) {
	_, ran, err := e.statement.Execute(ctx, tCtx)
	return ran, err
}

// NewOTTLSpanCondition creates a new OTTLCondition for a span with the given condition.
func NewOTTLSpanCondition(condition string, set component.TelemetrySettings) (*OTTLCondition[ottlspan.TransformContext], error) {
	statementStr := "noop() where " + condition
	statement, err := NewOTTLSpanStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLCondition[ottlspan.TransformContext]{
		statement: statement,
	}, nil
}

// NewOTTLDatapointCondition creates a new OTTLCondition for a datapoint with the given condition.
func NewOTTLDatapointCondition(condition string, set component.TelemetrySettings) (*OTTLCondition[ottldatapoint.TransformContext], error) {
	statementStr := "noop() where " + condition
	statement, err := NewOTTLDatapointStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLCondition[ottldatapoint.TransformContext]{
		statement: statement,
	}, nil
}

// NewOTTLLogRecordCondition creates a new OTTLCondition for a log record with the given condition.
func NewOTTLLogRecordCondition(condition string, set component.TelemetrySettings) (*OTTLCondition[ottllog.TransformContext], error) {
	statementStr := "noop() where " + condition
	statement, err := NewOTTLLogRecordStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLCondition[ottllog.TransformContext]{
		statement: statement,
	}, nil
}
