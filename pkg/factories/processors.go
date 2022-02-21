package factories

import (
	"github.com/GoogleCloudPlatform/opentelemetry-operations-collector/processor/normalizesumsprocessor"
	"github.com/observiq/observiq-collector/pkg/processor/resourceattributetransposerprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
)

var defaultProcessors = []component.ProcessorFactory{
	attributesprocessor.NewFactory(),
	batchprocessor.NewFactory(),
	componenttest.NewNopProcessorFactory(),
	groupbyattrsprocessor.NewFactory(),
	k8sattributesprocessor.NewFactory(),
	memorylimiterprocessor.NewFactory(),
	metricstransformprocessor.NewFactory(),
	normalizesumsprocessor.NewFactory(),
	probabilisticsamplerprocessor.NewFactory(),
	resourceattributetransposerprocessor.NewFactory(),
	resourcedetectionprocessor.NewFactory(),
	resourceprocessor.NewFactory(),
}
