package status

import (
	"fmt"
	"sync"

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
	Metrics       map[string]*Metric `json:"metrics"`
	sync.Mutex
}

// ToMessage converts a report into a message.
func (r *Report) ToMessage() *message.Message {
	msg, _ := message.New(message.StatusReport, r)
	return msg
}

// Get returns the status of the collector.
func Get() (*Report, error) {
	report := Report{
		ComponentType: "bpagent",
		ComponentID:   "bpagent",
		Status:        ACTIVE,
		Metrics:       make(map[string]*Metric),
	}
	err := report.addPerformanceMetrics()
	if err != nil {
		return nil, err
	}
	return &report, nil
}

type metricGatherer = func(sr *Report) error

func (sr *Report) addPerformanceMetrics() error {
	for _, metricGatherer := range []metricGatherer{AddCPUMetrics, AddMemoryMetrics, AddNetworkMetrics} {
		err := metricGatherer(sr)
		if err != nil {
			return fmt.Errorf("there was an error gathering performance metrics. %s", err)
		}
	}
	return nil
}

func (sr *Report) withMetric(m Metric) {
	sr.Lock()
	defer sr.Unlock()
	sr.Metrics[string(m.Type)] = &m
}
