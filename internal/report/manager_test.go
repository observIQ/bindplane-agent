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

package report

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagerSetClient(t *testing.T) {
	testCases := []struct {
		desc        string
		client      Client
		expectedErr error
	}{
		{
			desc:        "Nil client",
			client:      nil,
			expectedErr: errors.New("client must not be nil"),
		},
		{
			desc:        "Successful set",
			client:      http.DefaultClient,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m := &Manager{}
			err := m.SetClient(tc.client)
			if tc.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.client, m.client)
			} else {
				assert.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}
