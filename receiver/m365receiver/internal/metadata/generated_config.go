// Code generated by mdatagen. DO NOT EDIT.

package metadata

import "go.opentelemetry.io/collector/confmap"

// MetricConfig provides common config for a particular metric.
type MetricConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (ms *MetricConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(ms, confmap.WithErrorUnused())
	if err != nil {
		return err
	}
	ms.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// MetricsConfig provides config for m365 metrics.
type MetricsConfig struct {
	M365OnedriveFilesActiveCount    MetricConfig `mapstructure:"m365.onedrive.files.active.count"`
	M365OnedriveFilesCount          MetricConfig `mapstructure:"m365.onedrive.files.count"`
	M365OnedriveUserActivityCount   MetricConfig `mapstructure:"m365.onedrive.user_activity.count"`
	M365OutlookAppUserCount         MetricConfig `mapstructure:"m365.outlook.app.user.count"`
	M365OutlookEmailActivityCount   MetricConfig `mapstructure:"m365.outlook.email_activity.count"`
	M365OutlookMailboxesActiveCount MetricConfig `mapstructure:"m365.outlook.mailboxes.active.count"`
	M365OutlookQuotaStatusCount     MetricConfig `mapstructure:"m365.outlook.quota_status.count"`
	M365OutlookStorageUsed          MetricConfig `mapstructure:"m365.outlook.storage.used"`
	M365SharepointFilesActiveCount  MetricConfig `mapstructure:"m365.sharepoint.files.active.count"`
	M365SharepointFilesCount        MetricConfig `mapstructure:"m365.sharepoint.files.count"`
	M365SharepointPagesUniqueCount  MetricConfig `mapstructure:"m365.sharepoint.pages.unique.count"`
	M365SharepointPagesViewedCount  MetricConfig `mapstructure:"m365.sharepoint.pages.viewed.count"`
	M365SharepointSiteStorageUsed   MetricConfig `mapstructure:"m365.sharepoint.site.storage.used"`
	M365SharepointSitesActiveCount  MetricConfig `mapstructure:"m365.sharepoint.sites.active.count"`
	M365TeamsCallsCount             MetricConfig `mapstructure:"m365.teams.calls.count"`
	M365TeamsDeviceUsageUsers       MetricConfig `mapstructure:"m365.teams.device_usage.users"`
	M365TeamsMeetingsCount          MetricConfig `mapstructure:"m365.teams.meetings.count"`
	M365TeamsMessagesPrivateCount   MetricConfig `mapstructure:"m365.teams.messages.private.count"`
	M365TeamsMessagesTeamCount      MetricConfig `mapstructure:"m365.teams.messages.team.count"`
}

func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		M365OnedriveFilesActiveCount: MetricConfig{
			Enabled: true,
		},
		M365OnedriveFilesCount: MetricConfig{
			Enabled: true,
		},
		M365OnedriveUserActivityCount: MetricConfig{
			Enabled: true,
		},
		M365OutlookAppUserCount: MetricConfig{
			Enabled: true,
		},
		M365OutlookEmailActivityCount: MetricConfig{
			Enabled: true,
		},
		M365OutlookMailboxesActiveCount: MetricConfig{
			Enabled: true,
		},
		M365OutlookQuotaStatusCount: MetricConfig{
			Enabled: true,
		},
		M365OutlookStorageUsed: MetricConfig{
			Enabled: true,
		},
		M365SharepointFilesActiveCount: MetricConfig{
			Enabled: true,
		},
		M365SharepointFilesCount: MetricConfig{
			Enabled: true,
		},
		M365SharepointPagesUniqueCount: MetricConfig{
			Enabled: true,
		},
		M365SharepointPagesViewedCount: MetricConfig{
			Enabled: true,
		},
		M365SharepointSiteStorageUsed: MetricConfig{
			Enabled: true,
		},
		M365SharepointSitesActiveCount: MetricConfig{
			Enabled: true,
		},
		M365TeamsCallsCount: MetricConfig{
			Enabled: true,
		},
		M365TeamsDeviceUsageUsers: MetricConfig{
			Enabled: true,
		},
		M365TeamsMeetingsCount: MetricConfig{
			Enabled: true,
		},
		M365TeamsMessagesPrivateCount: MetricConfig{
			Enabled: true,
		},
		M365TeamsMessagesTeamCount: MetricConfig{
			Enabled: true,
		},
	}
}

// MetricsBuilderConfig is a configuration for m365 metrics builder.
type MetricsBuilderConfig struct {
	Metrics MetricsConfig `mapstructure:"metrics"`
}

func DefaultMetricsBuilderConfig() MetricsBuilderConfig {
	return MetricsBuilderConfig{
		Metrics: DefaultMetricsConfig(),
	}
}
