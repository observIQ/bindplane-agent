package status

import (
	"fmt"

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/manager/message"
	"go.uber.org/zap"
)

// Status is the status of the collector.
type Status int

const (
	DISABLED Status = 0
	ACTIVE   Status = 1
	ERROR    Status = 2
)

// Metric is an object that reflects a performance metric of the host that
// the collector is running on
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

// TODO: figure out the structure of this with the other team
// type ErrorContext struct {
// 	Class          string `json:"class"`
// 	UserMessage    string `json:"userMessage" mapstructure:"user_message"`
// 	Recommendation string `json:"recommendation" mapstructure:"recommendation"`
// }

// ToMessage converts the status report struct to a compatible message that can be sent up the websocket connection
func (r *Report) ToMessage() *message.Message {
	msg, _ := message.New(message.StatusReport, r)
	return msg
}

// Get creates a status report struct in order to be converted to a message
// to send up the websocket to reflect current running state of the collector.
func Get(agentID string, status collector.Status) *Report {
	report := &Report{
		ComponentType: "observiq-collector",
		ComponentID:   agentID,
		Status:        ACTIVE,
		Metrics:       make(map[string]*Metric),
	}
	// if status.Err != nil {
	// 	// TODO: fix service side configuration before we say things are in an error state
	// 	// report.Status = ERROR
	// 	report.Context = ErrorContext{
	// 		Class:          "unknown",
	// 		UserMessage:    "Something went wrong with the agent",
	// 		Recommendation: "We recommend that you contact support",
	// 	}
	// }
	return report
}

type metricGatherer struct {
	metricClass string
	collectFunc func(r *Report) error
}

var performanceIndicators = []metricGatherer{
	{
		metricClass: "cpu",
		collectFunc: AddCPUMetrics,
	},
	{
		metricClass: "memory",
		collectFunc: AddMemoryMetrics,
	},
}

// AddPerformanceMetrics will go through and attach performance metrics to the status report
func (sr *Report) AddPerformanceMetrics(logger *zap.Logger) {
	for _, pi := range performanceIndicators {
		err := pi.collectFunc(sr)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to gather performance data for %s", pi.metricClass), zap.Error(err))
			continue
		}
	}
}

func (sr *Report) withMetric(m Metric) {
	sr.Metrics[string(m.Type)] = &m
}
