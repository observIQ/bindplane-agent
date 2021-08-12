package logging

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCollectorLoggingOpts(t *testing.T) {
	opts := GetCollectorLoggingOpts(DefaultConfig())
	require.Len(t, opts, 1)

	opts = GetCollectorLoggingOpts(nil)
	require.Len(t, opts, 0)
}

func TestGetManagerLogger(t *testing.T) {
	// Basic test just tests logger doesn't panic on create
	_ = GetManagerLogger(DefaultConfig())
	_ = GetManagerLogger(nil)
}
