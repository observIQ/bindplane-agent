package status

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
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
		return fmt.Errorf("there was an error reading CPU metrics")
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

func AddNetworkMetrics(sr *Report) error {
	now := time.Now()
	netStat, err := net.IOCounters(true) // Per Nic
	if err != nil {
		return fmt.Errorf("error getting network interface statistics")
	}
	for _, nic := range cleanNetStat(netStat) {
		sr.withMetric(networkIn(nic.BytesRecv, now))
		sr.withMetric(networkOut(nic.BytesSent, now))
	}
	return nil
}

// In some cases, such as when running vmware fusion, it is possible for multiple nics
// to be returned with the same name. Aggregation is perhaps not a perfect solution,
// but it is preferrable to causing errors downstream, and it does still represent
// something meaningful, in that it would be the total for a virtualization program
func cleanNetStat(netStat []net.IOCountersStat) []net.IOCountersStat {
	aggregatedStats := map[string]net.IOCountersStat{}

	for _, nic := range netStat {
		if total, exists := aggregatedStats[nic.Name]; exists {
			aggregatedStats[nic.Name] = sumStats(total, nic)
		} else {
			aggregatedStats[nic.Name] = nic
		}
	}

	result := []net.IOCountersStat{}
	for _, stat := range aggregatedStats {
		result = append(result, stat)
	}
	return result
}

func sumStats(nic1, nic2 net.IOCountersStat) net.IOCountersStat {
	return net.IOCountersStat{
		Name:        nic1.Name,
		BytesSent:   nic1.BytesSent + nic2.BytesSent,
		BytesRecv:   nic1.BytesRecv + nic2.BytesRecv,
		PacketsSent: nic1.PacketsSent + nic2.PacketsSent,
		PacketsRecv: nic1.PacketsRecv + nic2.PacketsRecv,
		Errin:       nic1.Errin + nic2.Errin,
		Errout:      nic1.Errout + nic2.Errout,
		Dropin:      nic1.Dropin + nic2.Dropin,
		Dropout:     nic1.Dropout + nic2.Dropout,
		Fifoin:      nic1.Fifoin + nic2.Fifoin,
		Fifoout:     nic1.Fifoout + nic2.Fifoout,
	}
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

func networkIn(bytesR uint64, t time.Time) Metric {
	return Metric{
		Type:      NETWORK_DATA_IN,
		Value:     bytesR,
		Timestamp: t.Unix(),
	}
}

func networkOut(bytesO uint64, t time.Time) Metric {
	return Metric{
		Type:      NETWORK_DATA_OUT,
		Value:     bytesO,
		Timestamp: t.Unix(),
	}
}
