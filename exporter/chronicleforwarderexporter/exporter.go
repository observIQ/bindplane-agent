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

package chronicleforwarderexporter

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type chronicleForwarderExporter struct {
	cfg       *Config
	logger    *zap.Logger
	writer    io.Writer
	marshaler logMarshaler
	endpoint  string
}

func newExporter(cfg *Config, params exporter.CreateSettings) (*chronicleForwarderExporter, error) {
	exporter := &chronicleForwarderExporter{
		cfg:       cfg,
		logger:    params.Logger,
		marshaler: newMarshaler(*cfg, params.TelemetrySettings),
	}

	switch cfg.ExportType {
	case ExportTypeSyslog:
		endpoint := buildEndpoint(cfg)

		var conn net.Conn
		var err error
		if cfg.Syslog.TLSSetting != nil {
			tlsConfig, err := cfg.Syslog.TLSSetting.LoadTLSConfig()
			if err != nil {
				return nil, fmt.Errorf("load TLS config: %w", err)
			}
			conn, err = tls.Dial(cfg.Syslog.Network, endpoint, tlsConfig)
		} else {
			conn, err = net.Dial(cfg.Syslog.Network, endpoint)
		}

		if err != nil {
			return nil, fmt.Errorf("dial: %w", err)
		}
		exporter.writer = conn

	case ExportTypeFile:
		var err error
		exporter.writer, err = os.OpenFile(cfg.File.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}
	}

	return exporter, nil
}

func buildEndpoint(cfg *Config) string {
	return fmt.Sprintf("%s:%d", cfg.Syslog.Host, cfg.Syslog.Port)
}

func (ce *chronicleForwarderExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (ce *chronicleForwarderExporter) logsDataPusher(ctx context.Context, ld plog.Logs) error {
	payloads, err := ce.marshaler.MarshalRawLogs(ctx, ld)
	if err != nil {
		return fmt.Errorf("marshal logs: %w", err)
	}

	for _, payload := range payloads {
		if err := ce.send(payload); err != nil {
			return fmt.Errorf("upload to Chronicle forwarder: %w", err)
		}
	}

	return nil
}

func (s *chronicleForwarderExporter) send(msg string) error {
	if !strings.HasSuffix(msg, "\n") {
		msg = fmt.Sprintf("%s%s", msg, "\n")
	}
	_, err := fmt.Fprint(s.writer, msg)
	return err
}

func (s *chronicleForwarderExporter) Shutdown(ctx context.Context) error {
	if s.writer != nil {
		if err := s.writer.(io.Closer).Close(); err != nil {
			return fmt.Errorf("close writer: %w", err)
		}
	}
	return nil
}
