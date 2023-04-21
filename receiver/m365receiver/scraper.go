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
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/observiq/observiq-otel-collector/receiver/m365receiver/internal/metadata"
)

type m365Scraper struct {
	settings                  component.TelemetrySettings
	cfg                       *Config
	httpClient                *http.Client
	mb                        *metadata.MetricsBuilder
	metricsAuthorizationToken string
}

func newM365Scraper(
	settings receiver.CreateSettings,
	cfg *Config,
	token string,
) *m365Scraper {
	m := &m365Scraper{
		settings:                  settings.TelemetrySettings,
		cfg:                       cfg,
		mb:                        metadata.NewMetricsBuilder(cfg.MetricsBuilderConfig, settings),
		metricsAuthorizationToken: "Bearer " + token,
	}

	return m
}

func (m *m365Scraper) start(_ context.Context, host component.Host) error {
	httpClient, err := m.cfg.ToClient(host, m.settings)
	if err != nil {
		return err
	}
	m.httpClient = httpClient
	return nil
}

type reportPair struct {
	columns  map[string]int
	endpoint string
}

var (
	rootEndpoint = "https://graph.microsoft.com/v1.0/reports/"
	reports      = []reportPair{
		{
			columns:  map[string]int{"files": 2, "files.active": 3},
			endpoint: "getSharePointSiteUsageFileCounts(period='D7')",
		},
		{
			columns:  map[string]int{"sites.active": 3},
			endpoint: "getSharePointSiteUsageSiteCounts(period='D7')",
		},
	}

	//sharepointFiles = reportPair{map[string]int{"files": 2, "files.active": 3}, "getSharePointSiteUsageFileCounts(period='D7')"}
	//todo: make rest of maps
)

func (m *m365Scraper) scrape(context.Context) (pmetric.Metrics, error) {
	//todo

	//initialize all workers
	//run all workers
	//collect from all workers
	//loop through resulting data and map metrics

	m365Data := m.runWorkers()

	errs := &scrapererror.ScrapeErrors{}
	now := pcommon.NewTimestampFromTime(time.Now())
	for metricKey, metricValue := range m365Data {
		switch metricKey {
		case "files":
			addPartialIfError(errs, m.mb.RecordM365SharepointFilesCountDataPoint(now, strconv.Itoa(metricValue)))
		case "files.active":
			addPartialIfError(errs, m.mb.RecordM365SharepointFilesActiveCountDataPoint(now, strconv.Itoa(metricValue)))
		case "sites.active":
			addPartialIfError(errs, m.mb.RecordM365SharepointSitesActiveCountDataPoint(now, strconv.Itoa(metricValue)))
		}

	}

	return m.mb.Emit(), errs.Combine()
}

func addPartialIfError(errs *scrapererror.ScrapeErrors, err error) {
	if err != nil {
		errs.AddPartial(1, err)
	}
}

func (m *m365Scraper) runWorkers() map[string]int {
	done := make(chan map[string]int)
	var workers = []worker{}

	for _, r := range reports {
		workers = append(workers, *newWorkerStruct(r.columns, m.metricsAuthorizationToken, r.endpoint))
	}

	for _, w := range workers {
		go w.parseData(done, m.httpClient)
	}

	numWorkers := len(workers)
	var results = map[string]int{}

	for i := 0; i < numWorkers; i++ {
		data := <-done
		for k, v := range data {
			results[k] = v
		}
	}

	return results
}

/******************************************************************/
/*********************** Worker Definition ************************/
/******************************************************************/

type worker struct {
	endpoint   string         //api endpoint for this report
	token      string         //authorization token to access apis
	csvData    [][]string     //read the csv data directly into this to avoid downloading files
	csvColumns map[string]int //the column numbers that the relevant data for this report are in
}

// returns a new worker struct
func newWorkerStruct(i_csvColumns map[string]int, i_token string, i_endpoint string) *worker {
	return &worker{
		endpoint:   i_endpoint,
		token:      i_token,
		csvColumns: i_csvColumns,
	}
}

// retrieves csv data from endpoint and places data into csvData object
func (w *worker) getCSVData(httpC *http.Client) {
	req, err := http.NewRequest("GET", w.endpoint, nil)
	if err != nil {
		//todo ERROR HANDLING
	}

	req.Header.Set("Authorization", w.token)
	resp, err := httpC.Do(req)
	if err != nil {
		//todo ERROR HANDLING
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)

	//skip first line of csv data (contains headers)
	_, err = csvReader.Read()
	if err != nil {
		//todo ERROR HANDLING
	}

	//read rest of csv data
	w.csvData, err = csvReader.ReadAll()
	if err != nil {
		//todo ERROR HANDLING
	}
}

// parses through csv data and creates metric points with data
func (w *worker) parseData(done chan map[string]int, httpC *http.Client) {
	w.getCSVData(httpC)

	data := map[string]int{}

	for _, line := range w.csvData {
		for k, v := range w.csvColumns {
			val, err := strconv.Atoi(line[v])
			if err != nil {
				fmt.Println(err)
			}

			data[k] += val
		}
	}

	done <- data
}
