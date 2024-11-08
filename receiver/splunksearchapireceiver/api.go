package splunksearchapireceiver

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func createHttpClient() *http.Client {
	// TODO: Add functionality to configure TLS settings using config options
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disables TLS verification
	}
	return &http.Client{Transport: tr}
}

func createSearchJob(config *Config, search string) (CreateJobResponse, error) {
	// fmt.Println("Creating search job for search: ", search)
	endpoint := fmt.Sprintf("https://%s/services/search/jobs", config.Server)

	reqBody := fmt.Sprintf(`search=%s`, search)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return CreateJobResponse{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	client := createHttpClient()
	resp, err := client.Do(req)
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

func getJobStatus(config *Config, sid string) (JobStatusResponse, error) {
	// fmt.Println("Getting job status")
	endpoint := fmt.Sprintf("https://%s/services/search/v2/jobs/%s", config.Server, sid)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return JobStatusResponse{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	client := createHttpClient()
	resp, err := client.Do(req)
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

func getSearchResults(config *Config, sid string) (SearchResults, error) {
	endpoint := fmt.Sprintf("https://%s/services/search/v2/jobs/%s/results?output_mode=json", config.Server, sid)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return SearchResults{}, err
	}
	req.SetBasicAuth(config.Username, config.Password)

	client := createHttpClient()
	resp, err := client.Do(req)
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
