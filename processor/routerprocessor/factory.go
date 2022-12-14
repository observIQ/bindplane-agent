package routerprocessor

import "go.opentelemetry.io/collector/component"

const (
	// typeStr is the value of the "type" key in the configuration.
	typeStr = "router"

	// stability is the current state of the processor.
	stability = component.StabilityLevelAlpha
)
