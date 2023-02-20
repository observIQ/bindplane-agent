// Copyright  observIQ, Inc.
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

package logdeduplicationprocessor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDefaultProcessorConfig(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	require.Equal(t, defaultInterval, cfg.Interval)
	require.Equal(t, defaultLogCountAttribute, cfg.LogCountAttribute)
	require.Equal(t, defaultTimezone, cfg.Timezone)
	require.Equal(t, defaultMatchFields, cfg.MatchFields)
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		desc        string
		cfg         *Config
		expectedErr error
	}{
		{
			desc: "invalid LogCountAttribute config",
			cfg: &Config{
				LogCountAttribute: "",
				Interval:          defaultInterval,
				Timezone:          defaultTimezone,
				MatchFields:       defaultMatchFields,
			},
			expectedErr: errInvalidLogCountAttribute,
		},
		{
			desc: "invalid Interval config",
			cfg: &Config{
				LogCountAttribute: defaultLogCountAttribute,
				Interval:          -1,
				Timezone:          defaultTimezone,
				MatchFields:       defaultMatchFields,
			},
			expectedErr: errInvalidInterval,
		},
		{
			desc: "invalid Timezone config",
			cfg: &Config{
				LogCountAttribute: defaultLogCountAttribute,
				Interval:          defaultInterval,
				Timezone:          "not a timezone",
				MatchFields:       defaultMatchFields,
			},
			expectedErr: errors.New("timezone is invalid"),
		},
		{
			desc: "invalid match fields config",
			cfg: &Config{
				LogCountAttribute: defaultLogCountAttribute,
				Interval:          defaultInterval,
				Timezone:          defaultTimezone,
				MatchFields:       []string{"not.field"},
			},
			expectedErr: errors.New("a match_field must start with"),
		},
		{
			desc: "valid config",
			cfg: &Config{
				LogCountAttribute: defaultLogCountAttribute,
				Interval:          defaultInterval,
				Timezone:          defaultTimezone,
				MatchFields:       defaultMatchFields,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
