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

package proofpointreceiver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

var (
	proofpointAPIAllEventsURL = "https://tap-api-v2.proofpoint.com/v2/siem/all"

	// ISO8601Format
	ISO8601Format = "2006-01-02T15:04:05Z"
)

type proofpointLogsReceiver struct {
	cfg       Config
	client    httpClient
	consumer  consumer.Logs
	logger    *zap.Logger
	cancel    context.CancelFunc
	sinceTime string
	wg        *sync.WaitGroup
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

// newProofpointLogsReceiver returns a newly configured proofpointLogsReceiver
func newProofpointLogsReceiver(cfg *Config, logger *zap.Logger, consumer consumer.Logs) (*proofpointLogsReceiver, error) {
	return &proofpointLogsReceiver{
		cfg:      *cfg,
		client:   http.DefaultClient,
		consumer: consumer,
		logger:   logger,
		wg:       &sync.WaitGroup{},
	}, nil
}

func (r *proofpointLogsReceiver) Start(_ context.Context, _ component.Host) error {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.wg.Add(1)
	go r.startPolling(ctx)
	return nil
}

func (r *proofpointLogsReceiver) startPolling(ctx context.Context) {
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

func (r *proofpointLogsReceiver) poll(ctx context.Context) error {
	ppResponse, err := r.getLogs(ctx)
	if err != nil {
		return err
	}
	observedTime := pcommon.NewTimestampFromTime(time.Now())
	logs := r.processLogEvents(observedTime, ppResponse)
	if logs.LogRecordCount() > 0 {
		if err := r.consumer.ConsumeLogs(ctx, logs); err != nil {
			return err
		}
	}
	return nil
}

func (r *proofpointLogsReceiver) getLogs(ctx context.Context) (*proofpointResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", proofpointAPIAllEventsURL, nil)
	req.SetBasicAuth(string(r.cfg.Principal), string(r.cfg.Secret))
	if err != nil {
		return nil, err
	}

	if r.sinceTime == "" {
		r.sinceTime = time.Now().UTC().Add(-r.cfg.PollInterval).Format(ISO8601Format)
	}

	query := req.URL.Query()
	query.Add("sinceTime", r.sinceTime)
	query.Add("format", "json")
	req.URL.RawQuery = query.Encode()

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responseJson proofpointResponse
	err = json.Unmarshal(body, &responseJson)
	if err != nil {
		return nil, err
	}

	r.sinceTime = responseJson.QueryEndTime
	return &responseJson, nil
}

func (r *proofpointLogsReceiver) processLogEvents(observedTime pcommon.Timestamp, response *proofpointResponse) plog.Logs {
	logs := plog.NewLogs()

	resourceLogs := logs.ResourceLogs().AppendEmpty()
	resourceLogs.ScopeLogs().AppendEmpty()

	for _, event := range response.ClicksBlocked {
		r.processClickEvent(observedTime, resourceLogs, event, "clicksBlocked")
	}

	for _, event := range response.ClicksPermitted {
		r.processClickEvent(observedTime, resourceLogs, event, "clicksPermitted")
	}

	for _, event := range response.MessagesBlocked {
		r.processMessageEvent(observedTime, resourceLogs, event, "messagesBlocked")
	}

	for _, event := range response.MessagesDelivered {
		r.processMessageEvent(observedTime, resourceLogs, event, "messagesDelivered")
	}

	return logs
}

func (r *proofpointLogsReceiver) processClickEvent(observedTime pcommon.Timestamp, resourceLogs plog.ResourceLogs, event ClickEvent, eventType string) {
	logRecord := resourceLogs.ScopeLogs().At(0).LogRecords().AppendEmpty()

	// timestamps
	logRecord.SetObservedTimestamp(observedTime)
	timestamp := time.UnixMilli(event.ClickTime.UnixMilli())
	logRecord.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))

	// body
	eventBytes, err := json.Marshal(event)
	if err != nil {
		r.logger.Error("unable to marshal event", zap.Error(err))
	} else {
		logRecord.Body().SetStr(string(eventBytes))
	}

	// attributes
	attributes := logRecord.Attributes()
	attributes.PutStr("event.type", eventType)
	attributes.PutStr("id", event.ID)
	attributes.PutStr("GUID", event.GUID)
	attributes.PutStr("threat.id", event.ThreatID)
	attributes.PutStr("threat.status", event.ThreatStatus)
	attributes.PutStr("threat.url", event.ThreatURL)
	attributes.PutStr("classification", event.Classification)
	attributes.PutStr("recipient", event.Recipient)
	attributes.PutStr("sender", event.Sender)
	attributes.PutStr("url", event.URL)
}

func (r *proofpointLogsReceiver) processMessageEvent(observedTime pcommon.Timestamp, resourceLogs plog.ResourceLogs, event MessageEvent, eventType string) {
	logRecord := resourceLogs.ScopeLogs().At(0).LogRecords().AppendEmpty()

	// timestamps
	logRecord.SetObservedTimestamp(observedTime)
	timestamp := time.UnixMilli(event.MessageTime.UnixMilli())
	logRecord.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))

	// body
	eventBytes, err := json.Marshal(event)
	if err != nil {
		r.logger.Error("unable to marshal event", zap.Error(err))
	} else {
		logRecord.Body().SetStr(string(eventBytes))
	}

	// attributes
	attributes := logRecord.Attributes()
	attributes.PutStr("event.type", eventType)
	attributes.PutStr("GUID", event.GUID)
	attributes.PutStr("message.id", event.MessageID)
	attributes.PutInt("malwareScore", int64(event.MalwareScore))
	attributes.PutInt("imposterScore", int64(event.ImposterScore))
	attributes.PutInt("phishScore", int64(event.PhishScore))
	attributes.PutInt("spamScore", int64(event.SpamScore))
	attributes.PutStr("QID", event.QID)
	attributes.PutStr("recipient", event.Recipient)
	attributes.PutStr("sender", event.Sender)
	attributes.PutStr("subject", event.Subject)
}

func (r *proofpointLogsReceiver) Shutdown(_ context.Context) error {
	r.logger.Debug("shutting down logs receiver")
	if r.cancel != nil {
		r.cancel()
	}
	r.client.CloseIdleConnections()
	r.wg.Wait()
	return nil
}

type proofpointResponse struct {
	QueryEndTime      string         `json:"queryEndTime"`
	MessagesDelivered []MessageEvent `json:"messagesDelivered"`
	MessagesBlocked   []MessageEvent `json:"messagesBlocked"`
	ClicksPermitted   []ClickEvent   `json:"clicksPermitted"`
	ClicksBlocked     []ClickEvent   `json:"clicksBlocked"`
}

type MessageEvent struct {
	CcAddresses         string        `json:"ccAddresses"`
	ClusterId           string        `json:"clusterId"`
	CompletelyRewritten string        `json:"completelyRewritten"`
	FromAddress         string        `json:"fromAddress"`
	GUID                string        `json:"GUID"`
	HeaderFrom          string        `json:"headerFrom"`
	HeaderReplyTo       string        `json:"headerReplyTo"`
	ImposterScore       int           `json:"imposterScore"`
	MalwareScore        int           `json:"malwareScore"`
	MessageID           string        `json:"messageID"`
	MessageParts        []MessagePart `json:"messagePart"`
	MessageSize         int           `json:"messageSize"`
	MessageTime         time.Time     `json:"messageTime"`
	ModulesRun          []string      `json:"modulesRun"`
	PhishScore          int           `json:"phishScore"`
	PolicyRoutes        []string      `json:"policyRoutes"`
	QID                 string        `json:"QID"`
	QuarantineFolder    string        `json:"quarantineFolder"`
	QuarantineRule      string        `json:"quarantineRule"`
	Recipient           string        `json:"recipient"`
	ReplyToAddress      string        `json:"replyToAddress"`
	Sender              string        `json:"sender"`
	SenderIP            string        `json:"senderIP"`
	SpamScore           int           `json:"spamScore"`
	Subject             string        `json:"subject"`
	ThreatsInfoMap      []ThreatInfo  `json:"threatsInfoMap"`
}

type ThreatInfo struct {
	DetectionType  string    `json:"detectionType"`
	CampaignID     string    `json:"campaignId"`
	Classification string    `json:"classification"`
	Threat         string    `json:"threat"`
	ThreatID       string    `json:"threatId"`
	ThreatStatus   string    `json:"threatStatus"`
	ThreatTime     time.Time `json:"threatTime"`
	ThreatURL      string    `json:"threatUrl"`
	ToAddresses    []string  `json:"toAddresses"`
	XMailer        string    `json:"xmailer"`
}

type MessagePart struct {
	ContentType   string `json:"contentType"`
	Disposition   string `json:"disposition"`
	Filename      string `json:"filename"`
	MD5           string `json:"md5"`
	OContentType  string `json:"oContentType"`
	SandboxStatus string `json:"sandboxStatus"`
	SHA256        string `json:"sha256"`
}

type ClickEvent struct {
	CampaignID     string    `json:"campaignId"`
	Classification string    `json:"classification"`
	ClickIP        string    `json:"clickIP"`
	ClickTime      time.Time `json:"clickTime"`
	GUID           string    `json:"GUID"`
	ID             string    `json:"id"`
	Recipient      string    `json:"recipient"`
	Sender         string    `json:"sender"`
	SenderIP       string    `json:"senderIP"`
	ThreatID       string    `json:"threatID"`
	ThreatTime     time.Time `json:"threatTime"`
	ThreatURL      string    `json:"threatURL"`
	ThreatStatus   string    `json:"threatStatus"`
	URL            string    `json:"url"`
	UserAgent      string    `json:"userAgent"`
}
