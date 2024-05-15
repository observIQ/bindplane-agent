module github.com/observiq/bindplane-agent/plugindocgen

go 1.21.9

require (
	github.com/observiq/bindplane-agent/receiver/pluginreceiver v1.52.0
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/apache/arrow/go/v15 v15.0.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/knadh/koanf/v2 v2.1.1 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20190525184631-5f46317e436b // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage v0.100.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_golang v1.19.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.13.0 // indirect
	github.com/shirou/gopsutil/v3 v3.24.3 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/spf13/cobra v1.8.0 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.etcd.io/bbolt v1.3.9 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector v0.100.0 // indirect
	go.opentelemetry.io/collector/component v0.100.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/converter/expandconverter v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/envprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/fileprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/httpprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/httpsprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/yamlprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/connector v0.100.0 // indirect
	go.opentelemetry.io/collector/consumer v0.100.0 // indirect
	go.opentelemetry.io/collector/exporter v0.100.0 // indirect
	go.opentelemetry.io/collector/extension v0.100.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.7.0 // indirect
	go.opentelemetry.io/collector/otelcol v0.100.0 // indirect
	go.opentelemetry.io/collector/pdata v1.7.0 // indirect
	go.opentelemetry.io/collector/processor v0.100.0 // indirect
	go.opentelemetry.io/collector/receiver v0.100.0 // indirect
	go.opentelemetry.io/collector/semconv v0.100.0 // indirect
	go.opentelemetry.io/collector/service v0.100.0 // indirect
	go.opentelemetry.io/contrib/config v0.6.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.26.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/bridge/opencensus v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.48.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	go.opentelemetry.io/proto/otlp v1.2.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20240205201215-2c58cdc269a3 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gonum.org/v1/gonum v0.15.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.34.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/observiq/bindplane-agent/receiver/pluginreceiver => ../../receiver/pluginreceiver

// Fixes ambiguous import with google.cloud.com/go/compute/metadata
replace cloud.google.com/go => cloud.google.com/go v0.105.0
