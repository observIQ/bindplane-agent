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

// Package splunksearchapireceiver provides a receiver that uses the Splunk API to migrate event data.
package splunksearchapireceiver

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func (ssapir *splunksearchapireceiver) createSearchJob(config *Config, search string) (CreateJobResponse, error) {
	// fmt.Println("Creating search job for search: ", search)
	endpoint := fmt.Sprintf("%s/services/search/jobs", config.Endpoint)

	reqBody := fmt.Sprintf(`search=%s`, search)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return CreateJobResponse{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	resp, err := ssapir.client.Do(req)
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
		return CreateJobResponse{}, fmt.Errorf("failed to read search job status response: %v", err)
	}

	err = xml.Unmarshal(body, &jobResponse)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("failed to unmarshal search job response: %v", err)
	}
	return jobResponse, nil
}

func (ssapir *splunksearchapireceiver) getJobStatus(config *Config, sid string) (JobStatusResponse, error) {
	// fmt.Println("Getting job status")
	endpoint := fmt.Sprintf("%s/services/search/v2/jobs/%s", config.Endpoint, sid)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return JobStatusResponse{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	resp, err := ssapir.client.Do(req)
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
		return JobStatusResponse{}, fmt.Errorf("failed to unmarshal search job response: %v", err)
	}

	return jobStatusResponse, nil
}

func (ssapir *splunksearchapireceiver) getSearchResults(config *Config, sid string) (SearchResults, error) {
	endpoint := fmt.Sprintf("%s/services/search/v2/jobs/%s/results?output_mode=json", config.Endpoint, sid)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return SearchResults{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	resp, err := ssapir.client.Do(req)
	if err != nil {
		return SearchResults{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchResults{}, fmt.Errorf("failed to get search job results: %d", resp.StatusCode)
	}

	var searchResults SearchResults
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SearchResults{}, fmt.Errorf("failed to read search job results response: %v", err)
	}
	// fmt.Println("Body: ", string(body))
	err = json.Unmarshal(body, &searchResults)
	if err != nil {
		return SearchResults{}, fmt.Errorf("failed to unmarshal search job results: %v", err)
	}

	return searchResults, nil
}
