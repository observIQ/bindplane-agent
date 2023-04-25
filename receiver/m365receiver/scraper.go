// Copyright  OpenTelemetry Authors
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

package m365receiver // import "github.com/observiq/observiq-otel-collector/receiver/m365receiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/observiq/observiq-otel-collector/receiver/m365receiver/internal/metadata"
)

type reportPair struct {
	columns  map[string]int //relevant column name and index for csv file
	endpoint string         //unique report endpoint
}

var reports = []reportPair{
	{
		columns:  map[string]int{"files": 2, "files.active": 3},
		endpoint: "getSharePointSiteUsageFileCounts(period='D7')",
	},
	{
		columns:  map[string]int{"sites.active": 3},
		endpoint: "getSharePointSiteUsageSiteCounts(period='D7')",
	},
	//TODO: make rest of maps
}

type mClient interface {
	GetCSV(endpoint string) ([]string, error)
	GetToken() error
}

type m365Scraper struct {
	settings component.TelemetrySettings
	cfg      *Config
	client   mClient
	mb       *metadata.MetricsBuilder
}

func newM365Scraper(
	settings receiver.CreateSettings,
	cfg *Config,
) *m365Scraper {
	m := &m365Scraper{
		settings: settings.TelemetrySettings,
		cfg:      cfg,
		mb:       metadata.NewMetricsBuilder(cfg.MetricsBuilderConfig, settings),
	}
	return m
}

func (m *m365Scraper) start(_ context.Context, host component.Host) error {
	httpClient, err := m.cfg.ToClient(host, m.settings)
	if err != nil {
		return err
	}

	m.client = newM365Client(httpClient, m.cfg)

	err = m.client.GetToken()
	if err != nil {
		return err
	}

	return nil
}

// retrieves data, builds metrics & emits them
func (m *m365Scraper) scrape(context.Context) (pmetric.Metrics, error) {
	m365Data, err := m.getStats()
	if err != nil {
		//TODO: error handling
	}

	errs := &scrapererror.ScrapeErrors{}
	now := pcommon.NewTimestampFromTime(time.Now())

	//fmt.Println(m365Data)

	for metricKey, metricValue := range m365Data {
		switch metricKey {
		case "files":
			addPartialIfError(errs, m.mb.RecordM365SharepointFilesCountDataPoint(now, metricValue))
		case "files.active":
			addPartialIfError(errs, m.mb.RecordM365SharepointFilesActiveCountDataPoint(now, metricValue))
		case "sites.active":
			addPartialIfError(errs, m.mb.RecordM365SharepointSitesActiveCountDataPoint(now, metricValue))
		}
	}

	return m.mb.Emit(), errs.Combine()
}

func (m *m365Scraper) getStats() (map[string]string, error) {
	reportData := map[string]string{}

	for _, r := range reports {
		//fmt.Println(r.endpoint)
		line, err := m.client.GetCSV(r.endpoint)
		if err != nil {
			//TODO: error handling
		}

		for k, v := range r.columns {
			reportData[k] = line[v]
		}

	}

	return reportData, nil
}

// adds any errors from recording metric values
func addPartialIfError(errs *scrapererror.ScrapeErrors, err error) {
	if err != nil {
		errs.AddPartial(1, err)
	}
}
