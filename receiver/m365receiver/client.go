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
	"context"
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

// HTTP parsing structs
type auth struct {
	Token string `json:"access_token"`
}
type tokenError struct {
	Err string `json:"error"`
}
type defaultError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type jsonError struct {
	Message string `json:"Message"`
}

type logResp struct {
	Content string `json:"contentUri"`
}

// JSON parsing structs

type logData struct {
	logs []jsonLogs
	body []string
}

type jsonLogs struct {
	Workload                 string              `json:"Workload,omitempty"`
	UserID                   string              `json:"UserId"`
	UserType                 int                 `json:"UserType"`
	CreationTime             string              `json:"CreationTime"`
	ID                       string              `json:"Id"`
	Operation                string              `json:"Operation"`
	ResultStatus             string              `json:"ResultStatus,omitempty"`
	RecordType               int                 `json:"RecordType"`
	SharepointSite           string              `json:"Site,omitempty"`
	SharepointSourceFileName string              `json:"SourceFileName,omitempty"`
	ExchangeMailboxGUID      string              `json:"MailboxGuid,omitempty"`
	AzureActor               *[]AzureActor       `json:"Actor,omitempty"`
	DLPSharePointMetaData    *SharePointMetaData `json:"SharePointMetaData,omitempty"`
	DLPExchangeMetaData      *ExchangeMetaData   `json:"ExchangeMetaData,omitempty"`
	DLPPolicyDetails         *[]PolicyDetails    `json:"PolicyDetails,omitempty"`
	SecurityAlertID          string              `json:"AlertId,omitempty"`
	SecurityAlertName        string              `json:"Name,omitempty"` // conflict with teams and investigation DONE-combine
	YammerActorID            string              `json:"ActorUserId,omitempty"`
	YammerFileID             *int                `json:"FileId,omitempty"`
	DefenderEmail            *[]AttachmentData   `json:"AttachmentData,omitempty"`
	DefenderURL              string              `json:"URL,omitempty"` // conflict with investigation DONE-recordType
	DefenderFile             *FileData           `json:"FileData,omitempty"`
	DefenderFileSource       *int                `json:"SourceWorkload,omitempty"`
	InvestigationID          string              `json:"InvestigationId,omitempty"`
	InvestigationStatus      string              `json:"Status,omitempty"`
	PowerAppName             string              `json:"AppName,omitempty"` // conflict with defender DONE-recordType
	DynamicsEntityID         string              `json:"EntityId,omitempty"`
	DynamicsEntityName       string              `json:"EntityName,omitempty"`
	QuarantineSource         *int                `json:"RequestSource,omitempty"`
	FormID                   string              `json:"FormId,omitempty"`
	MIPLabelID               string              `json:"LabelId,omitempty"`
	EncryptedMessageID       string              `json:"MessageId,omitempty"` //conflicting field name with yammer and teams DONE-recordType
	CommCompliance           *ExchangeDetails    `json:"ExchangeDetails,omitempty"`
	ConnectorJobID           string              `json:"JobId,omitempty"`
	ConnectorTaskID          string              `json:"TaskId,omitempty"` //conflict with MS web DONE-combine
	DataShareInvitation      *Invitation         `json:"Invitation,omitempty"`
	MSGraphConsentAppID      string              `json:"ApplicationId,omitempty"` //lots of conflicts DONE-recordType
	VivaGoalsUsername        string              `json:"Username,omitempty"`
	VivaGoalsOrgName         string              `json:"OrganizationName,omitempty"` //conflicts DONE-combine
	MSToDoAppID              string              `json:"ActorAppId,omitempty"`
	MSToDoItemID             string              `json:"ItemID,omitempty"`
	MSWebProjectID           string              `json:"ProjectId,omitempty"`
	MSWebRoadmapID           string              `json:"RoadmapId,omitempty"`
	MSWebRoadmapItemID       string              `json:"RoadmapItemId,omitempty"`
}

// AzureActor json struct
type AzureActor struct {
	ID   string `json:"ID"`
	Type int    `json:"Type"`
}

// SharePointMetaData json struct
type SharePointMetaData struct {
	From string `json:"From"`
}

// ExchangeMetaData json struct
type ExchangeMetaData struct {
	MessageID string `json:"MessageID"`
}

// PolicyDetails json struct
type PolicyDetails struct {
	PolicyID   string `json:"PolicyId"`
	PolicyName string `json:"PolicyName"`
}

// AttachmentData json struct
type AttachmentData struct {
	FileName string `json:"FileName"`
}

// FileData json struct
type FileData struct {
	DocumentID  string `json:"DocumentId"`
	FileVerdict int    `json:"FileVerdict"`
}

// ExchangeDetails json struct
type ExchangeDetails struct {
	NetworkMessageID string `json:"NetworkMessageId,omitempty"`
}

// Invitation json struct
type Invitation struct {
	ShareID string `json:"ShareId"`
}

// client implementation
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
		var respErr defaultError
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, err
		}
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return []string{}, err
		}
		if respErr.Error.Code == "InvalidAuthenticationToken" {
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
func (m *m365Client) GetJSON(ctx context.Context, endpoint string) (logData, error) {
	// make request to Audit endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return logData{}, err
	}
	req.Header.Set("Authorization", "Bearer "+m.token)
	resp, err := m.client.Do(req)
	if err != nil {
		return logData{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	// troubleshoot error code
	if resp.StatusCode != 200 {
		var respErr jsonError
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return logData{}, err
		}
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return logData{}, err
		}
		if respErr.Message == "Authorization has been denied for this request." {
			return logData{}, fmt.Errorf("authorization denied")
		}
		return logData{}, fmt.Errorf("got non 200 status code from request, got %d", resp.StatusCode)
	}

	// check if contentUri field if present
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return logData{}, err
	}
	if string(body) == "[]" { // if body is empty, no new log data available
		return logData{}, nil
	}

	// extract contentUri field
	var lr []logResp
	err = json.Unmarshal(body, &lr)
	if err != nil {
		return logData{}, err
	}

	// read in json
	body, err = m.followLink(ctx, &lr[0])
	if err != nil {
		return logData{}, err
	}
	var data logData
	err = json.Unmarshal(body, &data.logs)
	if err != nil {
		return logData{}, err
	}
	data.body = m.parseBody(body)

	return data, nil
}

func (m *m365Client) StartSubscription(endpoint string) error {
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+m.token)
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// troubleshoot error handling, mainly sub already enabled
	// no error if sub already enabled, not troubleshooting stale token
	// only called while code is synchronous right after a GetT5oken call
	// if token is stale, regenerating it won't fix anything
	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 { // subscription already started possibly
			var respErr defaultError
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return err
			}
			if respErr.Error.Message == "The subscription is already enabled. No property change." {
				return nil
			}
		}
		//unsure what the error is
		return fmt.Errorf("got non 200 status code from request, got %d", resp.StatusCode)
	}
	//if StatusCode == 200, then subscription was successfully started
	return nil
}

// followLink will follow the response of a first request that has a link to the actual content
func (m *m365Client) followLink(ctx context.Context, lr *logResp) ([]byte, error) {
	redirectReq, err := http.NewRequestWithContext(ctx, "GET", lr.Content, nil)
	if err != nil {
		return []byte{}, err
	}
	redirectReq.Header.Set("Authorization", "Bearer "+m.token)
	redirectResp, err := m.client.Do(redirectReq)
	if err != nil {
		return []byte{}, err
	}
	defer func() { _ = redirectResp.Body.Close() }()

	// troubleshoot error code
	if redirectResp.StatusCode != 200 {
		var respErr jsonError
		body, err := io.ReadAll(redirectResp.Body)
		if err != nil {
			return []byte{}, err
		}
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return []byte{}, err
		}
		if respErr.Message == "Authorization has been denied for this request." {
			return []byte{}, fmt.Errorf("authorization denied")
		}
		return []byte{}, fmt.Errorf("got non 200 status code from request, got %d", redirectResp.StatusCode)
	}

	return io.ReadAll(redirectResp.Body)
}

// parseBody takes the byte[] response and parses it into string objects
func (m *m365Client) parseBody(body []byte) []string {
	data := strings.Split(string(body), "},{\"C")
	last := len(data) - 1
	data[0] = strings.TrimPrefix(data[0], "[{\"C")
	data[last] = strings.TrimSuffix(data[last], "}]")
	return data
}
