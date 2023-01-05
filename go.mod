module github.com/observiq/observiq-otel-collector

go 1.18

require (
	github.com/google/uuid v1.3.0
	github.com/mholt/archiver/v3 v3.5.1
	github.com/observiq/observiq-otel-collector/exporter/googlecloudexporter v1.14.1
	github.com/observiq/observiq-otel-collector/packagestate v1.14.1
	github.com/observiq/observiq-otel-collector/processor/logcountprocessor v1.14.1
	github.com/observiq/observiq-otel-collector/processor/maskprocessor v1.14.1
	github.com/observiq/observiq-otel-collector/processor/metricextractprocessor v1.14.1
	github.com/observiq/observiq-otel-collector/processor/resourceattributetransposerprocessor v1.14.1
	github.com/observiq/observiq-otel-collector/processor/samplingprocessor v1.14.1
	github.com/observiq/observiq-otel-collector/processor/throughputmeasurementprocessor v1.14.1
	github.com/observiq/observiq-otel-collector/receiver/pluginreceiver v1.14.1
	github.com/observiq/observiq-otel-collector/receiver/routereceiver v1.14.1
	github.com/observiq/observiq-otel-collector/receiver/sapnetweaverreceiver v0.14.1
	github.com/oklog/ulid/v2 v2.1.0
	github.com/open-telemetry/opamp-go v0.2.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alibabacloudlogserviceexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awscloudwatchlogsexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsemfexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awskinesisexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuremonitorexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/carbonexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/coralogixexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/dynatraceexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/f5cloudexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudpubsubexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlemanagedprometheusexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/influxdbexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sapmexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/logstransformprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureeventhubreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/bigipreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dotnetdiagnosticsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/iisreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/saphanareceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmpreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.67.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.67.0
	github.com/shirou/gopsutil/v3 v3.22.11
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.1
	go.opencensus.io v0.24.0
	go.opentelemetry.io/collector v0.67.0
	go.opentelemetry.io/collector/component v0.67.0
	go.opentelemetry.io/collector/confmap v0.67.0
	go.opentelemetry.io/collector/consumer v0.67.0
	go.opentelemetry.io/collector/exporter/loggingexporter v0.67.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.67.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.67.0
	go.opentelemetry.io/collector/extension/ballastextension v0.67.0
	go.opentelemetry.io/collector/extension/zpagesextension v0.67.0
	go.opentelemetry.io/collector/pdata v1.0.0-rc1
	go.opentelemetry.io/collector/processor/batchprocessor v0.67.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.67.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.67.0
	go.uber.org/multierr v1.9.0
	go.uber.org/zap v1.24.0
	golang.org/x/sys v0.3.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bmatcuk/doublestar/v4 v4.4.0 // indirect
	github.com/hooklift/gowsdl v0.5.0 // indirect
	github.com/observiq/observiq-otel-collector/internal/expr v1.14.1 // indirect
	github.com/ovh/go-ovh v1.3.0 // indirect
	github.com/relvacode/iso8601 v1.1.0 // indirect
	github.com/signalfx/signalfx-agent/pkg/apm v0.0.0-20220920175102-539ae8d8ba8e // indirect
)

require (
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.18 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.67.0 // indirect; indir7ct
	go.opentelemetry.io/collector/featuregate v0.67.0 // indirect; indir7ct
)

require (
	cloud.google.com/go v0.105.0 // indirect
	cloud.google.com/go/compute v1.14.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.2 // indirect
	cloud.google.com/go/iam v0.8.0 // indirect
	cloud.google.com/go/logging v1.6.1 // indirect
	cloud.google.com/go/longrunning v0.3.0 // indirect
	cloud.google.com/go/monitoring v1.9.0 // indirect
	cloud.google.com/go/pubsub v1.28.0 // indirect
	cloud.google.com/go/spanner v1.41.0 // indirect
	cloud.google.com/go/trace v1.4.0 // indirect
	code.cloudfoundry.org/clock v1.0.0 // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20211115184647-b584dd5df32c // indirect
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible // indirect
	code.cloudfoundry.org/rfc5424 v0.0.0-20201103192249-000122071b78 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.2 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14 // indirect
	github.com/Azure/azure-amqp-common-go/v3 v3.2.3 // indirect
	github.com/Azure/azure-event-hubs-go/v3 v3.3.20 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v65.0.0+incompatible // indirect
	github.com/Azure/azure-storage-blob-go v0.15.0 // indirect
	github.com/Azure/go-amqp v0.17.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.21 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v0.34.3-0.20221202192616-0186b89ba914 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector v0.34.3-0.20221202192616-0186b89ba914 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus v0.34.3-0.20221202192616-0186b89ba914 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.10.2 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.34.3-0.20221202192616-0186b89ba914 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ReneKroon/ttlcache/v2 v2.11.0 // indirect
	github.com/SAP/go-hdb v0.111.5 // indirect
	github.com/Shopify/sarama v1.37.2 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/aerospike/aerospike-client-go/v6 v6.7.0 // indirect
	github.com/alecthomas/participle/v2 v2.0.0-beta.5 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/aliyun/aliyun-log-go-sdk v0.1.41 // indirect
	github.com/andybalholm/brotli v1.0.1 // indirect
	github.com/antonmedv/expr v1.9.0 // indirect
	github.com/apache/arrow/go/arrow v0.0.0-20211112161151-bc219186db40 // indirect
	github.com/apache/thrift v0.17.0 // indirect
	github.com/armon/go-metrics v0.4.0 // indirect
	github.com/aws/aws-sdk-go v1.44.155 // indirect
	github.com/aws/aws-sdk-go-v2 v1.17.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.20 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.33 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.20 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.15.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.27.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.26 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.17.6 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/checkpoint-restore/go-criu/v5 v5.3.0 // indirect
	github.com/cilium/ebpf v0.7.0 // indirect
	github.com/cloudfoundry-incubator/uaago v0.0.0-20190307164349-8136b7bbe76e // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/cncf/xds/go v0.0.0-20220314180256-7f1daf1720fc // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/containerd/ttrpc v1.1.0 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-systemd/v22 v22.4.0 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.12.2 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/digitalocean/godo v1.88.0 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.21+incompatible // indirect
	github.com/docker/go-connections v0.4.1-0.20210727194412-58542c764a11 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/dynatrace-oss/dynatrace-metric-utils-go v0.5.0 // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.1.0 // indirect
	github.com/elastic/go-elasticsearch/v8 v8.5.0 // indirect
	github.com/elastic/go-structform v0.0.10 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/envoyproxy/go-control-plane v0.10.3 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.13 // indirect
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.1 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.22.1 // indirect
	github.com/go-redis/redis/v7 v7.4.1 // indirect
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cadvisor v0.46.0 // indirect
	github.com/google/flatbuffers v2.0.8+incompatible // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gophercloud/gophercloud v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gosnmp/gosnmp v1.35.0 // indirect
	github.com/grafana/regexp v0.0.0-20221005093135-b4c2bcb0a4b6 // indirect
	github.com/grobie/gomemcache v0.0.0-20180201122607-1f779c573665 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	github.com/hashicorp/consul/api v1.18.0 // indirect
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.3.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20221102143410-8a95f1239005 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hetznercloud/hcloud-go v1.35.3 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/influxdata/go-syslog/v3 v3.0.1-0.20210608084020-ac565dc76ba6 // indirect
	github.com/influxdata/influxdb-observability/common v0.2.32 // indirect
	github.com/influxdata/influxdb-observability/influx2otel v0.2.32 // indirect
	github.com/influxdata/influxdb-observability/otel2influx v0.2.32 // indirect
	github.com/influxdata/line-protocol/v2 v2.2.1 // indirect
	github.com/ionos-cloud/sdk-go/v6 v6.1.3 // indirect
	github.com/jaegertracing/jaeger v1.39.1-0.20221110195127-14c11365a856 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.3 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/karrick/godirwalk v1.17.0 // indirect
	github.com/klauspost/compress v1.15.12 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/knadh/koanf v1.4.4 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/leodido/ragel-machinery v0.0.0-20181214104525-299bdde78165 // indirect
	github.com/leoluk/perflib_exporter v0.2.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/linode/linodego v1.9.3 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-ieproxy v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/microsoft/ApplicationInsights-Go v0.4.4 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mindprince/gonvml v0.0.0-20190828220739-9ebdce4bb989 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mongodb-forks/digest v1.0.4 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/mostynb/go-grpc-compression v1.1.17 // indirect
	github.com/mrunalp/fileutils v0.5.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nginxinc/nginx-prometheus-exporter v0.8.1-0.20201110005315-f5a5f8086c19 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/observiq/ctimefmt v1.0.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/cwlogs v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/docker v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/loki v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/signalfx v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.67.0 // indirect; indir7.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/winperfcounters v0.67.0 // indirect; indir7.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/opencontainers/runc v1.1.4 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/opencontainers/selinux v1.10.1 // indirect
	github.com/openlyinc/pointy v1.2.0 // indirect
	github.com/openshift/api v0.0.0-20210521075222-e273a339932a // indirect
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/prometheus/prometheus v0.40.5 // indirect
	github.com/prometheus/statsd_exporter v0.22.7 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.9 // indirect
	github.com/seccomp/libseccomp-golang v0.9.2-0.20220502022130-f33da4d89646 // indirect
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.3 // indirect
	github.com/signalfx/gohistogram v0.0.0-20160107210732-1ccfd2ff5083 // indirect
	github.com/signalfx/golib/v3 v3.3.46 // indirect
	github.com/signalfx/sapm-proto v0.12.0 // indirect
	github.com/sijms/go-ora/v2 v2.5.17 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/snowflakedb/gosnowflake v1.6.15 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v1.6.1 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tidwall/gjson v1.10.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20210330154013-f5de75959ad5 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/vmware/govmomi v0.29.0 // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.mongodb.org/atlas v0.19.0 // indirect
	go.mongodb.org/mongo-driver v1.10.3 // indirect
	go.opentelemetry.io/collector/semconv v0.67.0 // indirect; indir7.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.36.4 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.4 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.11.1 // indirect
	go.opentelemetry.io/contrib/zpages v0.36.4 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.33.0 // indirect
	go.opentelemetry.io/otel/metric v0.33.0 // indirect
	go.opentelemetry.io/otel/sdk v1.11.1 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/goleak v1.2.0 // indirect
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/exp v0.0.0-20221205204356-47842c84f3db // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/net v0.3.0 // indirect
	golang.org/x/oauth2 v0.2.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.org/x/tools v0.2.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gonum.org/v1/gonum v0.12.0 // indirect
	google.golang.org/api v0.104.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221206210731-b1a01be3a5f6 // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.25.4 // indirect
	k8s.io/apimachinery v0.25.4 // indirect
	k8s.io/client-go v0.25.4 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/kubelet v0.25.4 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace github.com/observiq/observiq-otel-collector/processor/resourceattributetransposerprocessor => ./processor/resourceattributetransposerprocessor

replace github.com/observiq/observiq-otel-collector/receiver/pluginreceiver => ./receiver/pluginreceiver

replace github.com/observiq/observiq-otel-collector/receiver/routereceiver => ./receiver/routereceiver

replace github.com/observiq/observiq-otel-collector/receiver/sapnetweaverreceiver => ./receiver/sapnetweaverreceiver

replace github.com/observiq/observiq-otel-collector/exporter/googlecloudexporter => ./exporter/googlecloudexporter

replace github.com/observiq/observiq-otel-collector/packagestate => ./packagestate

replace github.com/observiq/observiq-otel-collector/processor/throughputmeasurementprocessor => ./processor/throughputmeasurementprocessor

replace github.com/observiq/observiq-otel-collector/processor/samplingprocessor => ./processor/samplingprocessor

replace github.com/observiq/observiq-otel-collector/processor/maskprocessor => ./processor/maskprocessor

replace github.com/observiq/observiq-otel-collector/processor/logcountprocessor => ./processor/logcountprocessor

replace github.com/observiq/observiq-otel-collector/processor/metricextractprocessor => ./processor/metricextractprocessor

replace github.com/observiq/observiq-otel-collector/internal/expr => ./internal/expr

// some dependencies attempt to bring something like v1.8.2-0.20220303173753-edfe657b5405, which is older than v0.38.0
// at the time of this inclusion, v0.38.0 was the latest version available (also tagged as v2.38.0)
//replace github.com/prometheus/prometheus => github.com/prometheus/prometheus v0.38.0

// Does not build with windows and only used in configschema executable
// Relevant issue https://github.com/mattn/go-ieproxy/issues/45
replace github.com/mattn/go-ieproxy v0.0.9 => github.com/mattn/go-ieproxy v0.0.1

// This is a fork of the official opentelemetry-operations-go exporter that removes the credential autodiscovery logic that conflicts with our own custom logic
replace github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector v0.34.3-0.20221202192616-0186b89ba914 => github.com/observiq/opentelemetry-operations-go/exporter/collector v0.26.1-0.20221215164112-d92100d9c728
