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

package m365receiver // import "github.com/observiq/observiq-otel-collector/receiver/m365receiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

type lClient interface {
	GetJSON(endpoint string) ([]jsonLogs, error)
	GetToken() error
	shutdown() error
}

type m365LogsReceiver struct {
	settings component.TelemetrySettings
	logger   *zap.Logger
	consumer consumer.Logs
	cfg      *Config
	client   lClient
}

func newM365Logs(cfg *Config, settings receiver.CreateSettings, consumer consumer.Logs) *m365LogsReceiver {
	return &m365LogsReceiver{
		settings: settings.TelemetrySettings,
		logger:   settings.Logger,
		consumer: consumer,
		cfg:      cfg,
	}
}

func (l *m365LogsReceiver) Start(_ context.Context, host component.Host) error {
	httpClient, err := l.cfg.ToClient(host, l.settings)
	if err != nil {
		l.logger.Error("error creating HTTP client", zap.Error(err))
		return err
	}

	l.client = newM365Client(httpClient, l.cfg, "https://manage.office.com/.default")
	err = l.client.GetToken()
	if err != nil {
		l.logger.Error("error creating authorization token", zap.Error(err))
		return err
	}
	return nil
}

func (l *m365LogsReceiver) Shutdown(_ context.Context) error {
	return nil
}
