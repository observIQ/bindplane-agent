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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetToken(t *testing.T) {
	m365Mock := newMockServerToken(t)
	testClient := newM365Client(m365Mock.Client(), &Config{})
	testClient.authEndpoint = m365Mock.URL + "/testTenantID"
	testClient.clientID = "testClientID"
	testClient.clientSecret = "testClientSecret"

	err := testClient.GetToken()
	require.NoError(t, err)
	require.Equal(t, "testAccessToken", testClient.token)
}

func TestGetCSV(t *testing.T) {
	m365Mock := newMockServerCSV(t)
	testClient := newM365Client(m365Mock.Client(), &Config{})
	testClient.token = "foo"

	testLine, err := testClient.GetCSV(m365Mock.URL + "/getSharePointSiteUsageFileCounts(period='D7')")
	require.NoError(t, err)
	expectedLine := []string{
		"2023-04-25", "All", "2", "0", "2023-04-25", "7",
	}
	require.Equal(t, expectedLine, testLine)
}

func newMockServerCSV(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/getSharePointSiteUsageFileCounts(period='D7')" {
			if req.Method != "GET" {
				t.Errorf("expected GET request, got %s", req.Method)
			}

			if a := req.Header.Get("Authorization"); a != "foo" {
				t.Errorf("incorrect authorization token, expected 'foo' got %s", a)
			}

			rw.WriteHeader(200)
			_, err := rw.Write([]byte(
				"Report Refresh Date,Site Type,Total,Active,Report Date,Report Period\n2023-04-25,All,2,0,2023-04-25,7\n2023-04-25,All,2,0,2023-04-24,7\n2023-04-25,All,2,0,2023-04-23,7\n2023-04-25,All,2,0,2023-04-22,7\n2023-04-25,All,2,0,2023-04-21,7\n2023-04-25,All,2,0,2023-04-20,7\n2023-04-25,All,2,0,2023-04-19,7\n",
			))
			require.NoError(t, err)
			return
		}
		rw.WriteHeader(404)
	}))
}

func newMockServerToken(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/testTenantID" {
			if req.Method != "POST" {
				t.Errorf("expected POST request, got %s", req.Method)
			}

			req.ParseForm()
			if gType := req.Form.Get("grant_type"); gType != "client_credentials" {
				t.Errorf("Expected request to have 'grant_type=client_credentials', got: %s", gType)
			}
			if scope := req.Form.Get("scope"); scope != "https://graph.microsoft.com/.default" {
				t.Errorf("Expected request to have 'scope=https://graph.microsoft.com/.default', got %s", scope)
			}
			if cID := req.Form.Get("client_id"); cID != "testClientID" {
				t.Errorf("Expected request to have 'client_id=testClientID', got %s", cID)
			}
			if cSec := req.Form.Get("client_secret"); cSec != "testClientSecret" {
				t.Errorf("Expected request to have 'client_secret=testClientSecret', got %s", cSec)
			}

			rw.WriteHeader(200)
			_, err := rw.Write([]byte(
				`{
					"token_type": "Bearer",
					"expires_in": 3599,
					"ext_expires_in": 3599,
					"access_token": "testAccessToken"
				 }`))
			require.NoError(t, err)
			return
		}
		rw.WriteHeader(404)
	}))
}

// r.ParseForm()
//     topic := r.Form.Get("topic")
//     if topic != "meaningful-topic" {
//       t.Errorf("Expected request to have ‘topic=meaningful-topic’, got: ‘%s’", topic)
//     }
