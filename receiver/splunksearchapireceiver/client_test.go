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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
	resp, err := testClient.CreateSearchJob("index=otel")
	require.NoError(t, err)
	require.Equal(t, "123456", resp.SID)

	// returns an error if the response status isn't 201
	resp, err = testClient.CreateSearchJob("index=fail_to_create_job")
	require.ErrorContains(t, err, "failed to create search job")
	require.Empty(t, resp)

	// returns an error if the response body can't be unmarshalled
	resp, err = testClient.CreateSearchJob("index=fail_to_unmarshal")
	require.ErrorContains(t, err, "failed to unmarshal search job create response")
	require.Empty(t, resp)

}

func TestGetJobStatus(t *testing.T) {
	resp, err := testClient.GetJobStatus("123456")
	require.NoError(t, err)
	require.Equal(t, "DONE", resp.Content.Dict.Keys[0].Value)
	require.Equal(t, "text/xml", resp.Content.Type)

	// returns an error if the response status isn't 200
	resp, err = testClient.GetJobStatus("654321")
	require.ErrorContains(t, err, "failed to get search job status")
	require.Empty(t, resp)

	// returns an error if the response body can't be unmarshalled
	resp, err = testClient.GetJobStatus("098765")
	require.ErrorContains(t, err, "failed to unmarshal search job status response")
	require.Empty(t, resp)
}

func TestGetSearchResults(t *testing.T) {
	resp, err := testClient.GetSearchResults("123456", 0, 5)
	require.NoError(t, err)
	require.Equal(t, 5, len(resp.Results))
	require.Equal(t, "Hello, world!", resp.Results[0].Raw)

	// returns an error if the response status isn't 200
	resp, err = testClient.GetSearchResults("654321", 0, 5)
	require.ErrorContains(t, err, "failed to get search job results")
	require.Empty(t, resp)

	// returns an error if the response body can't be unmarshalled
	resp, err = testClient.GetSearchResults("098765", 0, 5)
	require.ErrorContains(t, err, "failed to unmarshal search job results response")
	require.Empty(t, resp)
}

// mock Splunk servers
func newMockSplunkServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.String() {
		case "/services/search/jobs":
			body, _ := io.ReadAll(req.Body)
			if strings.Contains(string(body), "index=otel") {
				rw.Header().Set("Content-Type", "application/xml")
				rw.WriteHeader(http.StatusCreated)
				rw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
				<response>
					<sid>123456</sid>
				</response>
				`))
			}
			if strings.Contains(string(body), "index=fail_to_create_job") {
				rw.WriteHeader(http.StatusNotFound)
			}
			if strings.Contains(string(body), "index=fail_to_unmarshal") {
				rw.WriteHeader(http.StatusCreated)
				rw.Write([]byte(`invalid xml`))
				req.Body = &errorReader{}
			}
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
		case "/services/search/v2/jobs/654321":
			rw.WriteHeader(http.StatusNotFound)
		case "/services/search/v2/jobs/098765":
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`invalid xml`))
		case "/services/search/v2/jobs/123456/results?output_mode=json&offset=0&count=5":
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(splunkEventsResultsP1)
		case "/services/search/v2/jobs/654321/results?output_mode=json&offset=0&count=5":
			rw.WriteHeader(http.StatusNotFound)
		case "/services/search/v2/jobs/098765/results?output_mode=json&offset=0&count=5":
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`invalid json`))
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

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (e *errorReader) Close() error {
	return nil
}
