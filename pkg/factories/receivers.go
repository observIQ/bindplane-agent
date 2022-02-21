package factories

import (
	"github.com/observiq/observiq-collector/pkg/receiver/logsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

var defaultReceivers = []component.ReceiverFactory{
	componenttest.NewNopReceiverFactory(),
	filelogreceiver.NewFactory(),
	hostmetricsreceiver.NewFactory(),
	jmxreceiver.NewFactory(),
	k8sclusterreceiver.NewFactory(),
	kubeletstatsreceiver.NewFactory(),
	logsreceiver.NewFactory(),
	mongodbreceiver.NewFactory(),
	mysqlreceiver.NewFactory(),
	otlpreceiver.NewFactory(),
	postgresqlreceiver.NewFactory(),
	prometheusreceiver.NewFactory(),
	statsdreceiver.NewFactory(),
	syslogreceiver.NewFactory(),
	tcplogreceiver.NewFactory(),
	udplogreceiver.NewFactory(),
}
