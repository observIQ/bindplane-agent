package splunksearchapireceiver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		endpoint    string
		username    string
		password    string
		searches    []Search
		errExpected bool
		errText     string
	}{
		{
			desc:     "Missing endpoint",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "missing Splunk server endpoint",
		},
		{
			desc:     "Missing username",
			endpoint: "http://localhost:8089",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "missing Splunk username",
		},
		{
			desc:     "Missing password",
			endpoint: "http://localhost:8089",
			username: "user",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "missing Splunk password",
		},
		{
			desc:        "Missing searches",
			endpoint:    "http://localhost:8089",
			username:    "user",
			password:    "password",
			errExpected: true,
			errText:     "at least one search must be provided",
		},
		{
			desc:     "Missing query",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "missing query in search",
		},
		{
			desc:     "Missing earliest_time",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:      "search index=_internal",
					LatestTime: "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "missing earliest_time in search",
		},
		{
			desc:     "Unparsable earliest_time",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "-1hr",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "earliest_time failed to parse as RFC3339",
		},
		{
			desc:     "Missing latest_time",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "missing latest_time in search",
		},
		{
			desc:     "Unparsable latest_time",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "-1hr",
				},
			},
			errExpected: true,
			errText:     "latest_time failed to parse as RFC3339",
		},
		{
			desc:     "earliest_time after latest_time",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T14:00:00.000Z",
					LatestTime:   "2024-10-30T04:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "earliest_time must be earlier than latest_time",
		},
		{
			desc:     "earliest_time and latest_time equal",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T14:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "earliest_time must be earlier than latest_time",
		},
		{
			desc:     "earliest_time in the future",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
					LatestTime:   time.Now().Add(10 * time.Hour).Format(time.RFC3339),
				},
			},
			errExpected: true,
			errText:     "earliest_time must be earlier than current time",
		},
		{
			desc:     "latest_time in the future",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   time.Now().Add(10 * time.Hour).Format(time.RFC3339),
				},
			},
			errExpected: true,
			errText:     "latest_time must be earlier than or equal to current time",
		},
		{
			desc:     "Invalid query chaining",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal | stats count by sourcetype",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: true,
			errText:     "command chaining is not supported for queries",
		},
		{
			desc:     "Valid config",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: false,
		},
		{
			desc:     "Valid config with multiple searches",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
				{
					Query:        "search index=_audit",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
				},
			},
			errExpected: false,
		},
		{
			desc:     "Valid config with limit",
			endpoint: "http://localhost:8089",
			username: "user",
			password: "password",
			searches: []Search{
				{
					Query:        "search index=_internal",
					EarliestTime: "2024-10-30T04:00:00.000Z",
					LatestTime:   "2024-10-30T14:00:00.000Z",
					Limit:        10,
				},
			},
			errExpected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := NewFactory().CreateDefaultConfig().(*Config)
			cfg.Endpoint = tc.endpoint
			cfg.Username = tc.username
			cfg.Password = tc.password
			cfg.Searches = tc.searches
			err := cfg.Validate()
			if tc.errExpected && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.errExpected && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.errExpected {
				require.EqualError(t, err, tc.errText)
				return
			}
			require.NoError(t, err)
		})
	}
}
