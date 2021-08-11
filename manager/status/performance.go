package status

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// MetricKey is the status of the collector.
type MetricKey string

// list of metrics being collected/reported to observiq
const (
	CPU_PERCENT      MetricKey = "cpu_percent"
	MEMORY_USED      MetricKey = "memory_used"
	MEMORY_AVAILABLE MetricKey = "memory_available"
	NETWORK_DATA_IN  MetricKey = "network_data_in"
	NETWORK_DATA_OUT MetricKey = "network_data_out"
)

func AddCPUMetrics(sr *Report) error {
	percentPerCore, err := cpu.Percent(0, true)
	if err != nil {
		return fmt.Errorf("there was an error reading CPU metrics: %s", err.Error())
	}
	now := time.Now()
	for core, value := range percentPerCore {
		sr.withMetric(cpuPercent(value, core, now))
	}
	return nil
}

func AddMemoryMetrics(sr *Report) error {
	now := time.Now()
	mStat, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("error getting virtual memory statistics")
	}
	sr.withMetric(memoryUsed(float64(mStat.Used), now))
	sr.withMetric(memoryAvailable(float64(mStat.Available), now))
	return nil
}

func cpuPercent(percent float64, core int, t time.Time) Metric {
	return Metric{
		Type:      CPU_PERCENT,
		Value:     percent,
		Timestamp: t.Unix(),
	}
}

func memoryUsed(used float64, t time.Time) Metric {
	return Metric{
		Type:      MEMORY_USED,
		Value:     used,
		Timestamp: t.Unix(),
	}
}

func memoryAvailable(available float64, t time.Time) Metric {
	return Metric{
		Type:      MEMORY_AVAILABLE,
		Value:     available,
		Timestamp: t.Unix(),
	}
}
