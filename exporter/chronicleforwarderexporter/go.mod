module github.com/observiq/bindplane-agent/exporter/chronicleforwarderexporter

go 1.20

require (
	github.com/observiq/bindplane-agent/expr v1.41.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.91.0
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/collector/component v0.91.0
	go.opentelemetry.io/collector/config/configtls v0.91.0
	go.opentelemetry.io/collector/consumer v0.91.0
	go.opentelemetry.io/collector/exporter v0.91.0
	go.opentelemetry.io/collector/pdata v1.0.0
	go.uber.org/zap v1.26.0
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v0.91.0 // indirect
)

require (
	github.com/alecthomas/participle/v2 v2.1.1 // indirect
	github.com/antonmedv/expr v1.15.5 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.0.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20220423185008-bf980b35cac4 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.91.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector v0.91.0 // indirect
	go.opentelemetry.io/collector/config/confignet v0.91.0
	go.opentelemetry.io/collector/config/configtelemetry v0.91.0 // indirect
	go.opentelemetry.io/collector/confmap v0.91.0 // indirect
	go.opentelemetry.io/collector/extension v0.91.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.0.0 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20231127185646-65229373498e // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/observiq/bindplane-agent/expr => ../../expr
