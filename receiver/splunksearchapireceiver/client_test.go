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

package splunksearchapireceiver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	server     = newMockSplunkServer()
	testClient = defaultSplunkSearchAPIClient{
		client:   server.Client(),
		endpoint: server.URL,
	}
)

func TestCreateSearchJob(t *testing.T) {
	// valid search
	resp, err := testClient.CreateSearchJob("index=otel")
	require.NoError(t, err)
	require.Equal(t, "123456", resp.SID)
}

func TestGetJobStatus(t *testing.T) {
	resp, err := testClient.GetJobStatus("123456")
	require.NoError(t, err)
	require.Equal(t, "DONE", resp.Content.Dict.Keys[0].Value)
	require.Equal(t, "text/xml", resp.Content.Type)
}

func TestGetSearchResults(t *testing.T) {

}

// mock Splunk servers
func newMockSplunkServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/services/search/jobs":
			rw.Header().Set("Content-Type", "application/xml")
			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<response>
				<sid>123456</sid>
			</response>
			`))
		case "/services/search/v2/jobs/123456":
			rw.Header().Set("Content-Type", "application/xml")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<response>
				<content type="text/xml">
					<dict>
						<key name="dispatchState">DONE</key>
					</dict>
				</content>
			</response>`))
		case "/services/search/v2/jobs/123456/results":
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(splunkEventsResultsP1)
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}))
}

var splunkEventsResultsP1 = []byte(`{
	"init_offset": 0,
	"results": [
		{
			"_raw": "Hello, world!",
			"_time": "2024-11-14T13:02:31.000-05:00"
		},
		{
			"_raw": "Goodbye, world!",
			"_time": "2024-11-14T13:02:30.000-05:00"
		},
		{
			"_raw": "lorem ipsum",
			"_time": "2024-11-14T13:02:29.000-05:00"
		},
		{
			"_raw": "dolor sit amet",
			"_time": "2024-11-14T13:02:28.000-05:00"
		},
		{
			"_raw": "consectetur adipiscing elit",
			"_time": "2024-11-14T13:02:27.000-05:00"
		}
	]
}`)
