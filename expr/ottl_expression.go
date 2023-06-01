package expr

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

type OTTLExpression[T any] struct {
	statement *ottl.Statement[T]
}

func (e OTTLExpression[T]) Execute(ctx context.Context, tCtx T) (any, error) {
	val, _, err := e.statement.Execute(ctx, tCtx)
	return val, err
}
