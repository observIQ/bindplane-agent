package cloudwatch

import (
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-log-collection/operator/helper"
	"github.com/open-telemetry/opentelemetry-log-collection/testutil"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	basicConfig := func() *CloudwatchInputConfig {
		cfg := NewCloudwatchConfig("test_operator_id")
		cfg.Region = "test-region"
		return cfg
	}

	var stream = "test stream"
	var testStreams = []*string{&stream}
	cases := []struct {
		name      string
		input     *CloudwatchInputConfig
		expectErr bool
	}{
		{
			"default",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				return cfg
			}(),
			false,
		},
		{
			"log-stream-name-prefix",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.LogStreamNamePrefix = ""
				return cfg
			}(),
			false,
		},
		{
			"event-limit",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.EventLimit = 5000
				return cfg
			}(),
			false,
		},
		{
			"poll-interval",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.PollInterval = helper.Duration{Duration: 15 * time.Second}
				return cfg
			}(),
			false,
		},
		{
			"profile",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.Profile = "test"
				return cfg
			}(),
			false,
		},
		{
			"log-stream-names",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.LogStreamNames = testStreams
				return cfg
			}(),
			false,
		},
		{
			"startat-end",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.StartAt = "end"
				return cfg
			}(),
			false,
		},
		{
			"logStreamNames and logStreamNamePrefix both parameters Error",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.LogStreamNames = testStreams
				cfg.LogStreamNamePrefix = "some prefix"
				return cfg
			}(),
			true,
		},
		{
			"startat-beginning",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.StartAt = "beginning"
				cfg.LogStreamNamePrefix = "some prefix"
				return cfg
			}(),
			false,
		},
		{
			"poll-interval-invalid",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.PollInterval = helper.Duration{Duration: time.Second * 0}
				return cfg
			}(),
			true,
		},
		{
			"event-limit-invalid",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.EventLimit = 10001
				return cfg
			}(),
			true,
		},
		{
			"default-required-startat-invalid",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				cfg.StartAt = "invalid"
				return cfg
			}(),
			true,
		},
		{
			"log-group-name",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroupName = "test"
				return cfg
			}(),
			false,
		},
		{
			"log-groups",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				return cfg
			}(),
			false,
		},
		{
			"log-groups-and-log-group-name",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroupName = "test"
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				return cfg
			}(),
			false,
		},
		{
			"log-group-prefix",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroupPrefix = "/aws"
				return cfg
			}(),
			false,
		},
		{
			"log-group-prefix-and-log-groups",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroupPrefix = "/aws"
				cfg.LogGroups = []string{
					"test",
					"test-2",
				}
				return cfg
			}(),
			false,
		},
		{
			"log-group-prefix-and-log-group-name",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroupPrefix = "/aws"
				cfg.LogGroupName = "test"
				return cfg
			}(),
			false,
		},
		{
			"log_group_prefix-log_group_name-log_groups",
			func() *CloudwatchInputConfig {
				cfg := basicConfig()
				cfg.LogGroupPrefix = "/aws"
				cfg.LogGroupName = "test"
				cfg.LogGroups = []string{
					"test",
					"aws",
				}
				return cfg
			}(),
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := tc.input
			_, err := cfg.Build(testutil.NewBuildContext(t))
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCurrentTimeInUnixMilliseconds(t *testing.T) {
	timeNow := time.Now()
	timeNowUnixMillis := timeNow.UnixNano() / int64(time.Millisecond)
	cases := []struct {
		name     string
		input    time.Time
		expected int64
	}{
		{
			name:     "test",
			input:    timeNow,
			expected: timeNowUnixMillis,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			x := currentTimeInUnixMilliseconds(tc.input)
			require.Equal(t, tc.expected, x)
		})
	}
}

func TestFromUnixMilli(t *testing.T) {
	timeNow := time.Now()
	timeNowUnixMillis := currentTimeInUnixMilliseconds(timeNow)

	cases := []struct {
		name     string
		input    int64
		expected time.Time
	}{
		{
			name:     "Time Now",
			input:    timeNowUnixMillis,
			expected: timeNow,
		},
		{
			name:     "Specific Time",
			input:    1620842185279,
			expected: time.Unix(0, 1620842185279000000),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			x := fromUnixMilli(tc.input)
			require.Equal(t, tc.expected.Unix(), x.Unix())
		})
	}
}

func TestTimeLayoutParser(t *testing.T) {
	timeNow := time.Now().Unix()
	cases := []struct {
		name      string
		input     string
		timeToUse int64
		expected  string
	}{
		{
			name:      "Time Now",
			input:     "%Y/%m/%d",
			timeToUse: timeNow,
			expected:  time.Unix(timeNow, 0).Format("2006/01/02"),
		},
		{
			name:      "Year4Digigt-Month2Digit-Day2Digit",
			input:     "%Y-%m-%d",
			timeToUse: 1620843711,
			expected:  "2021-05-12",
		},
		{
			name:      "Year4Digigt-Month2Digit-Day2Digit-TrailingText",
			input:     "%Y-%m-%d/Test",
			timeToUse: 1620843711,
			expected:  "2021-05-12/Test",
		},
		{
			name:      "Layout repeated",
			input:     "%Y-%m-%d %Y-%m-%d",
			timeToUse: 1620843711,
			expected:  "2021-05-12 %Y-%m-%d",
		},
		{
			name:      "All Directives",
			input:     "%Y-%y-%m-%q-%b-%h-%B-%d-%g-%a-%A",
			timeToUse: 1639351311,
			expected:  "2021-21-12-12-Dec-Dec-December-12-12-Sun-Sunday",
		},
		{
			name:      "All Directives single digit day",
			input:     "%Y-%y-%m-%q-%b-%h-%B-%d-%g-%a-%A",
			timeToUse: 1619907711,
			expected:  "2021-21-05-5-May-May-May-01-1-Sat-Saturday",
		},
		{
			name:      "All Directives single digit month",
			input:     "%Y-%y-%m-%q-%b-%h-%B-%d-%g-%a-%A",
			timeToUse: 1620858111,
			expected:  "2021-21-05-5-May-May-May-12-12-Wed-Wednesday",
		},
		{
			name:      "Leap Year",
			input:     "%Y-%y-%m-%q-%b-%h-%B-%d-%g-%a-%A",
			timeToUse: 1583018511,
			expected:  "2020-20-02-2-Feb-Feb-February-29-29-Sat-Saturday",
		},
		{
			name:      "No Directives",
			input:     "2021-05-12",
			timeToUse: 1583018511,
			expected:  "2021-05-12",
		},
		{
			name:      "Empty string",
			input:     "",
			timeToUse: 1583018511,
			expected:  "",
		},
		{
			name:      "Symbols",
			input:     "%^&*!@#$()-=+_",
			timeToUse: 1583018511,
			expected:  "%^&*!@#$()-=+_",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, timeLayoutParser(tc.input, time.Unix(tc.timeToUse, 0)))
		})
	}
}
