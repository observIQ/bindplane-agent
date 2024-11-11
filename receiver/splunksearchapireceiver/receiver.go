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

package splunksearchapireceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type splunksearchapireceiver struct {
	host         component.Host
	logger       *zap.Logger
	logsConsumer consumer.Logs
	config       *Config
	settings     component.TelemetrySettings
	client       splunkSearchAPIClient
}

func (ssapir *splunksearchapireceiver) Start(ctx context.Context, host component.Host) error {
	ssapir.host = host
	var err error
	ssapir.client, err = newSplunkSearchAPIClient(ctx, ssapir.settings, *ssapir.config, ssapir.host)
	if err != nil {
		return err
	}
	go ssapir.runQueries(ctx)
	return nil
}

func (ssapir *splunksearchapireceiver) Shutdown(_ context.Context) error {
	return nil
}

func (ssapir *splunksearchapireceiver) runQueries(ctx context.Context) error {
	for _, search := range ssapir.config.Searches {
		// create search in Splunk
		searchID, err := ssapir.createSplunkSearch(search.Query)
		if err != nil {
			ssapir.logger.Error("error creating search", zap.Error(err))
		}
		// fmt.Println("Search created successfully with ID: ", searchID)

		// wait for search to complete
		for {
			done, err := ssapir.isSearchCompleted(searchID)
			if err != nil {
				ssapir.logger.Error("error checking search status", zap.Error(err))
			}
			if done {
				break
			}
			time.Sleep(2 * time.Second)
		}
		// fmt.Println("Search completed successfully")

		// fetch search results
		results, err := ssapir.getSplunkSearchResults(searchID)
		if err != nil {
			ssapir.logger.Error("error fetching search results", zap.Error(err))
		}
		// fmt.Println("Search results: ", results)

		// parse time strings to time.Time
		earliestTime, err := time.Parse(time.RFC3339, search.EarliestTime)
		if err != nil {
			// should be impossible to reach with config validation
			ssapir.logger.Error("earliest_time failed to be parsed as RFC3339", zap.Error(err))
		}

		latestTime, err := time.Parse(time.RFC3339, search.LatestTime)
		if err != nil {
			// should be impossible to reach with config validation
			ssapir.logger.Error("latest_time failed to be parsed as RFC3339", zap.Error(err))
		}

		logs := plog.NewLogs()
		for idx, splunkLog := range results.Results {
			if idx >= search.Limit && search.Limit != 0 {
				break
			}
			// convert log timestamp to ISO8601 (UTC() makes RFC3339 into ISO8601)
			logTimestamp, err := time.Parse(time.RFC3339, splunkLog.Time)
			if err != nil {
				ssapir.logger.Error("error parsing log timestamp", zap.Error(err))
				break
			}
			if logTimestamp.UTC().After(latestTime.UTC()) {
				ssapir.logger.Info("skipping log entry - timestamp after latestTime", zap.Time("time", logTimestamp.UTC()), zap.Time("latestTime", latestTime.UTC()))
				// logger.Info will only log up to 10 times for a given code block, known weird behavior
				continue
			}
			if logTimestamp.UTC().Before(earliestTime) {
				ssapir.logger.Info("skipping log entry - timestamp before earliestTime", zap.Time("time", logTimestamp.UTC()), zap.Time("earliestTime", earliestTime.UTC()))
				continue
			}
			log := logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()

			// convert time to timestamp
			timestamp := pcommon.NewTimestampFromTime(logTimestamp.UTC())
			log.SetTimestamp(timestamp)
			log.Body().SetStr(splunkLog.Raw)

		}
		if logs.ResourceLogs().Len() == 0 {
			ssapir.logger.Info("search returned no logs within the given time range")
			return nil
		}

		// pass logs, wait for exporter to confirm successful export to GCP
		err = ssapir.logsConsumer.ConsumeLogs(ctx, logs)
		if err != nil {
			// Error from down the pipeline, freak out
			ssapir.logger.Error("error consuming logs", zap.Error(err))
		}
	}
	return nil
}

func (ssapir *splunksearchapireceiver) createSplunkSearch(search string) (string, error) {
	resp, err := ssapir.client.CreateSearchJob(search)
	if err != nil {
		return "", err
	}
	return resp.SID, nil
}

func (ssapir *splunksearchapireceiver) isSearchCompleted(sid string) (bool, error) {
	resp, err := ssapir.client.GetJobStatus(sid)
	if err != nil {
		return false, err
	}

	for _, key := range resp.Content.Dict.Keys {
		if key.Name == "dispatchState" {
			if key.Value == "DONE" {
				return true, nil
			}
			break
		}
	}
	return false, nil
}

func (ssapir *splunksearchapireceiver) getSplunkSearchResults(sid string) (SearchResultsResponse, error) {
	resp, err := ssapir.client.GetSearchResults(sid)
	if err != nil {
		return SearchResultsResponse{}, err
	}
	return resp, nil
}
