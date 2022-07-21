package bindplaneexporter

import (
	"os"

	"go.opentelemetry.io/collector/config"
)

type Config struct {
	config.ExporterSettings `mapstructure:",squash"`
	LiveTailFile            string `mapstructure:"live_tail_file"`
}

func (c *Config) Validate() error {
	_, err := os.Stat(c.LiveTailFile)
	if err != nil {
		return err
	}
	return nil
}

func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		LiveTailFile:     "/opt/observiq-otel-collector/live_tail.yaml",
	}
}
