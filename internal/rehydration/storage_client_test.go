package rehydration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNopStorage(t *testing.T) {
	storage := NewNopStorage()
	require.NoError(t, storage.SaveCheckpoint(context.Background(), "key", &CheckPoint{}))

	checkpoint, err := storage.LoadCheckPoint(context.Background(), "key")
	require.Equal(t, &CheckPoint{}, checkpoint)
	require.NoError(t, err)
}
