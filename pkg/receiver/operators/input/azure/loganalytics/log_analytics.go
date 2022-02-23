package loganalytics

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/observiq/observiq-collector/pkg/receiver/operators/input/azure"
	"github.com/open-telemetry/opentelemetry-log-collection/operator"
	"github.com/open-telemetry/opentelemetry-log-collection/operator/helper"
)

const operatorName = "azure_log_analytics_input"

func init() {
	operator.Register(operatorName, func() operator.Builder { return NewLogAnalyticsConfig("") })
}

// NewLogAnalyticsConfig creates a new Azure Log Analytics input config with default values
func NewLogAnalyticsConfig(operatorID string) *InputConfig {
	return &InputConfig{
		InputConfig: helper.NewInputConfig(operatorID, operatorName),
		Config: azure.Config{
			PrefetchCount: 1000,
			StartAt:       "end",
		},
	}
}

// InputConfig is the configuration of a Azure Log Analytics input operator.
type InputConfig struct {
	helper.InputConfig `yaml:",inline"`
	azure.Config       `yaml:",inline"`
}

// Build will build a Azure Log Analytics input operator.
func (c *InputConfig) Build(buildContext operator.BuildContext) ([]operator.Operator, error) {
	if err := c.Config.Build(buildContext, c.InputConfig); err != nil {
		return nil, err
	}

	logAnalyticsInput := &Input{
		EventHub: azure.EventHub{
			Config: c.Config,
		},
		json: jsoniter.ConfigFastest,
	}
	return []operator.Operator{logAnalyticsInput}, nil
}

// Input is an operator that reads Azure Log Analytics input from Azure Event Hub.
type Input struct {
	azure.EventHub
	json jsoniter.API
}

// Start will start generating log entries.
func (l *Input) Start(persister operator.Persister) error {
	l.Handler = l.handleBatchedEvents
	l.Persist = &azure.Persister{DB: persister}
	return l.StartConsumers(context.Background())
}

// Stop will stop generating logs.
func (l *Input) Stop() error {
	return l.StopConsumers()
}
