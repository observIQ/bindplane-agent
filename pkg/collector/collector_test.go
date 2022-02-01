package collector

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/extension/extensionhelper"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"go.opentelemetry.io/collector/service"
)

func TestCollectorRunValid(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

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
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

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
	ctx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	collector := New("./test/invalid.yaml", "0.0.0", nil)
	err := collector.Run(ctx)
	require.EqualError(t, context.DeadlineExceeded, err.Error())

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

func TestCollectorRunInvalidFactory(t *testing.T) {
	extensions := defaultExtensions
	defer func() { defaultExtensions = extensions }()

	defaultExtensions = append(defaultExtensions, extensionhelper.NewFactory(
		"invalid",
		defaultInvalidConfig,
		createInvalidExtension,
	))

	ctx := context.Background()
	// collector.Run should return an error well before the generous deadline is exceeded
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	collector := New("./test/valid.yaml", "0.0.0", nil)
	err := collector.Run(ctx)
	require.Contains(t, err.Error(), "invalid config settings")
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

	err = collector.Run(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "service already running")
}

func TestCollectorRestart(t *testing.T) {
	ctx := context.Background()
	// context must live long enough for the collector to start and restart
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	collector := New("./test/valid.yaml", "0.0.0", nil)
	err := collector.Run(ctx)
	require.NoError(t, err)

	defer collector.Stop()

	err = collector.Restart(ctx)
	require.NoError(t, err)

	status := <-collector.Status()
	require.True(t, status.Running)
}

func TestCollectorPrematureStop(t *testing.T) {
	collector := New("./test/valid.yaml", "0.0.0", nil)
	collector.Stop()
	require.Equal(t, 0, len(collector.Status()))
}

func TestCollectorCreateServicePanic(t *testing.T) {
	defaultPanic := func() config.Receiver {
		panic("expected panic")
	}

	receiver := receiverhelper.NewFactory("panic", defaultPanic)
	factories := component.Factories{
		Receivers: map[config.Type]component.ReceiverFactory{
			"panic": receiver,
		},
	}

	collector := &collector{
		statusChan: make(chan *Status, 10),
		settings: service.CollectorSettings{
			Factories: factories,
		},
	}

	svc, err := collector.createService()
	require.Nil(t, svc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "panic during service creation")
}

func TestCollectorRunServicePanic(t *testing.T) {
	ctx := context.Background()
	collector := &collector{
		statusChan: make(chan *Status, 10),
		wg:         &sync.WaitGroup{},
	}

	collector.wg.Add(1)
	collector.runService(ctx)
	collector.wg.Wait()

	status := <-collector.statusChan
	require.False(t, status.Running)
	require.Error(t, status.Err)
	require.Contains(t, status.Err.Error(), "panic while running service")
}

// InvalidConfig is a config without a mapstructure tag.
type InvalidConfig struct {
	config.ExtensionSettings `mapstructure:",squash"`
	FieldWithoutTag          string
}

// defaultInvalidConfig creates the default invalid config.
func defaultInvalidConfig() config.Extension {
	return &InvalidConfig{
		ExtensionSettings: config.NewExtensionSettings(config.NewComponentID("invalid")),
	}
}

// createInvalidExtension always errors when creating an extension.
func createInvalidExtension(_ context.Context, _ component.ExtensionCreateSettings, _ config.Extension) (component.Extension, error) {
	return nil, errors.New("invalid extension")
}
