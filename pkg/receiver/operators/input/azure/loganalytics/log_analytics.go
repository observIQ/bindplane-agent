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
func NewLogAnalyticsConfig(operatorID string) *LogAnalyticsInputConfig {
	return &LogAnalyticsInputConfig{
		InputConfig: helper.NewInputConfig(operatorID, operatorName),
		AzureConfig: azure.AzureConfig{
			PrefetchCount: 1000,
			StartAt:       "end",
		},
	}
}

// LogAnalyticsInputConfig is the configuration of a Azure Log Analytics input operator.
type LogAnalyticsInputConfig struct {
	helper.InputConfig `yaml:",inline"`
	azure.AzureConfig  `yaml:",inline"`
}

// Build will build a Azure Log Analytics input operator.
func (c *LogAnalyticsInputConfig) Build(buildContext operator.BuildContext) ([]operator.Operator, error) {
	if err := c.AzureConfig.Build(buildContext, c.InputConfig); err != nil {
		return nil, err
	}

	logAnalyticsInput := &LogAnalyticsInput{
		EventHub: azure.EventHub{
			AzureConfig: c.AzureConfig,
		},
		json: jsoniter.ConfigFastest,
	}
	return []operator.Operator{logAnalyticsInput}, nil
}

// LogAnalyticsInput is an operator that reads Azure Log Analytics input from Azure Event Hub.
type LogAnalyticsInput struct {
	azure.EventHub
	json jsoniter.API
}

// Start will start generating log entries.
func (l *LogAnalyticsInput) Start(persister operator.Persister) error {
	l.Handler = l.handleBatchedEvents
	l.Persist = &azure.Persister{DB: persister}
	return l.StartConsumers(context.Background())
}

// Stop will stop generating logs.
func (l *LogAnalyticsInput) Stop() error {
	return l.StopConsumers()
}
