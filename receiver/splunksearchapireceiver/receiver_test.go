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
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

var (
	logger        = zap.NewNop()
	config        = &Config{}
	settings      = component.TelemetrySettings{}
	id            = component.ID{}
	ssapireceiver = newSSAPIReceiver(logger, config, settings, id)
)

func TestPolling(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.JobPollInterval = 1 * time.Second
	ssapireceiver.config = cfg

	client := &mockLogsClient{}
	ssapireceiver.client = client

	file := filepath.Join("testdata", "logs", "testPollJobStatus", "input.xml")
	client.On("GetJobStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(client.loadTestStatusResponse(t, file), nil)
	cancelCtx, cancel := context.WithCancel(context.Background())
	ssapireceiver.cancel = cancel
	ssapireceiver.checkpointRecord = &EventRecord{}

	err := ssapireceiver.pollSearchCompletion(cancelCtx, "123456")
	require.NoError(t, err)

	client.AssertCalled(t, "GetJobStatus", "123456")

	err = ssapireceiver.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestIsSearchCompleted(t *testing.T) {
	jobResponse := SearchStatusResponse{
		Content: Content{
			Dict: Dict{
				Keys: []Key{
					{
						Name:  "dispatchState",
						Value: "DONE",
					},
				},
			},
		},
	}

	emptyResponse := SearchStatusResponse{}

	done := ssapireceiver.isSearchCompleted(jobResponse)
	require.True(t, done)

	jobResponse.Content.Dict.Keys[0].Value = "RUNNING"
	done = ssapireceiver.isSearchCompleted(jobResponse)
	require.False(t, done)

	done = ssapireceiver.isSearchCompleted(emptyResponse)
	require.False(t, done)
}

func TestInitCheckpoint(t *testing.T) {

}

func TestCheckpoint(t *testing.T) {
	t.Skip("Not implemented")
}

func TestLoadCheckpoint(t *testing.T) {
	t.Skip("Not implemented")
}

type mockLogsClient struct {
	mock.Mock
}

func (m *mockLogsClient) loadTestStatusResponse(t *testing.T, file string) SearchStatusResponse {
	logBytes, err := os.ReadFile(file)
	require.NoError(t, err)
	var resp SearchStatusResponse
	err = xml.Unmarshal(logBytes, &resp)
	require.NoError(t, err)
	return resp
}

func (m *mockLogsClient) GetJobStatus(searchID string) (SearchStatusResponse, error) {
	args := m.Called(searchID)
	return args.Get(0).(SearchStatusResponse), args.Error(1)
}

func (m *mockLogsClient) CreateSearchJob(searchQuery string) (CreateJobResponse, error) {
	args := m.Called(searchQuery)
	return args.Get(0).(CreateJobResponse), args.Error(1)
}

func (m *mockLogsClient) GetSearchResults(searchID string, offset int, batchSize int) (SearchResultsResponse, error) {
	args := m.Called(searchID, offset, batchSize)
	return args.Get(0).(SearchResultsResponse), args.Error(1)
}
