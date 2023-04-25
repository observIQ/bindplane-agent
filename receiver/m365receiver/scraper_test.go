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
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-otel-collector/receiver/m365receiver/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestScraper(t *testing.T) {
	//mocks
	mc := &mockClient{}
	mc.On("GetToken").Return(nil)
	mc.On("GetCSV", "getSharePointSiteUsageFileCounts(period='D7')").Return([]string{
		"2023-04-23", "All", "2", "0", "2023-04-23", "7",
	}, nil)
	mc.On("GetCSV", "getSharePointSiteUsageSiteCounts(period='D7')").Return([]string{
		"2023-04-23", "All", "8", "0", "2023-04-23", "7",
	}, nil)

	scraper := newM365Scraper(
		receivertest.NewNopCreateSettings(),
		&Config{MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig()},
	)

	scraper.start(context.Background(), componenttest.NewNopHost())
	scraper.client = mc

	actualMetrics, err := scraper.scrape(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, actualMetrics)

	//generate testdata file
	/*m := pmetric.JSONMarshaler{}
	mBytes, err := m.MarshalMetrics(actualMetrics)
	require.NoError(t, err)
	goldenPath := filepath.Join("testdata", "metrics", "unit-test-metrics.json")
	err = os.WriteFile(goldenPath, mBytes, 0666)
	require.NoError(t, err)*/

	//validate output of scrape
	expectedFile := filepath.Join("testdata", "metrics", "unit-test-metrics.json")
	expectedMetrics, err := ReadMetrics(expectedFile)
	require.NoError(t, err)
	require.NoError(t, pmetrictest.CompareMetrics(expectedMetrics, actualMetrics,
		pmetrictest.IgnoreMetricDataPointsOrder(), pmetrictest.IgnoreStartTimestamp(), pmetrictest.IgnoreTimestamp()),
	)
}

type mockClient struct {
	mock.Mock
}

func (mw *mockClient) GetCSV(endpoint string) ([]string, error) {
	args := mw.Called(endpoint)
	return args.Get(0).([]string), args.Error(1)
}

func (mw *mockClient) GetToken() error {
	args := mw.Called()
	return args.Error(0)
}

func ReadMetrics(filePath string) (pmetric.Metrics, error) {
	expectedFileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return pmetric.Metrics{}, err
	}
	unmarshaller := &pmetric.JSONUnmarshaler{}
	return unmarshaller.UnmarshalMetrics(expectedFileBytes)
}
