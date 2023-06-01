package expr

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlmetric"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"go.opentelemetry.io/collector/component"
)

type OTTLCondition[T any] struct {
	statement *ottl.Statement[T]
}

func (e OTTLCondition[T]) Match(ctx context.Context, tCtx T) (bool, error) {
	_, ran, err := e.statement.Execute(ctx, tCtx)
	return ran, err
}

func NewOTTLSpanCondition(condition string, set component.TelemetrySettings) (*OTTLCondition[ottlspan.TransformContext], error) {
	statementStr := "drop() where " + condition
	statement, err := NewOTTLSpanStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLCondition[ottlspan.TransformContext]{
		statement: statement,
	}, nil
}

func NewOTTLMetricCondition(condition string, set component.TelemetrySettings) (*OTTLCondition[ottlmetric.TransformContext], error) {
	statementStr := "drop() where " + condition
	statement, err := NewOTTLMetricStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLCondition[ottlmetric.TransformContext]{
		statement: statement,
	}, nil
}

func NewOTTLDatapointCondition(condition string, set component.TelemetrySettings) (*OTTLCondition[ottldatapoint.TransformContext], error) {
	statementStr := "drop() where " + condition
	statement, err := NewOTTLDatapointStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLCondition[ottldatapoint.TransformContext]{
		statement: statement,
	}, nil
}

func NewOTTLLogCondition(condition string, set component.TelemetrySettings) (*OTTLCondition[ottllog.TransformContext], error) {
	statementStr := "drop() where " + condition
	statement, err := NewOTTLLogStatement(statementStr, set)
	if err != nil {
		return nil, err
	}

	return &OTTLCondition[ottllog.TransformContext]{
		statement: statement,
	}, nil
}
