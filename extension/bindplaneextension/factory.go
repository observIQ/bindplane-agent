package bindplaneextension

import (
	"context"

	"github.com/observiq/bindplane-agent/extension/bindplaneextension/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

func NewFactory() extension.Factory {
	return extension.NewFactory(
		metadata.Type,
		defaultConfig,
		createBindPlaneExtension,
		metadata.ExtensionStability,
	)
}

func defaultConfig() component.Config {
	return &Config{}
}

func createBindPlaneExtension(_ context.Context, _ extension.CreateSettings, _ component.Config) (extension.Extension, error) {
	return bindplaneExtension{}, nil
}
