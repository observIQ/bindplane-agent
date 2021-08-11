package context

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithInterruptSignal(t *testing.T) {
	ctx, cancel := WithInterrupt()
	defer cancel()

	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	require.NoError(t, err)

	_, ok := <-ctx.Done()
	require.False(t, ok)
}

func TestWithInterruptCancel(t *testing.T) {
	ctx, cancel := WithInterrupt()
	cancel()

	_, ok := <-ctx.Done()
	require.False(t, ok)
}
