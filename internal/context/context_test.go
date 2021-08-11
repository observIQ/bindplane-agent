package context

import (
	"context"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyContext(t *testing.T) {
	ctx := EmptyContext()
	require.Equal(t, context.Background(), ctx)
}

func TestWithParent(t *testing.T) {
	ctx := WithParent(-5)
	<-ctx.Done()
}

func TestWithInterruptSignal(t *testing.T) {
	ctx, cancel := WithInterrupt(context.Background())
	defer cancel()

	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	require.NoError(t, err)

	_, ok := <-ctx.Done()
	require.False(t, ok)
}
