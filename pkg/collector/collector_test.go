package collector

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCollectorRunValid(t *testing.T) {
	ctx := context.Background()

	collector := New("./test/valid.yaml", "0.0.0", nil)
	err := collector.Run(ctx)
	require.NoError(t, err)

	status := <-collector.Status()
	require.True(t, status.Running)
	require.NoError(t, status.Err)

	collector.Stop()
	status = <-collector.Status()
	require.False(t, status.Running)
}

func TestCollectorRunMultiple(t *testing.T) {
	collector := New("./test/valid.yaml", "0.0.0", nil)
	for i := 1; i < 5; i++ {
		ctx := context.Background()

		attempt := fmt.Sprintf("Attempt %d", i)
		t.Run(attempt, func(t *testing.T) {
			err := collector.Run(ctx)
			require.NoError(t, err)

			status := <-collector.Status()
			require.True(t, status.Running)
			require.NoError(t, status.Err)

			collector.Stop()
			status = <-collector.Status()
			require.False(t, status.Running)
		})
	}
}

func TestCollectorRunInvalidConfig(t *testing.T) {
	ctx := context.Background()

	collector := New("./test/invalid.yaml", "0.0.0", nil)
	err := collector.Run(ctx)
	require.Error(t, err)

	status := <-collector.Status()
	require.False(t, status.Running)
	require.Error(t, status.Err)
	require.Contains(t, status.Err.Error(), "cannot build receivers")
}

func TestCollectorRunCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	collector := New("./test/valid.yaml", "0.0.0", nil)
	err := collector.Run(ctx)
	require.EqualError(t, context.Canceled, err.Error())
}

func TestCollectorRunTwice(t *testing.T) {
	ctx := context.Background()
	// context must live long enough for the collector to start and attempt to start a second time
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	collector := New("./test/valid.yaml", "0.0.0", nil)
	err := collector.Run(ctx)
	require.NoError(t, err)
	defer collector.Stop()

	status := <-collector.Status()
	require.True(t, status.Running)
	require.NoError(t, status.Err)

	err = collector.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "service already running")

	collector.Stop()
	status = <-collector.Status()
	require.False(t, status.Running)
}

func TestCollectorRestart(t *testing.T) {
	ctx := context.Background()

	collector := New("./test/valid.yaml", "0.0.0", nil)
	err := collector.Run(ctx)
	require.NoError(t, err)

	status := <-collector.Status()
	require.True(t, status.Running)
	require.NoError(t, status.Err)

	err = collector.Restart(ctx)
	require.NoError(t, err)

	status = <-collector.Status()
	require.False(t, status.Running)

	status = <-collector.Status()
	require.True(t, status.Running)

	collector.Stop()
	status = <-collector.Status()
	require.False(t, status.Running)
}

func TestCollectorPrematureStop(t *testing.T) {
	collector := New("./test/valid.yaml", "0.0.0", nil)
	collector.Stop()
	require.Equal(t, 0, len(collector.Status()))
}
