module github.com/observiq/observiq-otel-collector

go 1.17

require (
	github.com/Azure/azure-event-hubs-go/v3 v3.3.17
	github.com/GoogleCloudPlatform/opentelemetry-operations-collector v0.0.3-0.20211123195618-f15f911ae5a1
	github.com/aws/aws-sdk-go v1.43.18
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jpillora/backoff v1.0.0
	github.com/json-iterator/go v1.1.12
	github.com/mitchellh/mapstructure v1.4.3
	github.com/observiq/goflow/v3 v3.4.4
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alibabacloudlogserviceexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awscloudwatchlogsexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsemfexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awskinesisexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsprometheusremotewriteexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuremonitorexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/carbonexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/f5cloudexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.46.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudpubsubexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/influxdbexporter v0.46.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/observiqexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dotnetdiagnosticsreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.46.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusexecreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.45.1
	github.com/open-telemetry/opentelemetry-log-collection v0.24.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.1
	github.com/testcontainers/testcontainers-go v0.12.0
	go.etcd.io/bbolt v1.3.6
	go.opentelemetry.io/collector v0.46.0
	go.opentelemetry.io/collector/model v0.46.0
	go.uber.org/multierr v1.8.0
	go.uber.org/zap v1.21.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
)

require (
	cloud.google.com/go v0.100.2 // indirect
	cloud.google.com/go/logging v1.4.2 // indirect
	cloud.google.com/go/spanner v1.29.0 // indirect
	code.cloudfoundry.org/clock v1.0.0 // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20211115184647-b584dd5df32c // indirect
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible // indirect
	code.cloudfoundry.org/rfc5424 v0.0.0-20201103192249-000122071b78 // indirect
	github.com/ReneKroon/ttlcache/v2 v2.11.0 // indirect
	github.com/Shopify/sarama v1.31.1 // indirect
	github.com/alecthomas/participle/v2 v2.0.0-alpha7 // indirect
	github.com/aliyun/aliyun-log-go-sdk v0.1.27 // indirect
	github.com/apache/thrift v0.15.0 // indirect
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/checkpoint-restore/go-criu/v5 v5.0.0 // indirect
	github.com/cilium/ebpf v0.6.2 // indirect
	github.com/cloudfoundry-incubator/uaago v0.0.0-20190307164349-8136b7bbe76e // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/containerd/console v1.0.2 // indirect
	github.com/containerd/ttrpc v1.1.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/elastic/go-elasticsearch/v7 v7.17.0 // indirect
	github.com/elastic/go-structform v0.0.9 // indirect
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/go-kit/kit v0.11.0 // indirect
	github.com/go-redis/redis/v7 v7.4.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godbus/dbus/v5 v5.0.4 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/google/cadvisor v0.43.0 // indirect
	github.com/grobie/gomemcache v0.0.0-20180201122607-1f779c573665 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/influxdata/influxdb-observability/common v0.2.14 // indirect
	github.com/influxdata/influxdb-observability/influx2otel v0.2.14 // indirect
	github.com/influxdata/influxdb-observability/otel2influx v0.2.14 // indirect
	github.com/influxdata/line-protocol/v2 v2.2.1 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/microsoft/ApplicationInsights-Go v0.4.4 // indirect
	github.com/mindprince/gonvml v0.0.0-20190828220739-9ebdce4bb989 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mrunalp/fileutils v0.5.0 // indirect
	github.com/nginxinc/nginx-prometheus-exporter v0.8.1-0.20201110005315-f5a5f8086c19 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/cwlogs v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/docker v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.45.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.45.1 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/opencontainers/selinux v1.8.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.0 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/seccomp/libseccomp-golang v0.9.1 // indirect
	github.com/signalfx/sapm-proto v0.7.2 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.10.1 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tidwall/gjson v1.10.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/vishvananda/netlink v1.1.1-0.20201029203352-d40f9887b852 // indirect
	github.com/vishvananda/netns v0.0.0-20200728191858-db3c7e526aae // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	k8s.io/klog v1.0.0 // indirect
)

require (
	cloud.google.com/go/compute v1.5.0 // indirect
	cloud.google.com/go/monitoring v1.1.0 // indirect
	cloud.google.com/go/trace v1.0.0 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.10 // indirect
	github.com/Azure/azure-amqp-common-go/v3 v3.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v61.1.0+incompatible // indirect
	github.com/Azure/go-amqp v0.17.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.23 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.18 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.3.0 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/Microsoft/hcsshim v0.8.23 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/antonmedv/expr v1.9.0 // indirect
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v3 v3.0.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cncf/xds/go v0.0.0-20211130200136-a8f946100490 // indirect
	github.com/containerd/cgroups v1.0.1 // indirect
	github.com/containerd/containerd v1.5.10 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/digitalocean/godo v1.73.0 // indirect
	github.com/docker/distribution v2.8.0+incompatible // indirect
	github.com/docker/docker v20.10.12+incompatible // indirect
	github.com/docker/go-connections v0.4.1-0.20210727194412-58542c764a11 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/envoyproxy/go-control-plane v0.10.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.2 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/go-zookeeper/zk v1.0.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gophercloud/gophercloud v0.24.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/consul/api v1.12.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.9.6 // indirect
	github.com/hetznercloud/hcloud-go v1.33.1 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jaegertracing/jaeger v1.31.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.14.4 // indirect
	github.com/knadh/koanf v1.4.0 // indirect
	github.com/kolo/xmlrpc v0.0.0-20201022064351-38db28db192b // indirect
	github.com/leoluk/perflib_exporter v0.1.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/libp2p/go-reuseport v0.0.1 // indirect
	github.com/linode/linodego v1.2.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/miekg/dns v1.1.45 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.5.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mongodb-forks/digest v1.0.3 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mostynb/go-grpc-compression v1.1.16 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/observiq/ctimefmt v1.0.0 // indirect
	github.com/observiq/go-syslog/v3 v3.0.2 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.45.1 // indirect; indir5.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.46.0 // indirect; indir5.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.46.0 // indirect; indir5.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.45.1 // indirect; indir5.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet v0.45.1 // indirect; indir5.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.45.1 // indirect; indir5.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/stanza v0.45.1 // indirect; indir5.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.45.1 // indirect; indir5.1
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.46.0 // indirect; indir5.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.0.3 // indirect
	github.com/openlyinc/pointy v1.1.2 // indirect
	github.com/openshift/api v0.0.0-20210521075222-e273a339932a // indirect
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/prometheus/prometheus v1.8.2-0.20220117154355-4855a0c067e2 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.7.0.20210223165440-c65ae3540d44 // indirect
	github.com/shirou/gopsutil/v3 v3.22.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/cobra v1.3.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.0 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.mongodb.org/atlas v0.15.0 // indirect
	go.mongodb.org/mongo-driver v1.8.3 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.29.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.29.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.29.0 // indirect
	go.opentelemetry.io/otel v1.4.1 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.27.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.27.0 // indirect
	go.opentelemetry.io/otel/metric v0.27.0 // indirect
	go.opentelemetry.io/otel/sdk v1.4.1 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.27.0 // indirect
	go.opentelemetry.io/otel/trace v1.4.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	golang.org/x/crypto v0.0.0-20220128200615-198e4374d7ed // indirect
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	golang.org/x/tools v0.1.9 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gonum.org/v1/gonum v0.9.3 // indirect
	google.golang.org/api v0.70.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220222213610-43724f9ea8cf // indirect
	google.golang.org/grpc v1.44.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.23.4 // indirect
	k8s.io/klog/v2 v2.40.1 // indirect
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65 // indirect
	k8s.io/kubelet v0.23.3 // indirect
	k8s.io/utils v0.0.0-20211116205334-6203023598ed // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.46.0 => github.com/observiq/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.0.0-20220304152956-bb36c08bd895
