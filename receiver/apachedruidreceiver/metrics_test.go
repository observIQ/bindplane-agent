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

package apachedruidreceiver

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"testing"

	"github.com/observiq/observiq-otel-collector/receiver/apachedruidreceiver/internal/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap/zaptest"
)

const (
	VALID_USERNAME      = "john.doe"
	VALID_PASSWORD      = "abcd1234"
	ENCODED_CREDENTIALS = "Basic am9obi5kb2U6YWJjZDEyMzQ="
)

func TestHandleRequest(t *testing.T) {
	cases := []struct {
		name               string
		request            *http.Request
		expectedStatusCode int
		metricExpected     bool
		consumerFailure    bool
		configBasicAuth    *BasicAuth
	}{
		{
			name: "No basic auth provided when configured to have it",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBufferString(`{"ClientIP": "127.0.0.1"}`)),
			},
			expectedStatusCode: http.StatusUnauthorized,
			metricExpected:     false,
			consumerFailure:    false,
			configBasicAuth: &BasicAuth{
				Username: VALID_USERNAME,
				Password: VALID_PASSWORD,
			},
		},
		{
			name: "Invalid basic auth",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBufferString(`{"ClientIP": "127.0.0.1"}`)),
				Header: map[string][]string{
					textproto.CanonicalMIMEHeaderKey("Authorization"): {"abc123"},
				},
			},
			expectedStatusCode: http.StatusUnauthorized,
			metricExpected:     false,
			consumerFailure:    false,
			configBasicAuth: &BasicAuth{
				Username: VALID_USERNAME,
				Password: VALID_PASSWORD,
			},
		},
		{
			name: "Non-POST request",
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBufferString(`{"ClientIP": "127.0.0.1"}`)),
				Header: map[string][]string{
					textproto.CanonicalMIMEHeaderKey("Authorization"): {ENCODED_CREDENTIALS},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			metricExpected:     false,
			consumerFailure:    false,
			configBasicAuth: &BasicAuth{
				Username: VALID_USERNAME,
				Password: VALID_PASSWORD,
			},
		},
		{
			name: "Non-json content",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBufferString(`{"ClientIP": "127.0.0.1"}`)),
				Header: map[string][]string{
					textproto.CanonicalMIMEHeaderKey("Authorization"): {ENCODED_CREDENTIALS},
					textproto.CanonicalMIMEHeaderKey("Content-Type"):  {"text/html"},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			metricExpected:     false,
			consumerFailure:    false,
			configBasicAuth: &BasicAuth{
				Username: VALID_USERNAME,
				Password: VALID_PASSWORD,
			},
		},
		{
			name: "Bad JSON",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8" standalone="no" ?>`)),
				Header: map[string][]string{
					textproto.CanonicalMIMEHeaderKey("Authorization"): {ENCODED_CREDENTIALS},
					textproto.CanonicalMIMEHeaderKey("Content-Type"):  {"application/json"},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			metricExpected:     false,
			consumerFailure:    false,
			configBasicAuth: &BasicAuth{
				Username: VALID_USERNAME,
				Password: VALID_PASSWORD,
			},
		},
		{
			name: "Consumer failure",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBufferString(`[{"ClientIP": "127.0.0.1"}]`)),
				Header: map[string][]string{
					textproto.CanonicalMIMEHeaderKey("Authorization"): {ENCODED_CREDENTIALS},
					textproto.CanonicalMIMEHeaderKey("Content-Type"):  {"application/json"},
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			metricExpected:     false,
			consumerFailure:    true,
			configBasicAuth: &BasicAuth{
				Username: VALID_USERNAME,
				Password: VALID_PASSWORD,
			},
		},
		{
			name: "Request succeeds",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBufferString(`[{"service":"druid/broker","metric":"query/count","value":123},{"service":"druid/broker","metric":"query/success/count","value":101},{"service":"druid/broker","metric":"query/failed/count","value":12},{"service":"druid/broker","metric":"query/interrupted/count","value":6},{"service":"druid/broker","metric":"query/timeout/count","value":4},{"service":"druid/broker","metric":"sqlQuery/time","dataSource":"table_1","value":97},{"service":"druid/broker","metric":"sqlQuery/bytes","dataSource":"table_1","value":450},{"service":"druid/broker","metric":"sqlQuery/time","dataSource":"table_1","value":115},{"service":"druid/broker","metric":"sqlQuery/bytes","dataSource":"table_1","value":1024},{"service":"druid/broker","metric":"sqlQuery/time","dataSource":"table_2","value":12},{"service":"druid/broker","metric":"sqlQuery/bytes","dataSource":"table_2","value":97},{"service":"druid/broker","metric":"sqlQuery/time","dataSource":"table_2","value":18},{"service":"druid/broker","metric":"sqlQuery/bytes","dataSource":"table_2","value":112}]`)),
				Header: map[string][]string{
					textproto.CanonicalMIMEHeaderKey("Authorization"): {ENCODED_CREDENTIALS},
					textproto.CanonicalMIMEHeaderKey("Content-Type"):  {"application/json"},
				},
			},
			expectedStatusCode: http.StatusOK,
			metricExpected:     true,
			consumerFailure:    false,
			configBasicAuth: &BasicAuth{
				Username: VALID_USERNAME,
				Password: VALID_PASSWORD,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var consumer consumer.Metrics
			if tc.consumerFailure {
				consumer = consumertest.NewErr(errors.New("consumer failure"))
			} else {
				consumer = &consumertest.MetricsSink{}
			}

			r := newReceiver(t, &Config{
				Metrics: MetricsConfig{
					Endpoint:             "0.0.0.0:12345",
					BasicAuth:            tc.configBasicAuth,
					MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
				},
			}, consumer)

			rec := httptest.NewRecorder()
			r.handleRequest(rec, tc.request)

			assert.Equal(t, tc.expectedStatusCode, rec.Code, "Status codes are not equal")

			if !tc.consumerFailure {
				if tc.metricExpected {
					assert.Equal(t, 11, consumer.(*consumertest.MetricsSink).DataPointCount(), "Did not receive metrics")
				} else {
					assert.Equal(t, 0, consumer.(*consumertest.MetricsSink).DataPointCount(), "Received metrics when they should have been dropped")
				}
			}
		})
	}
}

func newReceiver(t *testing.T, cfg *Config, nextConsumer consumer.Metrics) *metricsReceiver {
	set := receivertest.NewNopCreateSettings()
	set.Logger = zaptest.NewLogger(t)
	r, err := newMetricsReceiver(set, cfg, nextConsumer)
	require.NoError(t, err)
	return r
}
