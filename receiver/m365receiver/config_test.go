// Copyright  OpenTelemetry Authors
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

package m365receiver

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		desc         string
		errExpected  bool
		errText      string
		tenantID     string
		clientID     string
		clientSecret string
	}{
		{
			desc:         "expected case, correct",
			errExpected:  false,
			tenantID:     "b243371b-5514-4951-840f-152d525b6555",
			clientID:     "0b110766-c91d-4fe5-bac5-14a16ce3c351",
			clientSecret: "FZu8Q~gFWqIo_vo8tE36TYsymPGk9VfWvTlS7b-I",
		},
		{
			desc:         "incorrect clientID: empty",
			errExpected:  true,
			errText:      "missing client_id; required",
			tenantID:     "b243371b-5514-4951-840f-152d525b6555",
			clientID:     "",
			clientSecret: "FZu8Q~gFWqIo_vo8tE36TYsymPGk9VfWvTlS7b-I",
		},
		{
			desc:         "incorrect clientID: too short",
			errExpected:  true,
			errText:      "client_id is invalid; must be a GUID",
			tenantID:     "b243371b-5514-4951-840f-152d525b6555",
			clientID:     "foo-test",
			clientSecret: "FZu8Q~gFWqIo_vo8tE36TYsymPGk9VfWvTlS7b-I",
		},
		//todo: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := NewFactory().CreateDefaultConfig().(*Config)
			cfg.Tenant_id = tc.tenantID
			cfg.Client_id = tc.clientID
			cfg.Client_secret = tc.clientSecret

			err := component.ValidateConfig(cfg)

			if tc.errExpected {
				require.EqualError(t, err, tc.errText)
				return
			}

			require.NoError(t, err)
		})
	}
}
