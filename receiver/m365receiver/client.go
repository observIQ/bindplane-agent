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

// Package m365receiver import "github.com/observiq/observiq-otel-collector/receiver/m365receiver"
package m365receiver

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type m365Client struct {
	client       *http.Client
	authEndpoint string
	clientID     string
	clientSecret string
	token        string
	scope        string
}

type auth struct {
	Token string `json:"access_token"`
}
type tokenError struct {
	Err string `json:"error"`
}
type csvError struct {
	ErrorCSV struct {
		Code string `json:"code"`
	} `json:"error"`
}

type jsonError struct {
	Message string `json:"Message"`
}

type logResp struct {
	Content string `json:"contentUri"`
}

type jsonLogs struct {
	OrganizationId string `json:"OrganizationId"`
	Workload       string `json:"Workload,omitempty"`
	UserId         string `json:"UserId"`
	UserType       int    `json:"UserType"`
	CreationTime   string `json:"CreationTime"`
	Id             string `json:"Id"`
	Operation      string `json:"Operation"`
	ResultStatus   string `json:"ResultStatus,omitempty"`
}

func newM365Client(c *http.Client, cfg *Config, scope string) *m365Client {
	return &m365Client{
		client:       c,
		authEndpoint: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", cfg.TenantID),
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		scope:        scope,
	}
}

func (m *m365Client) shutdown() error {
	m.client.CloseIdleConnections()
	return nil
}

// Get authorization token
// metrics token has scope = "https://graph.microsoft.com/.default"
// logs token has scope = "https://manage.office.com/.default"
func (m *m365Client) GetToken() error {
	formData := url.Values{
		"grant_type":    {"client_credentials"},
		"scope":         {m.scope},
		"client_id":     {m.clientID},
		"client_secret": {m.clientSecret},
	}

	requestBody := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", m.authEndpoint, requestBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	//troubleshoot resp err if present
	if resp.StatusCode != 200 {
		//get error code
		var respErr tokenError
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return err
		}
		//match on error code
		switch respErr.Err {
		case "unauthorized_client":
			return fmt.Errorf("the provided client_id is incorrect or does not exist within the given tenant directory")
		case "invalid_client":
			return fmt.Errorf("the provided client_secret is incorrect or does not belong to the given client_id")
		case "invalid_request":
			return fmt.Errorf("the provided tenant_id is incorrect or does not exist")
		}
		return fmt.Errorf("got non 200 status code from request, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var token auth
	err = json.Unmarshal(body, &token)
	if err != nil {
		return err
	}

	m.token = token.Token

	return nil
}

// function for getting metrics data
func (m *m365Client) GetCSV(endpoint string) ([]string, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return []string{}, err
	}

	req.Header.Set("Authorization", m.token)
	resp, err := m.client.Do(req)
	if err != nil {
		return []string{}, err
	}

	//troubleshoot resp err if present
	if resp.StatusCode != 200 {
		//get error code
		var respErr csvError
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, err
		}
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return []string{}, err
		}
		if respErr.ErrorCSV.Code == "InvalidAuthenticationToken" {
			return []string{}, fmt.Errorf("access token invalid")
		}
		return []string{}, fmt.Errorf("got non 200 status code from request, got %d", resp.StatusCode)
	}

	defer func() { _ = resp.Body.Close() }()
	csvReader := csv.NewReader(resp.Body)

	//parse out 2nd line & return csv data
	_, err = csvReader.Read()
	if err != nil {
		return []string{}, err
	}
	data, err := csvReader.Read()
	if err != nil {
		if err == io.EOF { //no data in report, scraper should still run
			return []string{}, nil
		}
		return []string{}, err
	}

	return data, nil
}

// function for getting log data
func (m *m365Client) GetJSON(endpoint string) ([]jsonLogs, error) {
	// make request to Audit endpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return []jsonLogs{}, err
	}
	req.Header.Set("Authorization", "Bearer "+m.token)
	resp, err := m.client.Do(req)
	if err != nil {
		return []jsonLogs{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	// troubleshoot error code
	if resp.StatusCode != 200 {
		var respErr jsonError
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []jsonLogs{}, err
		}
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return []jsonLogs{}, err
		}
		if respErr.Message == "Authorization has been denied for this request." {
			return []jsonLogs{}, fmt.Errorf("authorization denied")
		}
		return []jsonLogs{}, fmt.Errorf("got non 200 status code from request, got %d", resp.StatusCode)
	}

	// check if contentUri field if present
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []jsonLogs{}, err
	}
	if len(body) == 0 { // if body is empty, no new log data available
		return []jsonLogs{}, nil
	}

	// extract contentUri field
	var contentUri logResp
	err = json.Unmarshal(body, &contentUri)
	if err != nil {
		return []jsonLogs{}, err
	}

	// make new GET request to contentUri, including token
	redirectReq, err := http.NewRequest("GET", contentUri.Content, nil)
	if err != nil {
		return []jsonLogs{}, err
	}
	redirectReq.Header.Set("Authorization", "Bearer "+m.token)
	redirectResp, err := m.client.Do(redirectReq)
	if err != nil {
		return []jsonLogs{}, err
	}
	defer func() { _ = redirectResp.Body.Close() }()

	// troubleshoot error code
	if redirectResp.StatusCode != 200 {
		var respErr jsonError
		body, err := io.ReadAll(redirectResp.Body)
		if err != nil {
			return []jsonLogs{}, err
		}
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return []jsonLogs{}, err
		}
		if respErr.Message == "Authorization has been denied for this request." {
			return []jsonLogs{}, fmt.Errorf("authorization denied")
		}
		return []jsonLogs{}, fmt.Errorf("got non 200 status code from request, got %d", redirectResp.StatusCode)
	}

	// read in json
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return []jsonLogs{}, err
	}
	var logData []jsonLogs
	err = json.Unmarshal(body, &logData)
	if err != nil {
		return []jsonLogs{}, err
	}

	return logData, nil
}
