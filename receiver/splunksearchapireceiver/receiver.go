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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
	"go.etcd.io/etcd/proxy/grpcproxy/adapter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const (
	eventStorageKey = "last_event_offset"
)

type splunksearchapireceiver struct {
	host          component.Host
	logger        *zap.Logger
	logsConsumer  consumer.Logs
	config        *Config
	settings      component.TelemetrySettings
	client        *http.Client
	storageClient adapter.StorageClient
	record        *eventRecord
}

type eventRecord struct {
	Offset string `json:"offset"`
}

func (ssapir *splunksearchapireceiver) Start(ctx context.Context, host component.Host) error {
	ssapir.host = host
	client, err := ssapir.config.ClientConfig.ToClient(ctx, host, ssapir.settings)
	if err != nil {
		return err
	}
	ssapir.client = client

	// create storage client
	storageClient, err := adapter.GetStorageClient(ssapir.config.StorageID)
	if err != nil {
		return fmt.Errorf("failed to get storage client: %w", err)
	}
	ssapir.storageClient = storageClient

	// check if a checkpoint already exists
	ssapir.loadCheckpoint(ctx)

	go ssapir.runQueries(ctx)
	return nil
}

func (ssapir *splunksearchapireceiver) Shutdown(_ context.Context) error {
	return nil
}

func (ssapir *splunksearchapireceiver) runQueries(ctx context.Context) error {
	for _, search := range ssapir.config.Searches {
		// create search in Splunk
		searchID, err := ssapir.createSplunkSearch(ssapir.config, search.Query)
		if err != nil {
			ssapir.logger.Error("error creating search", zap.Error(err))
		}

		// wait for search to complete
		for {
			done, err := ssapir.isSearchCompleted(ssapir.config, searchID)
			if err != nil {
				ssapir.logger.Error("error checking search status", zap.Error(err))
			}
			if done {
				break
			}
			time.Sleep(2 * time.Second)
		}

		// fetch search results
		results, err := ssapir.getSplunkSearchResults(ssapir.config, searchID)
		if err != nil {
			ssapir.logger.Error("error fetching search results", zap.Error(err))
		}

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
		ssapir.record.Offset = results.Results[len(results.Results)-1].Offset
		err = ssapir.checkpoint(ctx)
		if err != nil {
			ssapir.logger.Error("error writing checkpoint", zap.Error(err))
		}
	}
	return nil
}

func (ssapir *splunksearchapireceiver) createSplunkSearch(config *Config, search string) (string, error) {
	resp, err := ssapir.createSearchJob(config, search)
	if err != nil {
		return "", err
	}
	return resp.SID, nil
}

func (ssapir *splunksearchapireceiver) isSearchCompleted(config *Config, sid string) (bool, error) {
	resp, err := ssapir.getJobStatus(config, sid)
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

func (ssapir *splunksearchapireceiver) getSplunkSearchResults(config *Config, sid string) (SearchResults, error) {
	resp, err := ssapir.getSearchResults(config, sid)
	if err != nil {
		return SearchResults{}, err
	}
	return resp, nil
}

func (ssapir *splunksearchapireceiver) checkpoint(ctx context.Context) error {
	marshalBytes, err := json.Marshal(ssapir.record)
	if err != nil {
		return fmt.Errorf("failed to write checkpoint: %w", err)
	}
	return ssapir.storageClient.Set(ctx, eventStorageKey, marshalBytes)
}

func (ssapir *splunksearchapireceiver) loadCheckpoint(ctx context.Context) {
	marshalBytes, err := ssapir.storageClient.Get(ctx, eventStorageKey)
	if err != nil {
		ssapir.logger.Error("failed to read checkpoint", zap.Error(err))
		return
	}
	if marshalBytes == nil {
		ssapir.logger.Info("no checkpoint found")
		return
	}
	err = json.Unmarshal(marshalBytes, ssapir.record)
	if err != nil {
		ssapir.logger.Error("failed to unmarshal checkpoint", zap.Error(err))
	}
}
