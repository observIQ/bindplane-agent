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
	"github.com/GoogleCloudPlatform/opentelemetry-operations-collector/processor/normalizesumsprocessor"
	"github.com/observiq/observiq-otel-collector/processor/resourceattributetransposerprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
)

var defaultProcessors = []component.ProcessorFactory{
	attributesprocessor.NewFactory(),
	batchprocessor.NewFactory(),
	componenttest.NewNopProcessorFactory(),
	cumulativetodeltaprocessor.NewFactory(),
	deltatorateprocessor.NewFactory(),
	filterprocessor.NewFactory(),
	groupbyattrsprocessor.NewFactory(),
	groupbytraceprocessor.NewFactory(),
	k8sattributesprocessor.NewFactory(),
	memorylimiterprocessor.NewFactory(),
	metricsgenerationprocessor.NewFactory(),
	metricstransformprocessor.NewFactory(),
	normalizesumsprocessor.NewFactory(),
	probabilisticsamplerprocessor.NewFactory(),
	resourceattributetransposerprocessor.NewFactory(),
	resourcedetectionprocessor.NewFactory(),
	resourceprocessor.NewFactory(),
	routingprocessor.NewFactory(),
	spanmetricsprocessor.NewFactory(),
	spanprocessor.NewFactory(),
	tailsamplingprocessor.NewFactory(),
	transformprocessor.NewFactory(),
}
