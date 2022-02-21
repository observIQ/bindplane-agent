package factories

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
)

var defaultExporters = []component.ExporterFactory{
	fileexporter.NewFactory(),
	otlpexporter.NewFactory(),
	otlphttpexporter.NewFactory(),
	observiqexporter.NewFactory(),
	loggingexporter.NewFactory(),
	componenttest.NewNopExporterFactory(),
	googlecloudexporter.NewFactory(),
}
