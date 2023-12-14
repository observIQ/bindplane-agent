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
	writer    io.WriteCloser
	marshaler logMarshaler
	endpoint  string
}

func newExporter(cfg *Config, params exporter.CreateSettings) (*chronicleForwarderExporter, error) {
	return &chronicleForwarderExporter{
		cfg:       cfg,
		logger:    params.Logger,
		marshaler: newMarshaler(*cfg, params.TelemetrySettings),
	}, nil
}

func (ce *chronicleForwarderExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}
func (ce *chronicleForwarderExporter) logsDataPusher(ctx context.Context, ld plog.Logs) error {
	// Open connection or file before sending each payload
	if err := ce.openWriter(); err != nil {
		return fmt.Errorf("open writer: %w", err)
	}
	defer ce.writer.Close()

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

func (ce *chronicleForwarderExporter) openWriter() error {
	switch ce.cfg.ExportType {
	case ExportTypeSyslog:
		var conn net.Conn
		var err error
		if ce.cfg.Syslog.TLSSetting != nil {
			tlsConfig, err := ce.cfg.Syslog.TLSSetting.LoadTLSConfig()
			if err != nil {
				return fmt.Errorf("load TLS config: %w", err)
			}
			conn, err = tls.Dial(ce.cfg.Syslog.NetAddr.Transport, ce.cfg.Syslog.NetAddr.Endpoint, tlsConfig)
		} else {
			conn, err = net.Dial(ce.cfg.Syslog.NetAddr.Transport, ce.cfg.Syslog.NetAddr.Endpoint)
		}

		if err != nil {
			return fmt.Errorf("dial: %w", err)
		}
		ce.writer = conn

	case ExportTypeFile:
		var err error
		ce.writer, err = os.OpenFile(ce.cfg.File.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
	default:
		return fmt.Errorf("unsupported export type")
	}
	return nil
}

func (ce *chronicleForwarderExporter) send(msg string) error {
	if !strings.HasSuffix(msg, "\n") {
		msg = fmt.Sprintf("%s%s", msg, "\n")
	}

	_, err := io.WriteString(ce.writer, msg)
	return err
}

// Shutdown stops the exporter and is invoked during shutdown.
func (ce *chronicleForwarderExporter) Shutdown(_ context.Context) error {
	if ce.writer != nil {
		if err := ce.writer.Close(); err != nil {
			return fmt.Errorf("close writer: %w", err)
		}
	}
	return nil
}
