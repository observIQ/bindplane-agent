package process

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchesParent(t *testing.T) {
	ppid := os.Getppid()
	require.True(t, MatchesParent(ppid))
	require.False(t, MatchesParent(-5))
}
