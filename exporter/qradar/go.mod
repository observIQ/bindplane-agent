module github.com/observiq/bindplane-agent/exporter/qradar

go 1.22.6

require (
	github.com/observiq/bindplane-agent/expr v1.61.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.110.0
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/collector/component v0.110.0
	go.opentelemetry.io/collector/config/configretry v1.16.0
	go.opentelemetry.io/collector/config/configtls v1.16.0
	go.opentelemetry.io/collector/consumer v0.110.0
	go.opentelemetry.io/collector/exporter v0.110.0
	go.opentelemetry.io/collector/pdata v1.16.0
	go.uber.org/goleak v1.3.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/elastic/go-grok v0.3.1 // indirect
	github.com/elastic/lunes v0.1.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.1.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/ua-parser/uap-go v0.0.0-20240611065828-3a4781585db6 // indirect
	go.opentelemetry.io/collector/component/componentprofiles v0.110.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.16.0 // indirect
	go.opentelemetry.io/collector/consumer/consumerprofiles v0.110.0 // indirect
	go.opentelemetry.io/collector/consumer/consumertest v0.110.0 // indirect
	go.opentelemetry.io/collector/exporter/exporterprofiles v0.110.0 // indirect
	go.opentelemetry.io/collector/extension/experimental/storage v0.110.0 // indirect
	go.opentelemetry.io/collector/internal/globalsignal v0.110.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.110.0 // indirect
	go.opentelemetry.io/collector/pipeline v0.110.0 // indirect
	go.opentelemetry.io/collector/receiver v0.110.0 // indirect
	go.opentelemetry.io/collector/receiver/receiverprofiles v0.110.0 // indirect
	go.opentelemetry.io/collector/semconv v0.110.0 // indirect
	go.opentelemetry.io/otel/sdk v1.30.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/alecthomas/participle/v2 v2.1.1 // indirect
	github.com/antonmedv/expr v1.15.5 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.110.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/collector/config/confignet v1.16.0
	go.opentelemetry.io/collector/config/configtelemetry v0.110.0 // indirect
	go.opentelemetry.io/collector/confmap v1.16.0
	go.opentelemetry.io/collector/extension v0.110.0 // indirect
	go.opentelemetry.io/otel v1.30.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240822170219-fc7c04adadcd // indirect
	google.golang.org/grpc v1.66.2 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/observiq/bindplane-agent/expr => ../../expr
