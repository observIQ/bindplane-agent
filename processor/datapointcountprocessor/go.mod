module github.com/observiq/observiq-otel-collector/processor/datapointcountprocessor

go 1.19

require (
	github.com/observiq/observiq-otel-collector/counter v1.27.0
	github.com/observiq/observiq-otel-collector/expr v1.27.0
	github.com/observiq/observiq-otel-collector/receiver/routereceiver v1.27.0
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/collector v0.79.0
	go.opentelemetry.io/collector/component v0.79.0
	go.opentelemetry.io/collector/consumer v0.79.0
	go.opentelemetry.io/collector/pdata v1.0.0-rcv0012
	go.opentelemetry.io/collector/receiver v0.79.0
	go.uber.org/zap v1.24.0
)

require (
	github.com/antonmedv/expr v1.12.5 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/collector/confmap v0.79.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.0.0-rcv0012 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/observiq/observiq-otel-collector/receiver/routereceiver => ../../receiver/routereceiver

replace github.com/observiq/observiq-otel-collector/expr => ../../expr

replace github.com/observiq/observiq-otel-collector/counter => ../../counter
