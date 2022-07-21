package bindplaneexporter

import (
	"go.opentelemetry.io/collector/component"
)

const (
	typeStr   = "bindplane"
	stability = component.StabilityLevelAlpha
)

func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
	)
}
