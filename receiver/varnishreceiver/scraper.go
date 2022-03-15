// Copyright  The OpenTelemetry Authors
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

package varnishreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/varnishreceiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"

	"github.com/observIQ/observiq-otel-collector/pkg/receiver/varnishreceiver/internal/metadata"
)

type varnishScraper struct {
	client   client
	config   *Config
	settings component.TelemetrySettings
	mb       *metadata.MetricsBuilder
}

func (v *varnishScraper) start(_ context.Context, host component.Host) error {
	v.client = newVarnishClient(v.config, host, v.settings)
	return nil
}

func newVarnishScraper(settings component.TelemetrySettings, config *Config) *varnishScraper {
	return &varnishScraper{
		settings: settings,
		config:   config,
		mb:       metadata.NewMetricsBuilder(metadata.DefaultMetricsSettings()),
	}
}

func (v *varnishScraper) scrape(context.Context) (pdata.Metrics, error) {
	stats, err := v.client.GetStats()
	if err != nil {
		v.settings.Logger.Error("Failed to execute varnishstat",
			zap.String("Working Directory:", v.config.WorkingDir),
			zap.String("Executable Directory:", v.config.ExecDir),
			zap.Error(err),
		)
		return pdata.NewMetrics(), err
	}

	now := pdata.NewTimestampFromTime(time.Now())
	md := v.mb.NewMetricData()

	v.recordVarnishBackendConnectionsCountDataPoint(now, stats)
	v.recordVarnishCacheOperationsCountDataPoint(now, stats)
	v.recordVarnishThreadOperationsCountDataPoint(now, stats)
	v.recordVarnishSessionCountDataPoint(now, stats)

	v.mb.RecordVarnishObjectExpiredCountDataPoint(now, stats.Counters.MAINNExpired.Value)
	v.mb.RecordVarnishObjectNukedCountDataPoint(now, stats.Counters.MAINNLruNuked.Value)
	v.mb.RecordVarnishObjectMovedCountDataPoint(now, stats.Counters.MAINNLruMoved.Value)
	v.mb.RecordVarnishObjectCountDataPoint(now, stats.Counters.MAINNObject.Value)
	v.mb.RecordVarnishClientRequestsCountDataPoint(now, stats.Counters.MAINClientReq.Value)
	v.mb.RecordVarnishBackendRequestsCountDataPoint(now, stats.Counters.MAINBackendReq.Value)

	v.mb.Emit(md.ResourceMetrics().At(0).InstrumentationLibraryMetrics().At(0).Metrics())
	return md, nil
}
