package status

import (
	"context"
	"fmt"

	"github.com/observiq/observiq-collector/manager/message"
)

// Status is the status of the collector.
type Status int

type Metric struct {
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}

const (
	DISABLED Status = 0
	ACTIVE   Status = 1
	ERROR    Status = 2
)

// Report is a status report.
type Report struct {
	ComponentType string             `json:"componentType" mapstructure:"componentType"`
	ComponentID   string             `json:"componentID" mapstructure:"componentID"`
	Status        Status             `json:"status" mapstructure:"status"`
	Metrics       map[string]*Metric `json:"metrics"`
}

// ToMessage converts a report into a message.
func (r *Report) ToMessage() *message.Message {
	msg, _ := message.New(message.StatusReport, r)
	return msg
}

// Get returns the status of the collector.
func Get(ctx context.Context) (*Report, error) {
	report := Report{
		ComponentType: "bpagent",
		ComponentID:   "bpagent",
		Status:        ACTIVE,
		Metrics:       make(map[string]*Metric),
	}
	err := report.addPerformanceMetrics(ctx)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

type metricGatherer = func(ctx context.Context, sr *Report) error

func (sr *Report) addPerformanceMetrics(ctx context.Context) error {
	for _, metricGatherer := range []metricGatherer{AddCPUMetrics, AddMemoryMetrics} {
		err := metricGatherer(ctx, sr)
		if err != nil {
			return fmt.Errorf("there was an error gathering performance metrics. %s", err)
		}
	}
	return nil
}

func (sr *Report) withMetric(m Metric) {
	sr.Metrics[m.Type] = &m
}
