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

package factories

import (
	"github.com/observiq/observiq-otel-collector/receiver/logsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dotnetdiagnosticsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusexecreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

var defaultReceivers = []component.ReceiverFactory{
	apachereceiver.NewFactory(),
	awscontainerinsightreceiver.NewFactory(),
	awsecscontainermetricsreceiver.NewFactory(),
	awsfirehosereceiver.NewFactory(),
	awsxrayreceiver.NewFactory(),
	carbonreceiver.NewFactory(),
	cloudfoundryreceiver.NewFactory(),
	collectdreceiver.NewFactory(),
	componenttest.NewNopReceiverFactory(),
	couchdbreceiver.NewFactory(),
	dockerstatsreceiver.NewFactory(),
	dotnetdiagnosticsreceiver.NewFactory(),
	elasticsearchreceiver.NewFactory(),
	filelogreceiver.NewFactory(),
	fluentforwardreceiver.NewFactory(),
	googlecloudpubsubreceiver.NewFactory(),
	googlecloudspannerreceiver.NewFactory(),
	hostmetricsreceiver.NewFactory(),
	influxdbreceiver.NewFactory(),
	jaegerreceiver.NewFactory(),
	jmxreceiver.NewFactory(),
	journaldreceiver.NewFactory(),
	k8sclusterreceiver.NewFactory(),
	k8seventsreceiver.NewFactory(),
	kafkametricsreceiver.NewFactory(),
	kafkareceiver.NewFactory(),
	kubeletstatsreceiver.NewFactory(),
	logsreceiver.NewFactory(),
	memcachedreceiver.NewFactory(),
	mongodbatlasreceiver.NewFactory(),
	mongodbreceiver.NewFactory(),
	mysqlreceiver.NewFactory(),
	nginxreceiver.NewFactory(),
	opencensusreceiver.NewFactory(),
	otlpreceiver.NewFactory(),
	podmanreceiver.NewFactory(),
	postgresqlreceiver.NewFactory(),
	prometheusexecreceiver.NewFactory(),
	prometheusreceiver.NewFactory(),
	rabbitmqreceiver.NewFactory(),
	redisreceiver.NewFactory(),
	sapmreceiver.NewFactory(),
	simpleprometheusreceiver.NewFactory(),
	statsdreceiver.NewFactory(),
	syslogreceiver.NewFactory(),
	tcplogreceiver.NewFactory(),
	udplogreceiver.NewFactory(),
	windowsperfcountersreceiver.NewFactory(),
	zipkinreceiver.NewFactory(),
	zookeeperreceiver.NewFactory(),
}
