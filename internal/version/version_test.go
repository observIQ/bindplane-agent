package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaults(t *testing.T) {
	require.Equal(t, version, Version())
	require.Equal(t, gitHash, GitHash())
	require.Equal(t, date, Date())
}
