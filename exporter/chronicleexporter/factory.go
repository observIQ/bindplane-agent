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

package chronicleexporter

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/observiq/bindplane-otel-collector/exporter/chronicleexporter/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
)

// NewFactory creates a new Chronicle exporter factory.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability))
}

// createDefaultConfig creates the default configuration for the exporter.
func createDefaultConfig() component.Config {
	return &Config{
		Protocol:            protocolGRPC,
		TimeoutConfig:       exporterhelper.NewDefaultTimeoutConfig(),
		QueueConfig:         exporterhelper.NewDefaultQueueConfig(),
		BackOffConfig:       configretry.NewDefaultBackOffConfig(),
		OverrideLogType:     true,
		Endpoint:            baseEndpoint,
		Compression:         noCompression,
		CollectAgentMetrics: true,
	}
}

// createLogsExporter creates a new log exporter based on this config.
func createLogsExporter(
	ctx context.Context,
	params exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	chronicleCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("invalid config type")
	}

	var cID string
	sid, ok := params.Resource.Attributes().Get(semconv.AttributeServiceInstanceID)
	if ok {
		cID = sid.AsString()
	} else {
		cID = uuid.New().String()
	}

	exp, err := newExporter(chronicleCfg, params, cID, params.ID.String())
	if err != nil {
		return nil, err
	}

	pusher := exp.logsDataPusher
	if chronicleCfg.Protocol == protocolHTTPS {
		pusher = exp.logsHTTPDataPusher
	}
	return exporterhelper.NewLogs(
		ctx,
		params,
		chronicleCfg,
		pusher,
		exporterhelper.WithCapabilities(exp.Capabilities()),
		exporterhelper.WithTimeout(chronicleCfg.TimeoutConfig),
		exporterhelper.WithQueue(chronicleCfg.QueueConfig),
		exporterhelper.WithRetry(chronicleCfg.BackOffConfig),
		exporterhelper.WithShutdown(exp.Shutdown),
	)
}
