module github.com/observiq/observiq-collector

go 1.16

require (
	github.com/aws/aws-sdk-go v1.40.8 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/klauspost/compress v1.13.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sprocessor v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.33.0
	github.com/open-telemetry/opentelemetry-log-collection v0.20.0
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.33.0
	go.opentelemetry.io/collector/model v0.33.0
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.19.0
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.4 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.0.0-00010101000000-000000000000 => github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.33.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.24.1-0.20210408210148-736647af91e1 => github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.33.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter v0.33.0 => github.com/observiq/opentelemetry-collector-contrib/exporter/observiqexporter v0.0.0-20210813200348-1aefe9072336
