package statusextension

import (
	"context"

	"github.com/observiq/bindplane-agent/extension/statusextension/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

func NewFactory() extension.Factory {
	return extension.NewFactory(
		metadata.Type,
		createDefaultConfig,
		createExtension,
		metadata.ExtensionStability,
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createExtension(_ context.Context, set extension.CreateSettings, _ component.Config) (extension.Extension, error) {
	return newStatusExtension(set.TelemetrySettings), nil
}
