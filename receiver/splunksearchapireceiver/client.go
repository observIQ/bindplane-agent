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

// Package splunksearchapireceiver contains the Splunk Search API receiver.
package splunksearchapireceiver

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

type splunkSearchAPIClient interface {
	CreateSearchJob(search string) (CreateJobResponse, error)
	GetJobStatus(searchID string) (JobStatusResponse, error)
	GetSearchResults(searchID string) (SearchResultsResponse, error)
}

type defaultSplunkSearchAPIClient struct {
	client   *http.Client
	endpoint string
	logger   *zap.Logger
	username string
	password string
}

func newSplunkSearchAPIClient(ctx context.Context, settings component.TelemetrySettings, conf Config, host component.Host) (*defaultSplunkSearchAPIClient, error) {
	client, err := conf.ClientConfig.ToClient(ctx, host, settings)
	if err != nil {
		return nil, err
	}
	return &defaultSplunkSearchAPIClient{
		client:   client,
		endpoint: conf.Endpoint,
		logger:   settings.Logger,
		username: conf.Username,
		password: conf.Password,
	}, nil
}

func (c defaultSplunkSearchAPIClient) CreateSearchJob(search string) (CreateJobResponse, error) {
	endpoint := fmt.Sprintf("%s/services/search/jobs", c.endpoint)

	reqBody := fmt.Sprintf(`search=%s`, search)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return CreateJobResponse{}, err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.client.Do(req)
	if err != nil {
		return CreateJobResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return CreateJobResponse{}, fmt.Errorf("failed to create search job: %d", resp.StatusCode)
	}

	var jobResponse CreateJobResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("failed to read search job create response: %v", err)
	}

	err = xml.Unmarshal(body, &jobResponse)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("failed to unmarshal search job create response: %v", err)
	}
	return jobResponse, nil
}

func (c defaultSplunkSearchAPIClient) GetJobStatus(sid string) (JobStatusResponse, error) {
	endpoint := fmt.Sprintf("%s/services/search/v2/jobs/%s", c.endpoint, sid)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return JobStatusResponse{}, err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.client.Do(req)
	if err != nil {
		return JobStatusResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return JobStatusResponse{}, fmt.Errorf("failed to get search job status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return JobStatusResponse{}, fmt.Errorf("failed to read search job status response: %v", err)
	}
	var jobStatusResponse JobStatusResponse
	err = xml.Unmarshal(body, &jobStatusResponse)
	if err != nil {
		return JobStatusResponse{}, fmt.Errorf("failed to unmarshal search job status response: %v", err)
	}

	return jobStatusResponse, nil
}

func (c defaultSplunkSearchAPIClient) GetSearchResults(sid string) (SearchResultsResponse, error) {
	endpoint := fmt.Sprintf("%s/services/search/v2/jobs/%s/results?output_mode=json", c.endpoint, sid)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return SearchResultsResponse{}, err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.client.Do(req)
	if err != nil {
		return SearchResultsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchResultsResponse{}, fmt.Errorf("failed to get search job results: %d", resp.StatusCode)
	}

	var searchResults SearchResultsResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SearchResultsResponse{}, fmt.Errorf("failed to read search job results response: %v", err)
	}
	// fmt.Println("Body: ", string(body))
	err = json.Unmarshal(body, &searchResults)
	if err != nil {
		return SearchResultsResponse{}, fmt.Errorf("failed to unmarshal search job results response: %v", err)
	}

	return searchResults, nil
}
