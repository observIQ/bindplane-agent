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

//go:build solaris || illumos

package factories

import (
	"github.com/observiq/bindplane-agent/exporter/azureblobexporter"
	"github.com/observiq/bindplane-agent/exporter/chronicleexporter"
	"github.com/observiq/bindplane-agent/exporter/chronicleforwarderexporter"
	"github.com/observiq/bindplane-agent/exporter/googlecloudexporter"
	"github.com/observiq/bindplane-agent/exporter/googlemanagedprometheusexporter"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter"
	"github.com/observiq/bindplane-agent/internal/version"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alibabacloudlogserviceexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awscloudwatchlogsexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsemfexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awskinesisexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awss3exporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuremonitorexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/carbonexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/coralogixexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudpubsubexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/influxdbexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sapmexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sumologicexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/syslogexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/nopexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
)

var defaultExporters = []exporter.Factory{
	alibabacloudlogserviceexporter.NewFactory(),
	awscloudwatchlogsexporter.NewFactory(),
	awsemfexporter.NewFactory(),
	awskinesisexporter.NewFactory(),
	awss3exporter.NewFactory(),
	awsxrayexporter.NewFactory(),
	azureblobexporter.NewFactory(),
	azuremonitorexporter.NewFactory(),
	carbonexporter.NewFactory(),
	chronicleexporter.NewFactory(),
	chronicleforwarderexporter.NewFactory(),
	clickhouseexporter.NewFactory(),
	coralogixexporter.NewFactory(),
	// datadogexporter.NewFactory(),
	elasticsearchexporter.NewFactory(),
	fileexporter.NewFactory(),
	googlecloudexporter.NewFactory(version.Version()),
	googlecloudpubsubexporter.NewFactory(),
	googlemanagedprometheusexporter.NewFactory(version.Version()),
	influxdbexporter.NewFactory(),
	kafkaexporter.NewFactory(),
	loadbalancingexporter.NewFactory(),
	loggingexporter.NewFactory(),
	logzioexporter.NewFactory(),
	lokiexporter.NewFactory(),
	nopexporter.NewFactory(),
	opencensusexporter.NewFactory(),
	otlpexporter.NewFactory(),
	otlphttpexporter.NewFactory(),
	prometheusexporter.NewFactory(),
	prometheusremotewriteexporter.NewFactory(),
	sapmexporter.NewFactory(),
	signalfxexporter.NewFactory(),
	snowflakeexporter.NewFactory(),
	splunkhecexporter.NewFactory(),
	sumologicexporter.NewFactory(),
	syslogexporter.NewFactory(),
	zipkinexporter.NewFactory(),
}
