package collector

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/extension/extensionhelper"
)

func TestCollectorRunValid(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)

	collector := New(settings)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		collector.Run(context.Background())
	}()

	<-time.After(time.Millisecond * 30)
	require.True(t, collector.Running())
	require.NoError(t, collector.Error())

	collector.Stop()
	wg.Wait()
	require.False(t, collector.running)
	require.NoError(t, collector.Error())
}

func TestCollectorRunContext(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	collector := New(settings)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		collector.Run(ctx)
	}()

	<-time.After(time.Millisecond * 30)
	require.True(t, collector.Running())
	require.NoError(t, collector.Error())

	cancel()
	wg.Wait()
	require.False(t, collector.running)
	require.NoError(t, collector.Error())
}

func TestCollectorRunInvalidConfig(t *testing.T) {
	settings, err := NewSettings("./test/invalid.yaml", nil)
	require.NoError(t, err)

	collector := New(settings)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		collector.Run(context.Background())
	}()

	wg.Wait()
	require.False(t, collector.Running())
	require.Error(t, collector.Error())
	require.Contains(t, collector.Error().Error(), "cannot build pipelines")
}

func TestCollectorRunInvalidFactory(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)

	settings.Factories.Extensions["invalid"] = extensionhelper.NewFactory(
		"invalid",
		defaultInvalidConfig,
		createInvalidExtension,
	)

	collector := New(settings)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		collector.Run(context.Background())
	}()

	wg.Wait()
	require.False(t, collector.Running())
	require.Error(t, collector.Error())
	require.Contains(t, collector.Error().Error(), "invalid config settings")
}

// InvalidConfig is a config without a mapstructure tag.
type InvalidConfig struct {
	config.ExtensionSettings `mapstructure:",squash"`
	FieldWithoutTag          string
}

// defaultInvalidConfig creates the default invalid config.
func defaultInvalidConfig() config.Extension {
	return &InvalidConfig{
		ExtensionSettings: config.NewExtensionSettings(config.NewID("invalid")),
	}
}

// createInvalidExtension always errors when creating an extension.
func createInvalidExtension(_ context.Context, _ component.ExtensionCreateSettings, _ config.Extension) (component.Extension, error) {
	return nil, errors.New("invalid extension")
}
