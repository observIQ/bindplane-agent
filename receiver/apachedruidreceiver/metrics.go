// Copyright The OpenTelemetry Authors
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

package apachedruidreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachedruidreceiver"

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/observiq/observiq-otel-collector/receiver/apachedruidreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

const (
	DRUID_QUERY_COUNT             string = "query/count"
	DRUID_QUERY_SUCCESS_COUNT     string = "query/success/count"
	DRUID_QUERY_FAILED_COUNT      string = "query/failed/count"
	DRUID_QUERY_INTERRUPTED_COUNT string = "query/interrupted/count"
	DRUID_QUERY_TIMEOUT_COUNT     string = "query/timeout/count"
	DRUID_SQLQUERY_TIME           string = "sqlQuery/time"
	DRUID_SQLQUERY_BYTES          string = "sqlQuery/bytes"
)

type metricsReceiver struct {
	logger         *zap.Logger
	config         *MetricsConfig
	server         *http.Server
	consumer       consumer.Metrics
	wg             *sync.WaitGroup
	id             component.ID // ID of the receiver component
	metricsBuilder *metadata.MetricsBuilder
}

func newMetricsReceiver(params receiver.CreateSettings, cfg *Config, consumer consumer.Metrics) (*metricsReceiver, error) {
	var tlsConfig *tls.Config
	recv := &metricsReceiver{
		config:         &cfg.Metrics,
		consumer:       consumer,
		logger:         params.Logger,
		wg:             &sync.WaitGroup{},
		id:             params.ID,
		metricsBuilder: metadata.NewMetricsBuilder(cfg.Metrics.MetricsBuilderConfig, params),
	}

	if recv.config.TLS != nil {
		var err error

		tlsConfig, err = recv.config.TLS.LoadTLSConfig()
		if err != nil {
			return nil, err
		}
	}

	s := &http.Server{
		TLSConfig:         tlsConfig,
		Handler:           http.HandlerFunc(recv.handleRequest),
		ReadHeaderTimeout: 20 * time.Second,
	}

	recv.server = s
	return recv, nil
}

func (mReceiver *metricsReceiver) Start(ctx context.Context, host component.Host) error {
	return mReceiver.startListening(ctx, host)
}

func (mReceiver *metricsReceiver) Shutdown(ctx context.Context) error {
	mReceiver.logger.Debug("Shutting down server")
	err := mReceiver.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	mReceiver.logger.Debug("Waiting for shutdown to complete.")
	mReceiver.wg.Wait()
	return nil
}

func (mReceiver *metricsReceiver) startListening(ctx context.Context, host component.Host) error {
	mReceiver.logger.Debug("starting receiver HTTP server")
	// We use l.server.Serve* over l.server.ListenAndServe*
	// So that we can catch and return errors relating to binding to network interface on start.
	var listenConfig net.ListenConfig

	listener, err := listenConfig.Listen(ctx, "tcp", mReceiver.config.Endpoint)
	if err != nil {
		return err
	}

	mReceiver.wg.Add(1)
	if mReceiver.config.TLS != nil {
		go func() {
			defer mReceiver.wg.Done()

			mReceiver.logger.Debug("Starting ServeTLS",
				zap.String("address", mReceiver.config.Endpoint),
				zap.String("certfile", mReceiver.config.TLS.CertFile),
				zap.String("keyfile", mReceiver.config.TLS.KeyFile))

			err := mReceiver.server.ServeTLS(listener, mReceiver.config.TLS.CertFile, mReceiver.config.TLS.KeyFile)

			mReceiver.logger.Debug("Serve TLS done")

			if err != http.ErrServerClosed {
				mReceiver.logger.Error("ServeTLS failed", zap.Error(err))
				host.ReportFatalError(err)
			}
		}()
	} else {
		go func() {
			defer mReceiver.wg.Done()

			mReceiver.logger.Debug("Starting Serve",
				zap.String("address", mReceiver.config.Endpoint))

			err = mReceiver.server.Serve(listener)

			mReceiver.logger.Debug("Serve done")

			if err != http.ErrServerClosed {
				mReceiver.logger.Error("Serve failed", zap.Error(err))
				host.ReportFatalError(err)
			}
		}()
	}

	return nil
}

func (mReceiver *metricsReceiver) handleRequest(rw http.ResponseWriter, request *http.Request) {
	if mReceiver.config.BasicAuth != nil {
		auth := request.Header.Get("Authorization")
		credentials := fmt.Sprintf("%s:%s", mReceiver.config.BasicAuth.Username, mReceiver.config.BasicAuth.Password)
		configAuth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials)))
		if auth == "" {
			rw.WriteHeader(http.StatusUnauthorized)
			mReceiver.logger.Debug("Got request with no basic auth credentials when they were specified in config, dropping...")
			return
		} else if auth != configAuth {
			rw.WriteHeader(http.StatusUnauthorized)
			mReceiver.logger.Debug("Got request with incorrect basic auth credentials when they were specified in config, dropping...")
			return
		}
	}

	if request.Method != "POST" {
		rw.WriteHeader(http.StatusBadRequest)
		mReceiver.logger.Debug("Receiver server only accepts POST requests", zap.String("remote", request.RemoteAddr))
		return
	}

	if request.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		mReceiver.logger.Debug("Content type must be JSON", zap.String("remote", request.RemoteAddr))
		return
	}

	payload, err := io.ReadAll(request.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		mReceiver.logger.Debug("Failed to read metrics payload", zap.Error(err), zap.String("remote", request.RemoteAddr))
		return
	}

	var metrics []interface{}
	if err = json.Unmarshal(payload, &metrics); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		mReceiver.logger.Debug("Failed to convert metrics payload from JSON array to golang slice", zap.Error(err), zap.String("remote", request.RemoteAddr))
		return
	}

	if err := mReceiver.consumer.ConsumeMetrics(request.Context(), mReceiver.processMetrics(metrics)); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		mReceiver.logger.Error("Failed to consume payload as metric", zap.Error(err))
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (mReceiver *metricsReceiver) processMetrics(metrics []interface{}) pmetric.Metrics {
	now := pcommon.NewTimestampFromTime(time.Now())

	purgedMetrics, service := purgeMetrics(metrics)
	queryCountStats := getQueryCountStats(purgedMetrics)

	mReceiver.metricsBuilder.RecordApachedruidQueryCountDataPoint(now, int64(queryCountStats[DRUID_QUERY_COUNT]))
	mReceiver.metricsBuilder.RecordApachedruidSuccessQueryCountDataPoint(now, int64(queryCountStats[DRUID_QUERY_SUCCESS_COUNT]))
	mReceiver.metricsBuilder.RecordApachedruidFailedQueryCountDataPoint(now, int64(queryCountStats[DRUID_QUERY_FAILED_COUNT]))
	mReceiver.metricsBuilder.RecordApachedruidInterruptedQueryCountDataPoint(now, int64(queryCountStats[DRUID_QUERY_INTERRUPTED_COUNT]))
	mReceiver.metricsBuilder.RecordApachedruidTimeoutQueryCountDataPoint(now, int64(queryCountStats[DRUID_QUERY_TIMEOUT_COUNT]))
	mReceiver.recordSqlQueryDataPoints(now, purgedMetrics)

	return mReceiver.metricsBuilder.Emit(metadata.WithApachedruidService(service))
}

// remove all metrics published to the receiver besides the ones it cares about
func purgeMetrics(metrics []interface{}) ([]interface{}, string) {
	collectedMetrics := [7]string{DRUID_QUERY_COUNT, DRUID_QUERY_SUCCESS_COUNT, DRUID_QUERY_FAILED_COUNT, DRUID_QUERY_INTERRUPTED_COUNT, DRUID_QUERY_TIMEOUT_COUNT, DRUID_SQLQUERY_TIME, DRUID_SQLQUERY_BYTES}
	purgedArray := make([]interface{}, 0)
	var service string

	for _, dataPoint := range metrics {
		if dataPoint == nil {
			continue
		}

		currentPoint := dataPoint.(map[string]interface{})
		if metric, metricOk := currentPoint["metric"]; !metricOk || !contains(collectedMetrics[:], metric.(string)) {
			continue
		}

		purgedArray = append(purgedArray, currentPoint)
		if service == "" {
			if serviceInterface, ok := currentPoint["service"]; ok {
				service = serviceInterface.(string)
			}
		}
	}

	return purgedArray, service
}

func getQueryCountStats(metrics []interface{}) map[string]float64 {
	queryCountStats := make(map[string]float64)

	for _, dataPoint := range metrics {
		currentPoint := dataPoint.(map[string]interface{})
		metric, metricOk := currentPoint["metric"]
		value, valueOk := currentPoint["value"]
		if metricOk && valueOk {
			if v, vOk := value.(float64); vOk {
				queryCountStats[metric.(string)] = v
			}
		}
	}

	return queryCountStats
}

func (mReceiver *metricsReceiver) recordSqlQueryDataPoints(now pcommon.Timestamp, metrics []interface{}) {
	dataSources := make([]string, 0)
	sqlQueryCount := make(map[string]float64)
	sqlQueryTime := make(map[string]float64)
	sqlQueryBytes := make(map[string]float64)
	for _, dataPoint := range metrics {
		currentPoint := dataPoint.(map[string]interface{})
		metric, metricOk := currentPoint["metric"]
		value, valueOk := currentPoint["value"]
		dataSource, dsOk := currentPoint["dataSource"]
		if metricOk && valueOk && dsOk {
			if ds, dOk := dataSource.(string); dOk {
				if m := metric.(string); m == DRUID_SQLQUERY_TIME || m == DRUID_SQLQUERY_BYTES {
					if m == DRUID_SQLQUERY_TIME {
						sqlQueryTime[ds] += value.(float64)
						sqlQueryCount[ds]++
					} else {
						sqlQueryBytes[ds] += value.(float64)
					}
					if !contains(dataSources, ds) {
						dataSources = append(dataSources, ds)
					}
				}
			}
		}
	}

	for _, source := range dataSources {
		count := sqlQueryCount[source]
		mReceiver.metricsBuilder.RecordApachedruidAverageSQLQueryTimeDataPoint(now, sqlQueryTime[source]/count, source)
		mReceiver.metricsBuilder.RecordApachedruidAverageSQLQueryBytesDataPoint(now, sqlQueryBytes[source]/count, source)
		mReceiver.metricsBuilder.RecordApachedruidSQLQueryCountDataPoint(now, int64(count), source)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
