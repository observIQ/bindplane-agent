package cloudwatch

import (
	"context"
	"testing"

	"github.com/open-telemetry/opentelemetry-log-collection/operator"
	"github.com/open-telemetry/opentelemetry-log-collection/testutil"
	"github.com/stretchr/testify/require"
)

func TestPersisterCache(t *testing.T) {
	ctx := context.TODO()
	stubDatabase := testutil.NewMockPersister("stub")
	p := persister{
		DB: operator.NewScopedPersister("test", stubDatabase),
	}
	err := p.Write(ctx, "key", int64(1620666055012))
	require.NoError(t, err)

	value, readErr := p.Read(ctx, "key")
	require.NoError(t, readErr)
	require.Equal(t, int64(1620666055012), value)
}

func TestPersisterLoad(t *testing.T) {
	ctx := context.TODO()
	p := testutil.NewMockPersister("mock")
	cwP := persister{
		DB: p,
	}
	err := cwP.Write(ctx, "key", 1620666055012)
	require.NoError(t, err)

	value, syncErr := cwP.Read(ctx, "key")
	require.NoError(t, syncErr)
	require.Equal(t, int64(1620666055012), value)
}

func TestPersistentLoadNoKey(t *testing.T) {
	ctx := context.TODO()

	p := testutil.NewMockPersister("mock")
	cwP := persister{
		DB: p,
	}
	value, err := cwP.Read(ctx, "key")
	require.NoError(t, err)
	require.Equal(t, int64(0), value)
}
