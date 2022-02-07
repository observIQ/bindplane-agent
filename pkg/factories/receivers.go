package factories

import (
	"github.com/observiq/observiq-collector/pkg/receiver/logsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

var defaultReceivers = []component.ReceiverFactory{
	logsreceiver.NewFactory(),
	otlpreceiver.NewFactory(),
	filelogreceiver.NewFactory(),
	jmxreceiver.NewFactory(),
	syslogreceiver.NewFactory(),
	tcplogreceiver.NewFactory(),
	udplogreceiver.NewFactory(),
	componenttest.NewNopReceiverFactory(),
	kubeletstatsreceiver.NewFactory(),
	k8sclusterreceiver.NewFactory(),
	mysqlreceiver.NewFactory(),
	statsdreceiver.NewFactory(),
	postgresqlreceiver.NewFactory(),
	hostmetricsreceiver.NewFactory(),
}
