package expr

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"go.opentelemetry.io/collector/component"
)

// OTTLExpression evaluates an OTTL expression, returning it's resultant value.
type OTTLExpression[T any] struct {
	statement *ottl.Statement[T]
}

// Execute executes the expression with the given context, returning the value of the expression.
func (e OTTLExpression[T]) Execute(ctx context.Context, tCtx T) (any, error) {
	val, _, err := e.statement.Execute(ctx, tCtx)
	return val, err
}

// NewOTTLSpanExpression creates a new expression for spans.
// The expression is wrapped in an editor function, so only Converter functions and target expressions can be used.
func NewOTTLSpanExpression(expression string, set component.TelemetrySettings) (*OTTLExpression[ottlspan.TransformContext], error) {
	statementStr := fmt.Sprintf("value(%s) where 1==1", expression)
	statement, err := NewOTTLSpanStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLExpression[ottlspan.TransformContext]{
		statement: statement,
	}, nil
}

// NewOTTLDatapointExpression creates a new expression for datapoints.
// The expression is wrapped in an editor function, so only Converter functions and target expressions can be used.
func NewOTTLDatapointExpression(expression string, set component.TelemetrySettings) (*OTTLExpression[ottldatapoint.TransformContext], error) {
	statementStr := fmt.Sprintf("value(%s) where 1==1", expression)
	statement, err := NewOTTLDatapointStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLExpression[ottldatapoint.TransformContext]{
		statement: statement,
	}, nil
}

// NewOTTLLogRecordExpression creates a new expression for log records.
// The expression is wrapped in an editor function, so only Converter functions and target expressions can be used.
func NewOTTLLogRecordExpression(expression string, set component.TelemetrySettings) (*OTTLExpression[ottllog.TransformContext], error) {
	statementStr := fmt.Sprintf("value(%s) where 1==1", expression)
	statement, err := NewOTTLLogRecordStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLExpression[ottllog.TransformContext]{
		statement: statement,
	}, nil
}
