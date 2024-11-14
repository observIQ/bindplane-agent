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
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/extension/experimental/storage"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const (
	eventStorageKey = "last_event_offset"
)

var (
	offset         = 0     // offset for pagination and checkpointing
	exportedEvents = 0     // track the number of events returned by the results endpoint that are exported
	limitReached   = false // flag to stop processing search results when limit is reached
)

type splunksearchapireceiver struct {
	host             component.Host
	logger           *zap.Logger
	logsConsumer     consumer.Logs
	config           *Config
	settings         component.TelemetrySettings
	id               component.ID
	client           splunkSearchAPIClient
	storageClient    storage.Client
	checkpointRecord *EventRecord
}

func newSSAPIReceiver(
	logger *zap.Logger,
	config *Config,
	settings component.TelemetrySettings,
	id component.ID,
) *splunksearchapireceiver {
	return &splunksearchapireceiver{
		logger:           logger,
		config:           config,
		settings:         settings,
		id:               id,
		checkpointRecord: &EventRecord{},
	}
}

func (ssapir *splunksearchapireceiver) Start(ctx context.Context, host component.Host) error {
	ssapir.host = host
	var err error
	ssapir.client, err = newSplunkSearchAPIClient(ctx, ssapir.settings, *ssapir.config, ssapir.host)
	if err != nil {
		return err
	}
	ssapir.client = client

	// create storage client
	storageClient, err := adapter.GetStorageClient(ctx, host, ssapir.config.StorageID, ssapir.id)
	if err != nil {
		return fmt.Errorf("failed to get storage client: %w", err)
	}
	ssapir.storageClient = storageClient

	ssapir.initCheckpoint(ctx)
	go ssapir.runQueries(ctx)
	return nil
}

func (ssapir *splunksearchapireceiver) Shutdown(ctx context.Context) error {
	ssapir.logger.Debug("shutting down logs receiver")

	err := ssapir.checkpoint(ctx)
	if err != nil {
		ssapir.logger.Error("failed checkpoint", zap.Error(err))
	}

	return ssapir.storageClient.Close(ctx)
}

func (ssapir *splunksearchapireceiver) runQueries(ctx context.Context) error {
	for _, search := range ssapir.config.Searches {
		// create search in Splunk
		searchID, err := ssapir.createSplunkSearch(search)
		if err != nil {
			ssapir.logger.Error("error creating search", zap.Error(err))
		}

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

		// fetch search results
		results, err := ssapir.getSplunkSearchResults(ssapir.config, searchID)
		if err != nil {
			ssapir.logger.Error("error fetching search results", zap.Error(err))
		}

			// parse time strings to time.Time
			earliestTime, err := time.Parse(time.RFC3339, search.EarliestTime)
			if err != nil {
				ssapir.logger.Error("earliest_time failed to be parsed as RFC3339", zap.Error(err))
			}

			latestTime, err := time.Parse(time.RFC3339, search.LatestTime)
			if err != nil {
				ssapir.logger.Error("latest_time failed to be parsed as RFC3339", zap.Error(err))
			}

			logs := plog.NewLogs()
			for idx, splunkLog := range results.Results {
				if (idx+exportedEvents) >= search.Limit && search.Limit != 0 {
					limitReached = true
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
					// logger will only log up to 10 times for a given code block, known weird behavior
					continue
				}
				if logTimestamp.UTC().Before(earliestTime) {
					ssapir.logger.Info("skipping log entry - timestamp before earliestTime", zap.Time("time", logTimestamp.UTC()), zap.Time("earliestTime", earliestTime.UTC()))
					break
				}
				log := logs.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()

				// convert time to timestamp
				timestamp := pcommon.NewTimestampFromTime(logTimestamp.UTC())
				log.SetTimestamp(timestamp)
				log.Body().SetStr(splunkLog.Raw)

				if logs.ResourceLogs().Len() == 0 {
					ssapir.logger.Info("search returned no logs within the given time range")
					return nil
				}
			}

			// pass logs, wait for exporter to confirm successful export to GCP
			err = ssapir.logsConsumer.ConsumeLogs(ctx, logs)
			if err != nil {
				// Error from down the pipeline, freak out
				ssapir.logger.Error("error consuming logs", zap.Error(err))
			}
			// last batch of logs has been successfully exported
			exportedEvents += logs.ResourceLogs().Len()
			offset += len(results.Results)

			// update checkpoint
			ssapir.checkpointRecord.Offset = offset
			err = ssapir.checkpoint(ctx)
			if err != nil {
				ssapir.logger.Error("error writing checkpoint", zap.Error(err))
			}
			if limitReached {
				ssapir.logger.Info("limit reached, stopping search result export")
				break
			}
			// if the number of results is less than the results per request, we have queried all pages for the search
			if len(results.Results) < search.EventBatchSize {
				break
			}

		}
		ssapir.logger.Info("search results exported", zap.String("query", search.Query), zap.Int("total results", exportedEvents))
	}
	return nil
}

func (ssapir *splunksearchapireceiver) pollSearchCompletion(ctx context.Context, searchID string) error {
	t := time.NewTicker(ssapir.config.JobPollInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			ssapir.logger.Debug("polling for search completion")
			done, err := ssapir.isSearchCompleted(searchID)
			if err != nil {
				return fmt.Errorf("error polling for search completion: %v", err)
			}
			if done {
				ssapir.logger.Info("search completed")
				return nil
			}
			ssapir.logger.Debug("search not completed yet")
		case <-ctx.Done():
			return nil
		}
	}
}

func (ssapir *splunksearchapireceiver) createSplunkSearch(search Search) (string, error) {
	timeFormat := "%Y-%m-%dT%H:%M:%S"
	searchQuery := fmt.Sprintf("%s starttime=\"%s\" endtime=\"%s\" timeformat=\"%s\"", search.Query, search.EarliestTime, search.LatestTime, timeFormat)
	ssapir.logger.Info("creating search", zap.String("query", searchQuery))
	resp, err := ssapir.client.CreateSearchJob(searchQuery)
	if err != nil {
		return "", err
	}
	return resp.SID, nil
}

func (ssapir *splunksearchapireceiver) isSearchCompleted(sid string) (bool, error) {
	resp, err := ssapir.getJobStatus(ssapir.config, sid)
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

func (ssapir *splunksearchapireceiver) getSplunkSearchResults(sid string, offset int, batchSize int) (SearchResults, error) {
	resp, err := ssapir.getSearchResults(sid, offset, batchSize)
	if err != nil {
		return SearchResultsResponse{}, err
	}
	return resp, nil
}

func (ssapir *splunksearchapireceiver) initCheckpoint(ctx context.Context) error {
	// if a checkpoint already exists, use the offset from the checkpoint
	if err := ssapir.loadCheckpoint(ctx); err != nil {
		return fmt.Errorf("failed to load checkpoint: %w", err)
	}
	if ssapir.checkpointRecord.Offset != 0 {
		// check if the search query in the checkpoint record matches any of the search queries in the config
		for idx, search := range ssapir.config.Searches {
			if search.Query == ssapir.checkpointRecord.Search {
				ssapir.logger.Info("found offset checkpoint in storage extension", zap.Int("offset", ssapir.checkpointRecord.Offset), zap.String("search", ssapir.checkpointRecord.Search))
				// skip searches that have already been processed, use the offset from the checkpoint
				ssapir.config.Searches = ssapir.config.Searches[idx:]
				offset = ssapir.checkpointRecord.Offset
			}
		}
	}
	return nil
}

func (ssapir *splunksearchapireceiver) checkpoint(ctx context.Context) error {
	marshalBytes, err := json.Marshal(ssapir.checkpointRecord)
	if err != nil {
		return fmt.Errorf("failed to write checkpoint: %w", err)
	}
	return ssapir.storageClient.Set(ctx, eventStorageKey, marshalBytes)
}

func (ssapir *splunksearchapireceiver) loadCheckpoint(ctx context.Context) error {
	marshalBytes, err := ssapir.storageClient.Get(ctx, eventStorageKey)
	if err != nil {
		return err
	}
	if marshalBytes == nil {
		ssapir.logger.Info("no checkpoint found")
		return nil
	}
	return json.Unmarshal(marshalBytes, ssapir.checkpointRecord)
}
