// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bindplaneextension

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/opampcustommessages"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

const (
	reportMeasurementsCapability = "com.bindplane.measurements"
	measurementsV1               = "measurements.v1"
)

type bindplaneExtension struct {
	cfg                     *Config
	providers               map[component.ID]ThroughputMetricProvider
	customCapabilityHandler opampcustommessages.CustomCapabilityHandler
}

func newBindplaneExtension(cfg *Config) *bindplaneExtension {
	return &bindplaneExtension{
		cfg: cfg,
	}
}

func (b *bindplaneExtension) Start(_ context.Context, host component.Host) error {
	var emptyComponentID component.ID
	if b.cfg.OpAMP != emptyComponentID {
		err := b.setupCustomCapabilities(host)
		if err != nil {
			return fmt.Errorf("setup capability handler: %w", err)
		}
	}

	return nil
}

func (b *bindplaneExtension) setupCustomCapabilities(host component.Host) error {
	ext, ok := host.GetExtensions()[b.cfg.OpAMP]
	if !ok {
		return fmt.Errorf("opamp extension %q does not exist", b.cfg.OpAMP)
	}

	registry, ok := ext.(opampcustommessages.CustomCapabilityRegistry)
	if !ok {
		return fmt.Errorf("extension %q is not an custom message registry", b.cfg.OpAMP)
	}

	var err error
	b.customCapabilityHandler, err = registry.Register(reportMeasurementsCapability)
	if err != nil {
		return fmt.Errorf("register custom capability: %w", err)
	}
}

func (b *bindplaneExtension) Dependencies() []component.ID {
	var emptyComponentID component.ID
	if b.cfg.OpAMP == emptyComponentID {
		return nil
	}

	return []component.ID{b.cfg.OpAMP}
}

func (bindplaneExtension) reportMetricsLoop(_ context.Context) {}

func (b *bindplaneExtension) reportMetrics() {
	allMetrics := pmetric.NewMetrics()
	rm := allMetrics.ResourceMetrics()
	for _, m := range b.providers {
		m := m.Metrics()
		m.MoveAndAppendTo(rm)
	}

	//TODO: Shove metrics into custom message and shuttle over opamp
	return
}

func (bindplaneExtension) Shutdown(_ context.Context) error {
	return nil
}
