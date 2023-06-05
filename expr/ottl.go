package expr

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
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

func NewOTTLLogRecordStatement(statementStr string, set component.TelemetrySettings) (*ottl.Statement[ottllog.TransformContext], error) {
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
	return ottl.CreateFactoryMap[T](
		ottlfuncs.NewConcatFactory[T](),
		ottlfuncs.NewConvertCaseFactory[T](),
		ottlfuncs.NewIntFactory[T](),
		ottlfuncs.NewIsMatchFactory[T](),
		ottlfuncs.NewLogFactory[T](),
		ottlfuncs.NewParseJSONFactory[T](),
		ottlfuncs.NewSpanIDFactory[T](),
		ottlfuncs.NewSplitFactory[T](),
		ottlfuncs.NewSubstringFactory[T](),
		ottlfuncs.NewTraceIDFactory[T](),
		ottlfuncs.NewUUIDFactory[T](),

		newNoopFactory[T](),
		newValueFactory[T](),
	)
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

type valueArguments[K any] struct {
	Target ottl.Getter[K] `ottlarg:"0"`
}

func newValueFactory[K any]() ottl.Factory[K] {
	return ottl.NewFactory("value", &valueArguments[K]{}, createValueFunction[K])
}

func createValueFunction[K any](c ottl.FunctionContext, a ottl.Arguments) (ottl.ExprFunc[K], error) {
	args, ok := a.(*valueArguments[K])
	if !ok {
		return nil, fmt.Errorf("valueFactory args must be of type *valueArguments[K]")
	}

	return valueFn[K](args)
}

func valueFn[K any](c *valueArguments[K]) (ottl.ExprFunc[K], error) {
	return func(ctx context.Context, tCtx K) (interface{}, error) {
		return c.Target.Get(ctx, tCtx)
	}, nil
}
