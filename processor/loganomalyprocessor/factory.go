package loganomalyprocessor

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

// NewFactory returns a new factory for the Metrics Generation processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		component.MustNewType("anomaly"),
		createDefaultConfig,
		processor.WithLogs(createLogsProcessor, component.StabilityLevelDevelopment))
}

func createDefaultConfig() component.Config {
	return &Config{
		SampleInterval:   1 * time.Minute,
		MaxWindowAge:     1 * time.Hour,
		ZScoreThreshold:  3.0,
		MADThreshold:     3.5,
		EmergencyMaxSize: 1000,
		Enabled:          true,
		OpAMP:            defaultOpAMPExtensionID,
	}
}

func createLogsProcessor(_ context.Context, set processor.Settings, cfg component.Config, nextConsumer consumer.Logs) (processor.Logs, error) {
	processorConfig, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("configuration parsing error")
	}

	return newProcessor(processorConfig, set.Logger, nextConsumer), nil
}
