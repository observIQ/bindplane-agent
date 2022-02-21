package factories

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
)

var defaultExtensions = []component.ExtensionFactory{
	bearertokenauthextension.NewFactory(),
	healthcheckextension.NewFactory(),
	oidcauthextension.NewFactory(),
	pprofextension.NewFactory(),
	zpagesextension.NewFactory(),
	filestorage.NewFactory(),
	ballastextension.NewFactory(),
	componenttest.NewNopExtensionFactory(),
}
