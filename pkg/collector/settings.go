package collector

import (
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

const buildDescription = "observIQ's opentelemetry-collector distribution"

// NewSettings returns new settings for the collector with default values.
func NewSettings(configPath string, version string, loggingOpts []zap.Option) service.CollectorSettings {
	factories, _ := DefaultFactories()
	buildInfo := component.BuildInfo{
		Command:     os.Args[0],
		Description: buildDescription,
		Version:     version,
	}
	provider := service.NewDefaultConfigProvider([]string{configPath}, []string{})

	return service.CollectorSettings{
		Factories:               factories,
		BuildInfo:               buildInfo,
		LoggingOptions:          loggingOpts,
		ConfigProvider:          provider,
		DisableGracefulShutdown: true,
	}
}
