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

package m365receiver // import "github.com/observiq/observiq-otel-collector/receiver/m365receiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/observiq/observiq-otel-collector/receiver/m365receiver/internal/metadata"
)

type reportPair struct {
	columns  map[string]int //relevant metric name and index for csv file
	endpoint string         //unique report endpoint
}

var reports = []reportPair{
	{
		columns:  map[string]int{"sharepoint.files": 2, "sharepoint.files.active": 3},
		endpoint: "getSharePointSiteUsageFileCounts(period='D7')",
	},
	{
		columns:  map[string]int{"sharepoint.sites.active": 3},
		endpoint: "getSharePointSiteUsageSiteCounts(period='D7')",
	},
	{
		columns:  map[string]int{"sharepoint.pages.viewed": 2},
		endpoint: "getSharePointSiteUsagePages(period='D7')",
	},
	{
		columns:  map[string]int{"sharepoint.pages.unique": 1},
		endpoint: "getSharePointActivityPages(period='D7')",
	},
	{
		columns:  map[string]int{"sharepoint.site.storage": 2},
		endpoint: "getSharePointSiteUsageStorage(period='D7')",
	},
	{
		columns:  map[string]int{"device.web": 1, "device.android": 3, "device.ios": 4, "device.mac": 5, "device.windows": 6, "device.chromeos": 7, "device.linux": 8},
		endpoint: "getTeamsDeviceUsageDistributionUserCounts(period='D7')",
	},
	{
		columns:  map[string]int{"teams.message.team": 2, "teams.message.private": 5, "teams.calls": 6, "teams.meetings": 7},
		endpoint: "getTeamsUserActivityCounts(period='D7')",
	},
	{
		columns:  map[string]int{"onedrive.files": 2, "onedrive.files.active": 3},
		endpoint: "getOneDriveUsageFileCounts(period='D7')",
	},
	{
		columns:  map[string]int{"onedrive.activity.view_edit": 1, "onedrive.activity.synced": 2, "onedrive.activity.internal": 3, "onedrive.activity.external": 4},
		endpoint: "getOneDriveActivityUserCounts(period='D7')",
	},
	{
		columns:  map[string]int{"outlook.mailboxes.active": 2},
		endpoint: "getMailboxUsageMailboxCounts(period='D7')",
	},
	{
		columns:  map[string]int{"outlook.read": 3, "outlook.sent": 1, "outlook.received": 2},
		endpoint: "getEmailActivityCounts(period='D7')",
	},
	{
		columns:  map[string]int{"outlook.storage": 1},
		endpoint: "getMailboxUsageStorage(period='D7')",
	},
	{
		columns:  map[string]int{"outlook.pop3": 7, "outlook.imap4": 8, "outlook.smtp": 9, "outlook.windows": 3, "outlook.mac": 2, "outlook.web": 6, "outlook.mobile": 4, "outlook.other_mobile": 5},
		endpoint: "getEmailAppUsageAppsUserCounts(period='D7')",
	},
	{
		columns:  map[string]int{"outlook.under_limit": 1, "outlook.warning": 2, "outlook.send_prohibited": 3, "outlook.send_receive_prohibited": 4, "outlook.indeterminate": 5},
		endpoint: "getMailboxUsageQuotaStatusMailboxCounts(period='D7')",
	},
}

type mClient interface {
	GetCSV(endpoint string) ([]string, error)
	GetToken() error
	shutdown() error
}

type m365Scraper struct {
	settings component.TelemetrySettings
	logger   *zap.Logger
	cfg      *Config
	client   mClient
	mb       *metadata.MetricsBuilder
	root     string
}

func newM365Scraper(
	settings receiver.CreateSettings,
	cfg *Config,
) *m365Scraper {
	m := &m365Scraper{
		settings: settings.TelemetrySettings,
		logger:   settings.Logger,
		cfg:      cfg,
		mb:       metadata.NewMetricsBuilder(cfg.MetricsBuilderConfig, settings),
		root:     "https://graph.microsoft.com/v1.0/reports/",
	}
	return m
}

func (m *m365Scraper) start(_ context.Context, host component.Host) error {
	httpClient, err := m.cfg.ToClient(host, m.settings)
	if err != nil {
		m.logger.Error("error creating HTTP client", zap.Error(err))
		return err
	}

	m.client = newM365Client(httpClient, m.cfg)

	err = m.client.GetToken()
	if err != nil {
		m.logger.Error("error creating authorization token", zap.Error(err))
		return err
	}

	return nil
}

func (m *m365Scraper) shutdown(_ context.Context) error {
	return m.client.shutdown()
}

// retrieves data, builds metrics & emits them
func (m *m365Scraper) scrape(context.Context) (pmetric.Metrics, error) {
	m365Data, err := m.getStats()
	if err != nil {
		//troubleshoot stale token
		m.logger.Error("error retrieving stats", zap.Error(err))
		if err.Error() == "access token invalid" {
			m.logger.Error("possible stale token; attempting to regenerate")
			err = m.client.GetToken()
			if err != nil {
				//something went wrong with generating token.
				m.logger.Error("error creating authorization token", zap.Error(err))
				return pmetric.Metrics{}, err
			}
			//retry data retrieval with fresh token
			m365Data, err = m.getStats()
			if err != nil {
				//not an error with the access token
				m.logger.Error("unable to retrieve stats", zap.Error(err))
				return pmetric.Metrics{}, err
			}
		} else {
			//not an error with the access token
			return pmetric.Metrics{}, err
		}
	}

	now := pcommon.NewTimestampFromTime(time.Now())
	for metricKey, metricValue := range m365Data {
		switch metricKey {
		case "sharepoint.files":
			m.mb.RecordM365SharepointFilesCountDataPoint(now, metricValue)
		case "sharepoint.files.active":
			m.mb.RecordM365SharepointFilesActiveCountDataPoint(now, metricValue)
		case "sharepoint.sites.active":
			m.mb.RecordM365SharepointSitesActiveCountDataPoint(now, metricValue)
		case "sharepoint.pages.viewed":
			m.mb.RecordM365SharepointPagesViewedCountDataPoint(now, metricValue)
		case "sharepoint.pages.unique":
			m.mb.RecordM365SharepointPagesUniqueCountDataPoint(now, metricValue)
		case "sharepoint.site.storage":
			m.mb.RecordM365SharepointSiteStorageCountDataPoint(now, metricValue)
		case "device.web":
			m.mb.RecordM365TeamsDeviceUsageCountDataPoint(now, metricValue, metadata.AttributeTeamsDevicesWeb)
		case "device.android":
			m.mb.RecordM365TeamsDeviceUsageCountDataPoint(now, metricValue, metadata.AttributeTeamsDevicesAndroid)
		case "device.ios":
			m.mb.RecordM365TeamsDeviceUsageCountDataPoint(now, metricValue, metadata.AttributeTeamsDevicesIOS)
		case "device.mac":
			m.mb.RecordM365TeamsDeviceUsageCountDataPoint(now, metricValue, metadata.AttributeTeamsDevicesMac)
		case "device.windows":
			m.mb.RecordM365TeamsDeviceUsageCountDataPoint(now, metricValue, metadata.AttributeTeamsDevicesWindows)
		case "device.chromeos":
			m.mb.RecordM365TeamsDeviceUsageCountDataPoint(now, metricValue, metadata.AttributeTeamsDevicesChromeOS)
		case "device.linux":
			m.mb.RecordM365TeamsDeviceUsageCountDataPoint(now, metricValue, metadata.AttributeTeamsDevicesLinux)
		case "teams.message.team":
			m.mb.RecordM365TeamsMessageTeamCountDataPoint(now, metricValue)
		case "teams.message.private":
			m.mb.RecordM365TeamsMessagesPrivateCountDataPoint(now, metricValue)
		case "teams.calls":
			m.mb.RecordM365TeamsCallsCountDataPoint(now, metricValue)
		case "teams.meetings":
			m.mb.RecordM365TeamsMeetingsCountDataPoint(now, metricValue)
		case "onedrive.files":
			m.mb.RecordM365OnedriveFilesCountDataPoint(now, metricValue)
		case "onedrive.files.active":
			m.mb.RecordM365OnedriveFilesActiveCountDataPoint(now, metricValue)
		case "onedrive.activity.view_edit":
			m.mb.RecordM365OnedriveUserActivityCountDataPoint(now, metricValue, metadata.AttributeOnedriveActivityViewEdit)
		case "onedrive.activity.synced":
			m.mb.RecordM365OnedriveUserActivityCountDataPoint(now, metricValue, metadata.AttributeOnedriveActivitySynced)
		case "onedrive.activity.internal":
			m.mb.RecordM365OnedriveUserActivityCountDataPoint(now, metricValue, metadata.AttributeOnedriveActivityInternalShare)
		case "onedrive.activity.external":
			m.mb.RecordM365OnedriveUserActivityCountDataPoint(now, metricValue, metadata.AttributeOnedriveActivityExternalShare)
		case "outlook.mailboxes.active":
			m.mb.RecordM365OutlookMailboxesActiveCountDataPoint(now, metricValue)
		case "outlook.read":
			m.mb.RecordM365OutlookEmailActivityCountDataPoint(now, metricValue, metadata.AttributeOutlookActivityRead)
		case "outlook.sent":
			m.mb.RecordM365OutlookEmailActivityCountDataPoint(now, metricValue, metadata.AttributeOutlookActivitySent)
		case "outlook.received":
			m.mb.RecordM365OutlookEmailActivityCountDataPoint(now, metricValue, metadata.AttributeOutlookActivityReceived)
		case "outlook.storage":
			m.mb.RecordM365OutlookStorageCountDataPoint(now, metricValue)
		case "outlook.pop3":
			m.mb.RecordM365OutlookAppUserCountDataPoint(now, metricValue, metadata.AttributeOutlookAppsPop3)
		case "outlook.imap4":
			m.mb.RecordM365OutlookAppUserCountDataPoint(now, metricValue, metadata.AttributeOutlookAppsImap4)
		case "outlook.smtp":
			m.mb.RecordM365OutlookAppUserCountDataPoint(now, metricValue, metadata.AttributeOutlookAppsSmtp)
		case "outlook.windows":
			m.mb.RecordM365OutlookAppUserCountDataPoint(now, metricValue, metadata.AttributeOutlookAppsWindows)
		case "outlook.mac":
			m.mb.RecordM365OutlookAppUserCountDataPoint(now, metricValue, metadata.AttributeOutlookAppsMac)
		case "outlook.web":
			m.mb.RecordM365OutlookAppUserCountDataPoint(now, metricValue, metadata.AttributeOutlookAppsWeb)
		case "outlook.mobile":
			m.mb.RecordM365OutlookAppUserCountDataPoint(now, metricValue, metadata.AttributeOutlookAppsMobile)
		case "outlook.other_mobile":
			m.mb.RecordM365OutlookAppUserCountDataPoint(now, metricValue, metadata.AttributeOutlookAppsOtherMobile)
		case "outlook.under_limit":
			m.mb.RecordM365OutlookQuotaStatusCountDataPoint(now, metricValue, metadata.AttributeOutlookQuotasUnderLimit)
		case "outlook.warning":
			m.mb.RecordM365OutlookQuotaStatusCountDataPoint(now, metricValue, metadata.AttributeOutlookQuotasWarning)
		case "outlook.send_prohibited":
			m.mb.RecordM365OutlookQuotaStatusCountDataPoint(now, metricValue, metadata.AttributeOutlookQuotasSendProhibited)
		case "outlook.send_receive_prohibited":
			m.mb.RecordM365OutlookQuotaStatusCountDataPoint(now, metricValue, metadata.AttributeOutlookQuotasSendReceiveProhibited)
		case "outlook.indeterminate":
			m.mb.RecordM365OutlookQuotaStatusCountDataPoint(now, metricValue, metadata.AttributeOutlookQuotasIndeterminate)
		}
	}

	return m.mb.Emit(), nil
}

func (m *m365Scraper) getStats() (map[string]string, error) {
	reportData := map[string]string{}

	for _, r := range reports {
		line, err := m.client.GetCSV(m.root + r.endpoint)
		if err != nil {
			return map[string]string{}, err
		}
		if len(line) == 0 {
			m.logger.Sugar().Errorf("no data available from %s endpoint and associated metrics", r.endpoint, zap.Error(err))
			continue
		}

		for k, v := range r.columns {
			reportData[k] = line[v]
		}
	}

	return reportData, nil
}
