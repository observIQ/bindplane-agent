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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/extension/experimental/storage"
	"go.uber.org/zap"
)

// Test the case where some data is exported, but a subsequent call for paginated data fails
func TestSplunkResultsPaginationFailure(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)
	cfg.Searches = []Search{
		{
			Query:          "search index=otel",
			EarliestTime:   "2024-11-14T00:00:00.000Z",
			LatestTime:     "2024-11-14T23:59:59.000Z",
			EventBatchSize: 5,
		},
	}
	var callCount int = 0
	server := newMockSplunkServer(&callCount)
	defer server.Close()
	settings := componenttest.NewNopTelemetrySettings()
	ssapir := newSSAPIReceiver(zap.NewNop(), cfg, settings, component.NewID(typeStr))
	ssapir.client, _ = newSplunkSearchAPIClient(context.Background(), settings, *cfg, componenttest.NewNopHost())
	ssapir.client.(*defaultSplunkSearchAPIClient).client = server.Client()
	ssapir.client.(*defaultSplunkSearchAPIClient).endpoint = server.URL
	ssapir.logsConsumer = &consumertest.LogsSink{}

	ssapir.storageClient = storage.NewNopClient()

	ssapir.initCheckpoint(context.Background())
	ssapir.runQueries(context.Background())
	require.Equal(t, 5, ssapir.checkpointRecord.Offset)
	require.Equal(t, 1, callCount)
}

func newMockSplunkServer(callCount *int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/services/search/jobs" {
			rw.Header().Set("Content-Type", "application/xml")
			rw.WriteHeader(201)
			rw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<response>
				<sid>123456</sid>
			</response>
			`))
		}
		if req.URL.String() == "/services/search/v2/jobs/123456" {
			rw.Header().Set("Content-Type", "application/xml")
			rw.WriteHeader(200)
			rw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<response>
				<content>
					<type>DISPATCH</type>
					<dict>
						<key name="dispatchState">DONE</key>
					</dict>
				</content>
			</response>`))
		}
		if req.URL.String() == "/services/search/v2/jobs/123456/results?output_mode=json&offset=0&count=5" && req.URL.Query().Get("offset") == "0" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(200)
			rw.Write(splunkEventsResultsP1)
			*callCount++
		}
		if req.URL.String() == "/services/search/v2/jobs/123456/results?output_mode=json&offset=5&count=5" && req.URL.Query().Get("offset") == "5" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(400)
			rw.Write([]byte("error, bad request"))
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
