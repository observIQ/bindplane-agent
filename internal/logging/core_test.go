package logging

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateFileCore(t *testing.T) {
	// Test that CreateFileCore doesn't panic for default config
	config, err := DefaultConfig()
	require.NoError(t, err)

	_ = CreateFileCore(&config.Collector)
	_ = CreateFileCore(&config.Manager)
}
