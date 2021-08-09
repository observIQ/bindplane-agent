package status

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func AddCPUMetrics(ctx context.Context, sr *Report) error {
	percentPerCore, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return fmt.Errorf("there was an error reading CPU metrics")
	}
	now := time.Now()
	for core, value := range percentPerCore {
		sr.withMetric(CPUPercent(value, core, now))
	}
	return nil
}

func AddMemoryMetrics(ctx context.Context, sr *Report) error {
	now := time.Now()
	mStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return fmt.Errorf("error getting virtual memory statistics")
	}
	sr.withMetric(MemoryUsed(float64(mStat.Used), now))
	sr.withMetric(MemoryAvailable(float64(mStat.Available), now))
	return nil
}

func CPUPercent(percent float64, core int, t time.Time) Metric {
	return Metric{
		Type:      "cpu_percent",
		Value:     percent,
		Timestamp: t.Unix(),
	}
}

func MemoryUsed(used float64, t time.Time) Metric {
	return Metric{
		Type:      "memory_used",
		Value:     used,
		Timestamp: t.Unix(),
	}
}

func MemoryAvailable(available float64, t time.Time) Metric {
	return Metric{
		Type:      "memory_available",
		Value:     available,
		Timestamp: t.Unix(),
	}
}
