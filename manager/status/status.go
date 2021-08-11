package status

import (
	"fmt"

	"github.com/observiq/observiq-collector/collector"
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
	Context       interface{}        `json:"context,omitempty" mapstrucuture:"context,omitempty"`
}

type ErrorContext struct {
	Class          string `json:"class"`
	UserMessage    string `json:"userMessage" mapstructure:"user_message"`
	Recommendation string `json:"recommendation" mapstructure:"recommendation"`
}

func (r *Report) ToMessage() *message.Message {
	msg, _ := message.New(message.StatusReport, r)
	return msg
}

func Get(status collector.Status) *Report {
	report := &Report{
		ComponentType: "bpagent",
		ComponentID:   "bpagent",
		Status:        ACTIVE,
		Metrics:       make(map[string]*Metric),
	}
	if status.Err != nil {
		report.Status = ERROR
		report.Context = ErrorContext{
			Class:          "unknown",
			UserMessage:    "Something went wrong with the agent",
			Recommendation: "We recommend that you contact support",
		}
	}
	return report
}

type metricGatherer = func(sr *Report) error

// AddPerformanceMetrics will go through and attach
func (sr *Report) AddPerformanceMetrics() error {
	for _, metricGatherer := range []metricGatherer{AddCPUMetrics, AddMemoryMetrics} {
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
