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
	"strings"
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
	// fmt.Print("Request header: ")
	// fmt.Println(request.Header)
	// fmt.Print("Request body: ")
	// fmt.Printf("%s %s %s\n", request.Method, request.Host, request.URL) TODO

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

	mReceiver.logger.Debug("Request body: ", zap.String("payload", string(payload)))

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

	fmt.Println("Pre-processed: ")
	fmt.Println(metrics)
	processedMetrics := mReceiver.mapMetrics(metrics)
	fmt.Println("Post-processed: ")
	fmt.Println(processedMetrics)

	mReceiver.recordApachedruidBrokerAverageQueryTimeDataPoint(now, processedMetrics)
	mReceiver.recordApachedruidBrokerFailedQueryCountDataPoint(now, processedMetrics)
	mReceiver.recordApachedruidBrokerQueryCountDataPoint(now, processedMetrics)
	mReceiver.recordApachedruidHistoricalAverageQueryTimeDataPoint(now, processedMetrics)
	mReceiver.recordApachedruidHistoricalFailedQueryCountDataPoint(now, processedMetrics)
	mReceiver.recordApachedruidHistoricalQueryCountDataPoint(now, processedMetrics)

	// not quite sure this makes sense. seems odd to be using the "Druid" node name when the Endpoint is actually the one the server is listening on...
	return mReceiver.metricsBuilder.Emit(metadata.WithApachedruidNodeName(mReceiver.config.Endpoint))
}

func (mReceiver *metricsReceiver) mapMetrics(metrics []interface{}) map[string]float64 {
	metricsMap := make(map[string]float64)
	for _, dataPoint := range metrics {
		if dataPoint == nil {
			continue
		}

		currentPoint := dataPoint.(map[string]interface{})
		service, serviceOk := currentPoint["service"]
		metric, metricOk := currentPoint["metric"]
		valueInterface, valueOk := currentPoint["value"]
		if !(serviceOk && metricOk && valueOk) {
			continue
		}
		value, ok := valueInterface.(float64)
		if !ok {
			mReceiver.logger.Error("Failed to parse metric value as float")
		}

		metricsMap[strings.TrimSpace(service.(string))+"/"+strings.TrimSpace(metric.(string))] += value
	}

	return metricsMap
}

func (mReceiver *metricsReceiver) recordApachedruidBrokerAverageQueryTimeDataPoint(now pcommon.Timestamp, metrics map[string]float64) {
	var averageQueryTime float64
	if totalQueries, ok := metrics["druid/broker/query/count"]; ok && totalQueries > 0 {
		if totalQueryTime, ok := metrics["druid/broker/query/time"]; ok {
			averageQueryTime = totalQueryTime / totalQueries
		}
	}

	mReceiver.metricsBuilder.RecordApachedruidBrokerAverageQueryTimeDataPoint(now, averageQueryTime)
}

func (mReceiver *metricsReceiver) recordApachedruidBrokerFailedQueryCountDataPoint(now pcommon.Timestamp, metrics map[string]float64) {
	mReceiver.metricsBuilder.RecordApachedruidBrokerFailedQueryCountDataPoint(now, int64(metrics["druid/broker/query/failed/count"]))
}

func (mReceiver *metricsReceiver) recordApachedruidBrokerQueryCountDataPoint(now pcommon.Timestamp, metrics map[string]float64) {
	mReceiver.metricsBuilder.RecordApachedruidBrokerQueryCountDataPoint(now, int64(metrics["druid/broker/query/count"]))
}

func (mReceiver *metricsReceiver) recordApachedruidHistoricalAverageQueryTimeDataPoint(now pcommon.Timestamp, metrics map[string]float64) {
	var averageQueryTime float64
	if totalQueries, ok := metrics["druid/historical/query/count"]; ok && totalQueries > 0 {
		if totalQueryTime, ok := metrics["druid/historical/query/time"]; ok {
			averageQueryTime = totalQueryTime / totalQueries
		}
	}

	mReceiver.metricsBuilder.RecordApachedruidHistoricalAverageQueryTimeDataPoint(now, averageQueryTime)
}

func (mReceiver *metricsReceiver) recordApachedruidHistoricalFailedQueryCountDataPoint(now pcommon.Timestamp, metrics map[string]float64) {
	mReceiver.metricsBuilder.RecordApachedruidHistoricalFailedQueryCountDataPoint(now, int64(metrics["druid/historical/query/failed/count"]))
}

func (mReceiver *metricsReceiver) recordApachedruidHistoricalQueryCountDataPoint(now pcommon.Timestamp, metrics map[string]float64) {
	mReceiver.metricsBuilder.RecordApachedruidHistoricalQueryCountDataPoint(now, int64(metrics["druid/historical/query/count"]))
}
