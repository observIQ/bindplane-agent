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

package httpreceiver

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

type httpLogsReceiver struct {
	address  string
	path     string
	server   *http.Server
	tls      *configtls.TLSServerSetting
	consumer consumer.Logs
	wg       *sync.WaitGroup
	logger   *zap.Logger
}

// newHTTPLogsReceiver returns a newly configured httpLogsReceiver
func newHTTPLogsReceiver(params receiver.CreateSettings, cfg *Config, consumer consumer.Logs) (*httpLogsReceiver, error) {
	var TLSConfig *tls.Config
	var err error
	if cfg.TLS != nil {
		TLSConfig, err = cfg.TLS.LoadTLSConfig()
		if err != nil {
			return nil, err
		}
	}

	r := &httpLogsReceiver{
		address:  cfg.Endpoint,
		path:     cfg.Path,
		tls:      cfg.TLS,
		consumer: consumer,
		wg:       &sync.WaitGroup{},
		logger:   params.Logger,
	}
	s := &http.Server{
		TLSConfig:         TLSConfig,
		Handler:           http.HandlerFunc(r.handleRequest),
		ReadHeaderTimeout: 20 * time.Second,
	}
	r.server = s
	return r, nil
}

// Start calls startListening
func (r *httpLogsReceiver) Start(ctx context.Context, host component.Host) error {
	return r.startListening(ctx, host)
}

// startListening starts serve on the server using TLS depending on receiver configuration
func (r *httpLogsReceiver) startListening(ctx context.Context, host component.Host) error {
	r.logger.Debug("starting receiver HTTP server")

	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", r.address)
	if err != nil {
		return err
	}

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()

		if r.tls != nil {
			r.logger.Debug("starting ServeTLS",
				zap.String("address", r.address),
				zap.String("cert_file", r.tls.CertFile),
				zap.String("key_file", r.tls.KeyFile),
			)

			err := r.server.ServeTLS(listener, r.tls.CertFile, r.tls.KeyFile)
			r.logger.Debug("ServeTLS done")
			if err != http.ErrServerClosed {
				r.logger.Error("ServeTLS failed", zap.Error(err))
				host.ReportFatalError(err)
			}
		} else {
			r.logger.Debug("starting to serve",
				zap.String("address", r.address),
			)

			err := r.server.Serve(listener)
			r.logger.Debug("Serve done")
			if err != http.ErrServerClosed {
				r.logger.Error("Serve failed", zap.Error(err))
				host.ReportFatalError(err)
			}
		}
	}()

	return nil
}

// Shutdown calls shutdownListener
func (r *httpLogsReceiver) Shutdown(ctx context.Context) error {
	return r.shutdownListener(ctx)
}

// shutdownLIstener tells the server to stop serving and waits for it to stop
func (r *httpLogsReceiver) shutdownListener(ctx context.Context) error {
	r.logger.Debug("shutting down server")

	if err := r.server.Shutdown(ctx); err != nil {
		return err
	}
	r.logger.Debug("waiting for shutdown to complete")
	r.wg.Wait()
	return nil

}

// handleRequest is the function the server uses for requests; calls ConsumeLogs
func (r *httpLogsReceiver) handleRequest(rw http.ResponseWriter, req *http.Request) {
	// path was configured && this req.URL does not match it
	if r.path != "" && req.URL.Path != r.path {
		rw.WriteHeader(http.StatusNotFound)
		r.logger.Debug("received request to path that does not match the configured path", zap.String("request path", req.URL.Path))
		return
	}

	// read in request body
	var payload []byte
	if req.Header.Get("Content-Encoding") == "gzip" {
		r.logger.Debug("req header has Content-Encoding set to gzip")
		reader, err := gzip.NewReader(req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusUnprocessableEntity)
			r.logger.Debug("got payload with gzip compression but failed to read", zap.Error(err))
			return
		}
		defer reader.Close()

		// read the decompressed response body
		payload, err = io.ReadAll(reader)
		if err != nil {
			rw.WriteHeader(http.StatusUnprocessableEntity)
			r.logger.Debug("got payload with gzip compression but failed to read uncompressed payload", zap.Error(err))
			return
		}
	} else {
		r.logger.Debug("req header does not have Content-Encoding set to gzip")
		var err error
		payload, err = io.ReadAll(req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusUnprocessableEntity)
			r.logger.Debug("failed to read logs payload", zap.Error(err), zap.String("remote", req.RemoteAddr))
			return
		}
	}

	// parse []byte into map structure
	logs, err := parsePayload(payload)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		r.logger.Error("failed to convert log request payload to maps", zap.Error(err))
		return
	}

	// consume logs after processing
	if err := r.consumer.ConsumeLogs(req.Context(), r.processLogs(pcommon.NewTimestampFromTime(time.Now()), logs)); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		r.logger.Error("failed to consume logs", zap.Error(err))
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// processLogs transforms the parsed payload into plog.Logs
func (r *httpLogsReceiver) processLogs(now pcommon.Timestamp, logs []map[string]interface{}) plog.Logs {
	pLogs := plog.NewLogs()
	resourceLogs := pLogs.ResourceLogs().AppendEmpty()
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()

	for _, log := range logs {
		logRecord := scopeLogs.LogRecords().AppendEmpty()

		logRecord.SetObservedTimestamp(now)

		if err := logRecord.Body().SetEmptyMap().FromRaw(log); err != nil {
			r.logger.Warn("unable to set log body", zap.Error(err))
		}
	}

	return pLogs
}

// parsePayload transforms the payload into []map[string]interface structure
func parsePayload(payload []byte) ([]map[string]interface{}, error) {
	rawLogsArray := []json.RawMessage{}
	rawLogObject := json.RawMessage{}
	if err := json.Unmarshal(payload, &rawLogsArray); err != nil {
		if err.Error() == "json: cannot unmarshal object into Go value of type []json.RawMessage" {
			if err = json.Unmarshal(payload, &rawLogObject); err != nil {
				return nil, err
			}
			return parseJSONObject(rawLogObject)
		}
		return nil, err
	}
	return parseJSONArray(rawLogsArray)
}

func parseJSONObject(rawLog json.RawMessage) ([]map[string]interface{}, error) {
	var logs []map[string]interface{}
	var log map[string]interface{}
	if err := json.Unmarshal(rawLog, &log); err != nil {
		return nil, err
	}
	return append(logs, log), nil
}

func parseJSONArray(rawLogs []json.RawMessage) ([]map[string]interface{}, error) {
	logs := make([]map[string]interface{}, 0, len(rawLogs))
	for _, l := range rawLogs {
		if len(l) == 0 {
			continue
		}
		var log map[string]interface{}
		if err := json.Unmarshal(l, &log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
