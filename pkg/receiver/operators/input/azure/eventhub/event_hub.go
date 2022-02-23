package eventhub

import (
	"context"

	"github.com/observiq/observiq-collector/pkg/receiver/operators/input/azure"
	"github.com/open-telemetry/opentelemetry-log-collection/operator"
	"github.com/open-telemetry/opentelemetry-log-collection/operator/helper"
)

const operatorName = "azure_event_hub_input"

func init() {
	operator.Register(operatorName, func() operator.Builder { return NewEventHubConfig("") })
}

// NewEventHubConfig creates a new Azure Event Hub input config with default values
func NewEventHubConfig(operatorID string) *InputConfig {
	return &InputConfig{
		InputConfig: helper.NewInputConfig(operatorID, operatorName),
		Config: azure.Config{
			PrefetchCount: 1000,
			StartAt:       "end",
		},
	}
}

// InputConfig is the configuration of a Azure Event Hub input operator.
type InputConfig struct {
	helper.InputConfig `yaml:",inline"`
	azure.Config       `yaml:",inline"`
}

// Build will build a Azure Event Hub input operator.
func (c *InputConfig) Build(buildContext operator.BuildContext) ([]operator.Operator, error) {
	if err := c.Config.Build(buildContext, c.InputConfig); err != nil {
		return nil, err
	}

	eventHubInput := &Input{
		EventHub: azure.EventHub{
			Config: c.Config,
		},
	}
	return []operator.Operator{eventHubInput}, nil
}

// Input is an operator that reads input from Azure Event Hub.
type Input struct {
	azure.EventHub
}

// Start will start generating log entries.
func (e *Input) Start(persister operator.Persister) error {
	e.Handler = e.handleEvent
	e.Persist = &azure.Persister{DB: persister}
	return e.StartConsumers(context.Background())
}

// Stop will stop generating logs.
func (e *Input) Stop() error {
	return e.StopConsumers()
}
