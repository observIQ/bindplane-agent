package status

import (
	"fmt"

	"github.com/observiq/observiq-collector/manager/message"
)

// Status is the status of the collector.
type Status int

const (
	DISABLED Status = 0
	ACTIVE   Status = 1
	ERROR    Status = 2
)

type Metric struct {
	Type      MetricKey   `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}

// Report is a status report.
type Report struct {
	ComponentType string             `json:"componentType" mapstructure:"componentType"`
	ComponentID   string             `json:"componentID" mapstructure:"componentID"`
	Status        Status             `json:"status" mapstructure:"status"`
	Metrics       map[string]*Metric `json:"metrics" mapstructure:"metrics"`
}

// ErrorReport is a Report with an additional ErrorContext
type ErrorReport struct {
	Context ErrorContext `json:"context" mapstrcuture:"context"`
	*Report
}

// ErrorContext is an object that gives more user friendly context
type ErrorContext struct {
	Class          string      `json:"class,omitempty" mapstructure:"class,omitempty"`
	UserMessage    string      `json:"userMessage,omitempty" mapstructure:"user_message,omitempty"`
	Recommendation string      `json:"recommendation,omitempty" mapstructure:"recommendation,omitempty"`
	Context        interface{} `json:"context,omitempty" mapstructure:"context,omitempty"`
}

func Get() *Report {
	return &Report{
		ComponentType: "bpagent",
		ComponentID:   "bpagent",
		Status:        ACTIVE,
		Metrics:       make(map[string]*Metric),
	}
}

func (r *Report) ToMessage() *message.Message {
	msg, _ := message.New(message.StatusReport, r)
	return msg
}

func Error(err error) *ErrorReport {
	report := &ErrorReport{
		// TODO: classify errors so that we can better present human information
		Context: ErrorContext{
			Class:          "unknown",
			UserMessage:    "Something went wrong with the agent",
			Recommendation: "We recommend that you contact support",
			Context:        make(map[string]interface{}),
		},
	}
	report.ComponentType = "bpagent"
	report.ComponentID = "bpagent"
	report.Status = ERROR
	return report
}

type metricGatherer = func(sr *Report) error

// AddPerformanceMetrics will go through and attach
func (sr *Report) AddPerformanceMetrics() error {
	for _, metricGatherer := range []metricGatherer{AddCPUMetrics, AddMemoryMetrics, AddNetworkMetrics} {
		err := metricGatherer(sr)
		if err != nil {
			return fmt.Errorf("there was an error gathering performance metrics. %s", err)
		}
	}
	return nil
}

func (sr *Report) withMetric(m Metric) {
	sr.Metrics[string(m.Type)] = &m
}
