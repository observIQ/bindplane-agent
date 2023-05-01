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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetToken(t *testing.T) {
	m365Mock := newMockServerToken()
	testClient := newM365Client(m365Mock.Client(), &Config{})
	testClient.authEndpoint = m365Mock.URL + "/testTenantID"
	testClient.clientID = "testClientID"
	testClient.clientSecret = "testClientSecret"

	err := testClient.GetToken()
	require.NoError(t, err)
	require.Equal(t, "testAccessToken", testClient.token)

	//err testing
	testClient.clientSecret = "err"
	err = testClient.GetToken()
	assert.EqualError(t, err, "got non 200 status code from request, got 400")

}

func TestGetCSV(t *testing.T) {
	m365Mock := newMockServerCSV()
	testClient := newM365Client(m365Mock.Client(), &Config{})
	testClient.token = "foo"

	//expected behavior
	testLine, err := testClient.GetCSV(m365Mock.URL + "/getSharePointSiteUsageFileCounts(period='D7')")
	require.NoError(t, err)
	expectedLine := []string{"2023-04-25", "All", "2", "0", "2023-04-25", "7"}
	require.Equal(t, expectedLine, testLine)

	//test no returned data
	testLine, err = testClient.GetCSV(m365Mock.URL + "/testNoData")
	require.NoError(t, err)
	expectedLine = []string{}
	require.Equal(t, expectedLine, testLine)

	//err testing
	testClient.token = "err"
	_, err = testClient.GetCSV(m365Mock.URL + "/getSharePointSiteUsageFileCounts(period='D7')")
	assert.EqualError(t, err, "got non 200 status code from request, got 400")
}

//	Mock Servers

func newMockServerCSV() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/getSharePointSiteUsageFileCounts(period='D7')" {
			if req.Method != "GET" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, incorrect HTTP method"))
				return
			}

			if a := req.Header.Get("Authorization"); a != "foo" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, not authorized"))
				return
			}

			rw.WriteHeader(200)
			rw.Write([]byte(
				"Report Refresh Date,Site Type,Total,Active,Report Date,Report Period\n2023-04-25,All,2,0,2023-04-25,7\n2023-04-25,All,2,0,2023-04-24,7\n2023-04-25,All,2,0,2023-04-23,7\n2023-04-25,All,2,0,2023-04-22,7\n2023-04-25,All,2,0,2023-04-21,7\n2023-04-25,All,2,0,2023-04-20,7\n2023-04-25,All,2,0,2023-04-19,7\n",
			))
			return
		}
		if req.URL.String() == "/testNoData" {
			if req.Method != "GET" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, incorrect HTTP method"))
				return
			}

			if a := req.Header.Get("Authorization"); a != "foo" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, not authorized"))
				return
			}

			rw.WriteHeader(200)
			rw.Write([]byte(
				"Report Refresh Date,Site Type,Total,Active,Report Date,Report Period\n",
			))
			return
		}
		rw.WriteHeader(404)
		return
	}))
}

func newMockServerToken() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/testTenantID" {
			if req.Method != "POST" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, incorrect HTTP method"))
			}

			req.ParseForm()
			if gType := req.Form.Get("grant_type"); gType != "client_credentials" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, incorrect grant_type"))
				return
			}
			if scope := req.Form.Get("scope"); scope != "https://graph.microsoft.com/.default" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, incorrect scope"))
				return
			}
			if cID := req.Form.Get("client_id"); cID != "testClientID" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, incorrect ID"))
				return
			}
			if cSec := req.Form.Get("client_secret"); cSec != "testClientSecret" {
				rw.WriteHeader(400)
				rw.Write([]byte("Error, incorrect secret"))
				return
			}

			rw.WriteHeader(200)
			rw.Write([]byte(
				`{
					"token_type": "Bearer",
					"expires_in": 3599,
					"ext_expires_in": 3599,
					"access_token": "testAccessToken"
				}`))
			return
		}
		rw.WriteHeader(404)
		return
	}))
}
