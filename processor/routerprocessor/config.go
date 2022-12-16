// Package routerprocessor contains a processor that routes OTel objects to route receivers based on a match expression.
package routerprocessor

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

// errNoRoutesSpecified is a configuration error when no routes are specified.
var errNoRoutesSpecified = errors.New("must specify at least one route")

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

// Validate validates the configuration.
func (c Config) Validate() error {
	// Ensure at least one route has been specified
	if len(c.Routes) == 0 {
		return errNoRoutesSpecified
	}

	// Check for duplicate routes/match expressions
	routeMatchesLookup := make(map[string]struct{}, len(c.Routes))
	routeNamesLookup := make(map[string]struct{}, len(c.Routes))

	for _, routes := range c.Routes {
		if _, ok := routeMatchesLookup[routes.Match]; ok {
			return fmt.Errorf("duplicate match expression '%s'", routes.Match)
		}
		routeMatchesLookup[routes.Match] = struct{}{}

		if _, ok := routeNamesLookup[routes.Route]; ok {
			return fmt.Errorf("duplicate route name '%s'", routes.Route)
		}
		routeNamesLookup[routes.Route] = struct{}{}
	}

	return nil
}

func createDefaultConfig() component.Config {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
		Routes:            make([]*RoutingRule, 0),
	}
}
