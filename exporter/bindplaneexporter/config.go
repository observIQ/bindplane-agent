package bindplaneexporter

import (
	"time"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// defaultTimeout is the default timeout of the exporter
const defaultTimeout = 10 * time.Second

// Config is the config of the bindplane exporter
type Config struct {
	config.ExporterSettings        `mapstructure:",squash"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"`
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	Endpoint string `mapstructure:"endpoint"`
	LiveTail string `mapstructure:"live_tail"`
}

// Validate validates the config
func (c *Config) Validate() error {
	return nil
}

// createDefaultConfig creates a default exporter config
func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		Endpoint:         "localhost:3001/v1/livetail",
		LiveTail:         "/opt/observiq-otel-collector/livetail.yaml",
		TimeoutSettings:  exporterhelper.TimeoutSettings{Timeout: defaultTimeout},
		RetrySettings:    exporterhelper.NewDefaultRetrySettings(),
		QueueSettings:    exporterhelper.NewDefaultQueueSettings(),
	}
}
