package collector

import (
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmapprovider"
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
	//provider := NewFileProvider(configPath)
	provider := configmapprovider.NewDefault(configPath, []string{})

	return service.CollectorSettings{
		Factories:               factories,
		BuildInfo:               buildInfo,
		LoggingOptions:          loggingOpts,
		ConfigMapProvider:       provider,
		DisableGracefulShutdown: true,
	}
}
