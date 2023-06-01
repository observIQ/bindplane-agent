package expr

import (
	"context"

	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlmetric"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
)

func NewOTTLSpanStatement(statementStr string, set component.TelemetrySettings) (*ottl.Statement[ottlspan.TransformContext], error) {
	parser, err := ottlspan.NewParser(functions[ottlspan.TransformContext](), set)
	if err != nil {
		return nil, err
	}

	statement, err := parser.ParseStatement(statementStr)
	if err != nil {
		return nil, err
	}

	return statement, nil
}

func NewOTTLMetricStatement(statementStr string, set component.TelemetrySettings) (*ottl.Statement[ottlmetric.TransformContext], error) {
	parser, err := ottlmetric.NewParser(functions[ottlmetric.TransformContext](), set)
	if err != nil {
		return nil, err
	}
	statement, err := parser.ParseStatement(statementStr)
	if err != nil {
		return nil, err
	}

	return statement, nil
}

func NewOTTLDatapointStatement(statementStr string, set component.TelemetrySettings) (*ottl.Statement[ottldatapoint.TransformContext], error) {
	parser, err := ottldatapoint.NewParser(functions[ottldatapoint.TransformContext](), set)
	if err != nil {
		return nil, err
	}
	statement, err := parser.ParseStatement(statementStr)
	if err != nil {
		return nil, err
	}

	return statement, nil
}

func NewOTTLLogStatement(statementStr string, set component.TelemetrySettings) (*ottl.Statement[ottllog.TransformContext], error) {
	parser, err := ottllog.NewParser(functions[ottllog.TransformContext](), set)
	if err != nil {
		return nil, err
	}

	statement, err := parser.ParseStatement(statementStr)
	if err != nil {
		return nil, err
	}

	return statement, nil
}

func functions[T any]() map[string]ottl.Factory[T] {
	return map[string]ottl.Factory[T]{
		"drop": newNoopFactory[T](),
	}
}

func newNoopFactory[K any]() ottl.Factory[K] {
	return ottl.NewFactory("noop", nil, createNoopFunction[K])
}

func createNoopFunction[K any](_ ottl.FunctionContext, _ ottl.Arguments) (ottl.ExprFunc[K], error) {
	return noopFn[K]()
}

func noopFn[K any]() (ottl.ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return true, nil
	}, nil
}
