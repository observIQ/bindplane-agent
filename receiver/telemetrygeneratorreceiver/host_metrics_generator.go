// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetrygeneratorreceiver //import "github.com/observiq/bindplane-agent/receiver/telemetrygeneratorreceiver"

import (
	"go.uber.org/zap"
)

// hostMetricsGenerator is a generator for host metrics. It generates a sampling of host metrics
// emulating the Host Metrics receiver: github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver

var defaultConfig = GeneratorConfig{
	Type: "host_metrics",
	AdditionalConfig: map[string]any{
		"metrics": []any{
			map[string]any{"name": "system.memory.utilization", "value_min": 0, "value_max": 100, "type": "Gauge", "unit": "1", "attributes": map[string]any{"state": "slab_unreclaimable"}},
			map[string]any{"name": "system.memory.utilization", "value_min": 0, "value_max": 100, "type": "Gauge", "unit": "1", "attributes": map[string]any{"state": "cached"}},
			map[string]any{"name": "system.memory.utilization", "value_min": 0, "value_max": 100, "type": "Gauge", "unit": "1", "attributes": map[string]any{"state": "slab_reclaimable"}},
			map[string]any{"name": "system.memory.utilization", "value_min": 0, "value_max": 100, "type": "Gauge", "unit": "1", "attributes": map[string]any{"state": "buffered"}},
			map[string]any{"name": "system.memory.usage", "value_min": 100000, "value_max": 1000000000, "type": "Sum", "unit": "By", "attributes": map[string]any{"state": "buffered"}},
			map[string]any{"name": "system.memory.usage", "value_min": 100000, "value_max": 1000000000, "type": "Sum", "unit": "By", "attributes": map[string]any{"state": "slab_reclaimable"}},
			map[string]any{"name": "system.memory.usage", "value_min": 100000, "value_max": 1000000000, "type": "Sum", "unit": "By", "attributes": map[string]any{"state": "slab_unreclaimable"}},
			map[string]any{"name": "system.memory.usage", "value_min": 100000, "value_max": 1000000000, "type": "Sum", "unit": "By", "attributes": map[string]any{"state": "cached"}},
			map[string]any{"name": "system.cpu.load_average.1m", "value_min": 0, "value_max": 1, "type": "Gauge", "unit": "{thread}"},
			map[string]any{"name": "system.filesystem.usage", "value_min": 0, "value_max": 15616700416, "type": "Sum", "unit": "By", "attributes": map[string]any{"device": "/dev/vda1", "mode": "rw", "mountpoint": "/etc/hosts", "state": "reserved", "type": "ext4"}},
			map[string]any{"name": "system.filesystem.usage", "value_min": 0, "value_max": 15616700416, "type": "Sum", "unit": "By", "attributes": map[string]any{"device": "/dev/vda1", "mode": "rw", "mountpoint": "/etc/hosts", "state": "free", "type": "ext4"}},
			map[string]any{"name": "system.filesystem.utilization", "value_min": 0, "value_max": 1, "type": "Gauge", "unit": "1", "attributes": map[string]any{"device": "/dev/vda1", "mode": "rw", "mountpoint": "/etc/hosts", "state": "free", "type": "ext4"}},
			map[string]any{"name": "system.filesystem.utilization", "value_min": 0, "value_max": 1, "type": "Gauge", "unit": "1", "attributes": map[string]any{"device": "/dev/vda1", "mode": "rw", "mountpoint": "/etc/hosts", "state": "free", "type": "ext4"}},
			map[string]any{"name": "system.network.packets", "value_min": 0, "value_max": 1000000, "type": "Sum", "unit": "{packets}", "attributes": map[string]any{"device": "eth0", "direction": "receive"}},
			map[string]any{"name": "system.network.packets", "value_min": 0, "value_max": 1000000, "type": "Sum", "unit": "{packets}", "attributes": map[string]any{"device": "eth0", "direction": "send"}},
			map[string]any{"name": "system.network.io", "value_min": 0, "value_max": 100000000, "type": "Sum", "unit": "By", "attributes": map[string]any{"device": "eth0", "direction": "send"}},
			map[string]any{"name": "system.network.io", "value_min": 0, "value_max": 100000000, "type": "Sum", "unit": "By", "attributes": map[string]any{"device": "eth0", "direction": "receive"}},
			map[string]any{"name": "system.network.errors", "value_min": 0, "value_max": 1000, "type": "Sum", "unit": "{errors}", "attributes": map[string]any{"device": "eth0", "direction": "receive"}},
			map[string]any{"name": "system.network.errors", "value_min": 0, "value_max": 1000, "type": "Sum", "unit": "{errors}", "attributes": map[string]any{"device": "eth0", "direction": "transmit"}},
			map[string]any{"name": "system.network.dropped", "value_min": 0, "value_max": 1000, "type": "Sum", "unit": "{packets}", "attributes": map[string]any{"device": "eth0", "direction": "transmit"}},
			map[string]any{"name": "system.network.dropped", "value_min": 0, "value_max": 1000, "type": "Sum", "unit": "{packets}", "attributes": map[string]any{"device": "eth0", "direction": "receive"}},
			map[string]any{"name": "system.network.conntrack.max", "value_min": 65536, "value_max": 65536, "type": "Sum", "unit": "{entries}"},
			map[string]any{"name": "system.network.conntrack.count", "value_min": 8, "value_max": 64, "type": "Sum", "unit": "{entries}"},
			map[string]any{"name": "system.network.connections", "value_min": 0, "value_max": 64, "type": "Sum", "unit": "{connections}", "attributes": map[string]any{"protocol": "tcp", "state": "ESTABLISHED"}},
			map[string]any{"name": "system.network.connections", "value_min": 0, "value_max": 64, "type": "Sum", "unit": "{connections}", "attributes": map[string]any{"protocol": "tcp", "state": "LISTEN"}},
			map[string]any{"name": "system.cpu.time", "value_min": 0, "value_max": 10000, "type": "Sum", "unit": "s", "attributes": map[string]any{"state": "user", "cpu": "cpu0"}},
			map[string]any{"name": "system.cpu.time", "value_min": 0, "value_max": 10000, "type": "Sum", "unit": "s", "attributes": map[string]any{"state": "system", "cpu": "cpu0"}},
			map[string]any{"name": "system.cpu.time", "value_min": 0, "value_max": 10000, "type": "Sum", "unit": "s", "attributes": map[string]any{"state": "idle", "cpu": "cpu0"}},
			map[string]any{"name": "system.cpu.time", "value_min": 0, "value_max": 10000, "type": "Sum", "unit": "s", "attributes": map[string]any{"state": "interrupt", "cpu": "cpu0"}},
			map[string]any{"name": "system.cpu.time", "value_min": 0, "value_max": 10000, "type": "Sum", "unit": "s", "attributes": map[string]any{"state": "nice", "cpu": "cpu0"}},
			map[string]any{"name": "system.cpu.time", "value_min": 0, "value_max": 10000, "type": "Sum", "unit": "s", "attributes": map[string]any{"state": "softirq", "cpu": "cpu0"}},
			map[string]any{"name": "system.cpu.time", "value_min": 0, "value_max": 10000, "type": "Sum", "unit": "s", "attributes": map[string]any{"state": "steal", "cpu": "cpu0"}},
			map[string]any{"name": "system.cpu.time", "value_min": 0, "value_max": 10000, "type": "Sum", "unit": "s", "attributes": map[string]any{"state": "wait", "cpu": "cpu0"}},
		},
	},
}

func newHostMetricsGenerator(cfg GeneratorConfig, logger *zap.Logger) metricGenerator {
	// ignore the config, which only contains the type and resource attributes
	defaultConfig.ResourceAttributes = cfg.ResourceAttributes
	g := newMetricsGenerator(defaultConfig, logger)
	return g
}
