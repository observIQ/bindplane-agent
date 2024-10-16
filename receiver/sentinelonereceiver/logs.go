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

package sentinelonereceiver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type sentinelOneLogsReceiver struct {
	apiURL         string
	cfg            Config
	client         httpClient
	consumer       consumer.Logs
	logger         *zap.Logger
	cancel         context.CancelFunc
	activitiesChan chan []Activity
	since          time.Time
	wg             *sync.WaitGroup
}

type sentinelOneAPIResponse struct {
	Pagination PaginationData `json:"pagination"`
	Data       []Activity     `json:"data"`
	Errors     []string       `json:"errors"`
}

// PaginationData contains pagination details for the response
type PaginationData struct {
	TotalItems int    `json:"totalItems"`
	NextCursor string `json:"nextCursor"`
}

// Activity represents individual activity data in the response
type Activity struct {
	AccountID            string `json:"accountId"`
	AccountName          string `json:"accountName"`
	ActivityType         int    `json:"activityType"`
	ActivityUUID         string `json:"activityUuid"`
	AgentID              string `json:"agentId"`
	AgentUpdatedVersion  string `json:"agentUpdatedVersion"`
	Comments             string `json:"comments"`
	CreatedAt            string `json:"createdAt"`
	Data                 any    `json:"data"`
	Description          string `json:"description"`
	GroupID              string `json:"groupId"`
	GroupName            string `json:"groupName"`
	Hash                 string `json:"hash"`
	ID                   string `json:"id"`
	OSFamily             string `json:"osFamily"`
	PrimaryDescription   string `json:"primaryDescription"`
	SecondaryDescription string `json:"secondaryDescription"`
	SiteID               string `json:"siteId"`
	SiteName             string `json:"siteName"`
	ThreatID             string `json:"threatId"`
	UpdatedAt            string `json:"updatedAt"`
	UserID               string `json:"userId"`
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

var (
	activities         = "activities"
	agents             = "agents"
	threats            = "threats"
	activitiesMaxLimit = 1000
)

// newSentinelOneLogsReceiver returns a newly configured sentinelOneLogsReceiver
func newSentinelOneLogsReceiver(cfg *Config, logger *zap.Logger, consumer consumer.Logs) (*sentinelOneLogsReceiver, error) {
	return &sentinelOneLogsReceiver{
		apiURL:         cfg.BaseURL + "/web/api/v2.1/",
		cfg:            *cfg,
		client:         http.DefaultClient,
		consumer:       consumer,
		logger:         logger,
		since:          time.Now().Add(-cfg.PollInterval),
		activitiesChan: make(chan []Activity, len(cfg.APIs)),
		wg:             &sync.WaitGroup{},
	}, nil
}

func (r *sentinelOneLogsReceiver) Start(_ context.Context, _ component.Host) error {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.wg.Add(1)
	go r.startPolling(ctx)
	return nil
}

func (r *sentinelOneLogsReceiver) startPolling(ctx context.Context) {
	defer r.wg.Done()
	t := time.NewTicker(r.cfg.PollInterval)

	err := r.poll(ctx)
	if err != nil {
		r.logger.Error("there was an error during the first poll", zap.Error(err))
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			err := r.poll(ctx)
			if err != nil {
				r.logger.Error("there was an error during the poll", zap.Error(err))
			}
		}
	}
}

func (r *sentinelOneLogsReceiver) poll(ctx context.Context) error {
	logEvents := r.getLogs(ctx)
	observedTime := pcommon.NewTimestampFromTime(time.Now())
	logs := r.processLogEvents(observedTime, logEvents)
	if logs.LogRecordCount() > 0 {
		if err := r.consumer.ConsumeLogs(ctx, logs); err != nil {
			return err
		}
	}
	return nil
}

func (r *sentinelOneLogsReceiver) getLogs(ctx context.Context) []Activity {
	if slices.Contains(r.cfg.APIs, activities) {
		go r.getActivities(ctx)
	}

	// if slices.Contains(r.cfg.APIs, agents) {
	// 	go r.getAgents(ctx)
	// }

	// if slices.Contains(r.cfg.APIs, threats) {
	// 	go r.getThreats(ctx)
	// }

	combinedLogs := []Activity{}
	for logs := range r.activitiesChan {
		combinedLogs = append(combinedLogs, logs...)
	}

	// TODO sort logs

	return combinedLogs
}

func (r *sentinelOneLogsReceiver) getActivities(ctx context.Context) {
	var activityLogs []Activity
	for {
		req, err := http.NewRequestWithContext(ctx, "GET", r.apiURL+activities, nil)
		if err != nil {
			r.logger.Error("error creating sentinelone activities request", zap.Error(err))
			break
		}
		query := req.URL.Query()
		query.Add("limit", strconv.Itoa(activitiesMaxLimit))
		// until := time.Now().Unix()
		// query.Add("createdAt__between", fmt.Sprintf("%d-%d", r.since.Unix(), until))
		req.URL.RawQuery = query.Encode()

		req.Header.Add("Authorization", "ApiToken "+string(r.cfg.APIToken))

		res, err := r.client.Do(req)
		if err != nil {
			r.logger.Error("error performing sentinelone activities request", zap.Error(err))
			break
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			r.logger.Error("okta logs endpoint returned non-200 statuscode: " + res.Status)
			break
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			r.logger.Error("error reading response body", zap.Error(err))
			break
		}

		var apiResponse sentinelOneAPIResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			r.logger.Error("unable to unmarshal log events", zap.Error(err))
			break
		}

		fmt.Println("TotalItems:", apiResponse.Pagination.TotalItems)

		activityLogs = append(activityLogs, apiResponse.Data...)

		if apiResponse.Pagination.NextCursor == "" {
			break
		}

		// TODO cursor logic & until logic
	}
	r.activitiesChan <- activityLogs
}

// func (r *sentinelOneLogsReceiver) getAgents(ctx context.Context) {
// 	for {
// 		req, err := http.NewRequestWithContext(ctx, "GET")
// 	}
// 	r.logsChan <- logs
// }

// func (r *sentinelOneLogsReceiver) getThreats(ctx context.Context) {
// 	for {
// 		req, err := http.NewRequestWithContext(ctx, "GET")
// 	}
// 	r.logsChan <- logs
// }

func (r *sentinelOneLogsReceiver) processLogEvents(observedTime pcommon.Timestamp, logEvents []Activity) plog.Logs {
	return plog.NewLogs()
}

func (r *sentinelOneLogsReceiver) Shutdown(_ context.Context) error {
	r.logger.Debug("shutting down logs receiver")
	if r.cancel != nil {
		r.cancel()
	}
	r.client.CloseIdleConnections()
	r.wg.Wait()
	return nil
}
