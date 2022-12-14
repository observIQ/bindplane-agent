// Package routerprocessor contains a processor that routes OTel objects to route receivers based on a match expression.
package routerprocessor

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

// Config is the config of the router processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
	Routes                   []*RoutingRule `mapstructure:"routes"`
}

// RoutingRule represents a rule for routing OTel objects to route receivers.
type RoutingRule struct {
	Match string `mapstructure:"match"`
	Route string `mapstructure:"route"`
}

func createDefaultConfig() component.Config {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
		Routes:            make([]*RoutingRule, 0),
	}
}
