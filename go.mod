module github.com/observiq/bindplane-agent

go 1.20

require (
	github.com/google/uuid v1.3.0
	github.com/mholt/archiver/v3 v3.5.1
	github.com/observiq/bindplane-agent/exporter/googlecloudexporter v1.33.0
	github.com/observiq/bindplane-agent/exporter/googlemanagedprometheusexporter v1.33.0
	github.com/observiq/bindplane-agent/packagestate v1.33.0
	github.com/observiq/bindplane-agent/processor/datapointcountprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/logcountprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/logdeduplicationprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/maskprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/metricextractprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/metricstatsprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/removeemptyvaluesprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/resourceattributetransposerprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/samplingprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/spancountprocessor v1.33.0
	github.com/observiq/bindplane-agent/processor/throughputmeasurementprocessor v1.33.0
	github.com/observiq/bindplane-agent/receiver/m365receiver v1.33.0
	github.com/observiq/bindplane-agent/receiver/pluginreceiver v1.33.0
	github.com/observiq/bindplane-agent/receiver/routereceiver v1.33.0
	github.com/observiq/bindplane-agent/receiver/sapnetweaverreceiver v1.33.0
	github.com/oklog/ulid/v2 v2.1.0
	github.com/open-telemetry/opamp-go v0.2.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/countconnector v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/servicegraphconnector v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alibabacloudlogserviceexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awscloudwatchlogsexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsemfexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awskinesisexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awss3exporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuremonitorexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/carbonexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/coralogixexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/dynatraceexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/f5cloudexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudpubsubexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/influxdbexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sapmexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/basicauthextension v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/logstransformprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachesparkreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureeventhubreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/bigipreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudflarereceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/iisreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/saphanareceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmpreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.83.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.83.0
	github.com/shirou/gopsutil/v3 v3.23.7
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	go.opencensus.io v0.24.0
	go.opentelemetry.io/collector v0.83.0
	go.opentelemetry.io/collector/component v0.83.0
	go.opentelemetry.io/collector/confmap v0.83.0
	go.opentelemetry.io/collector/connector v0.83.0
	go.opentelemetry.io/collector/connector/forwardconnector v0.83.0
	go.opentelemetry.io/collector/consumer v0.83.0
	go.opentelemetry.io/collector/exporter v0.83.0
	go.opentelemetry.io/collector/exporter/loggingexporter v0.83.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.83.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.83.0
	go.opentelemetry.io/collector/extension v0.83.0
	go.opentelemetry.io/collector/extension/ballastextension v0.83.0
	go.opentelemetry.io/collector/extension/zpagesextension v0.83.0
	go.opentelemetry.io/collector/featuregate v1.0.0-rcv0014
	go.opentelemetry.io/collector/pdata v1.0.0-rcv0014
	go.opentelemetry.io/collector/processor v0.83.0
	go.opentelemetry.io/collector/processor/batchprocessor v0.83.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.83.0
	go.opentelemetry.io/collector/receiver v0.83.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.83.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.25.0
	golang.org/x/sys v0.11.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.6.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.1.0 // indirect
	github.com/DataDog/agent-payload/v5 v5.0.89 // indirect
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.48.0-beta.1 // indirect
	github.com/DataDog/datadog-agent/pkg/proto v0.48.0-beta.1 // indirect
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.48.0-beta.1 // indirect
	github.com/DataDog/datadog-agent/pkg/trace v0.48.0-devel // indirect
	github.com/DataDog/datadog-agent/pkg/util/cgroups v0.48.0-beta.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/log v0.48.0-beta.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/pointer v0.48.0-beta.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.48.0-beta.1 // indirect
	github.com/DataDog/datadog-api-client-go/v2 v2.14.0 // indirect
	github.com/DataDog/datadog-go/v5 v5.1.1 // indirect
	github.com/DataDog/go-tuf v1.0.1-0.5.2 // indirect
	github.com/DataDog/gohai v0.0.0-20220718130825-1776f9beb9cc // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata v0.7.0 // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.7.0 // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/logs v0.7.0 // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metrics v0.7.0 // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile v0.7.0 // indirect
	github.com/DataDog/sketches-go v1.4.2 // indirect
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/GehirnInc/crypt v0.0.0-20200316065508-bb7000b8a962 // indirect
	github.com/IBM/sarama v1.40.1 // indirect
	github.com/JohnCGriffin/overflow v0.0.0-20211019200055-46fa312c352c // indirect
	github.com/apache/arrow/go/v12 v12.0.1 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/containerd/cgroups v1.0.4 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elastic/go-elasticsearch/v7 v7.17.10 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/grafana/loki/pkg/push v0.0.0-20230127072203-4e8cc8d71928 // indirect
	github.com/influxdata/influxdb-observability/otel2influx v0.5.5 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.3 // indirect
	github.com/knadh/koanf/v2 v2.0.1 // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/observiq/bindplane-agent/counter v1.33.0 // indirect
	github.com/observiq/bindplane-agent/expr v1.33.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlemanagedprometheusexporter v0.83.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/servicegraphprocessor v0.83.0 // indirect
	github.com/outcaste-io/ristretto v0.2.1 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.7.0 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tg123/go-htpasswd v1.2.1 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opentelemetry.io/collector/config/configauth v0.83.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v0.83.0 // indirect
	go.opentelemetry.io/collector/config/configgrpc v0.83.0 // indirect
	go.opentelemetry.io/collector/config/confighttp v0.83.0 // indirect
	go.opentelemetry.io/collector/config/confignet v0.83.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v0.83.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.83.0 // indirect
	go.opentelemetry.io/collector/config/configtls v0.83.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.83.0 // indirect
	go.opentelemetry.io/collector/extension/auth v0.83.0 // indirect
	go.opentelemetry.io/otel/bridge/opencensus v0.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v0.39.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.39.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.16.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230807174057-1744710a1577 // indirect
	gopkg.in/zorkian/go-datadog-api.v2 v2.30.0 // indirect
	sigs.k8s.io/controller-runtime v0.15.1 // indirect
)

require (
	github.com/Azure/azure-amqp-common-go/v4 v4.2.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.0 // indirect
	github.com/hooklift/gowsdl v0.5.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.83.0 // indirect; indir2ct
	github.com/ovh/go-ovh v1.4.1 // indirect
	github.com/relvacode/iso8601 v1.3.0 // indirect
	github.com/signalfx/signalfx-agent/pkg/apm v0.0.0-20230214151822-6a6813cf5bf1 // indirect
	github.com/tilinna/clock v1.1.0 // indirect
)

require (
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.26 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.83.0 // indirect; indir2ct
)

require (
	cloud.google.com/go v0.110.6 // indirect
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.4-0.20230617002413-005d2dfb6b68 // indirect
	cloud.google.com/go/iam v1.1.1 // indirect
	cloud.google.com/go/logging v1.7.0 // indirect
	cloud.google.com/go/longrunning v0.5.1 // indirect
	cloud.google.com/go/monitoring v1.15.1 // indirect
	cloud.google.com/go/pubsub v1.33.0 // indirect
	cloud.google.com/go/spanner v1.47.0 // indirect
	cloud.google.com/go/trace v1.10.1 // indirect
	code.cloudfoundry.org/clock v1.0.0 // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20211115184647-b584dd5df32c // indirect
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible // indirect
	code.cloudfoundry.org/rfc5424 v0.0.0-20201103192249-000122071b78 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.2 // indirect
	github.com/Azure/azure-event-hubs-go/v3 v3.6.1 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/go-amqp v1.0.1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.23 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.18.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector v0.42.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus v0.42.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.18.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.42.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ReneKroon/ttlcache/v2 v2.11.0 // indirect
	github.com/SAP/go-hdb v1.3.10 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/aerospike/aerospike-client-go/v6 v6.13.0 // indirect
	github.com/alecthomas/participle/v2 v2.0.0 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/aliyun/aliyun-log-go-sdk v0.1.51 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/antonmedv/expr v1.13.0 // indirect
	github.com/apache/thrift v0.18.1 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/aws/aws-sdk-go v1.44.323 // indirect
	github.com/aws/aws-sdk-go-v2 v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.12 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.33 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.32 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.59 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.39 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.18.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.31.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.2 // indirect
	github.com/aws/smithy-go v1.14.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/checkpoint-restore/go-criu/v5 v5.3.0 // indirect
	github.com/cilium/ebpf v0.7.0 // indirect
	github.com/cloudfoundry-incubator/uaago v0.0.0-20190307164349-8136b7bbe76e // indirect
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe // indirect
	github.com/cncf/xds/go v0.0.0-20230607035331-e9ce68804cb4 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/containerd/ttrpc v1.1.0 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denisenkom/go-mssqldb v0.12.2 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/digitalocean/godo v1.98.0 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v24.0.5+incompatible // indirect
	github.com/docker/go-connections v0.4.1-0.20210727194412-58542c764a11 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/dynatrace-oss/dynatrace-metric-utils-go v0.5.0 // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230111030713-bf00bc1b83b6 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/elastic/go-structform v0.0.10 // indirect
	github.com/emicklei/go-restful/v3 v3.10.1 // indirect
	github.com/envoyproxy/go-control-plane v0.11.1-0.20230524094728-9239064ad72f // indirect
	github.com/envoyproxy/protoc-gen-validate v0.10.1 // indirect
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-redis/redis/v7 v7.4.1 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cadvisor v0.47.3 // indirect
	github.com/google/flatbuffers v23.1.21+incompatible // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.5 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gophercloud/gophercloud v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gosnmp/gosnmp v1.35.0 // indirect
	github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd // indirect
	github.com/grobie/gomemcache v0.0.0-20180201122607-1f779c573665 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.2 // indirect
	github.com/hashicorp/consul/api v1.24.0 // indirect
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.2 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20230418003350-3067191c5197 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hetznercloud/hcloud-go v1.42.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/influxdata/go-syslog/v3 v3.0.1-0.20210608084020-ac565dc76ba6 // indirect
	github.com/influxdata/influxdb-observability/common v0.5.5 // indirect
	github.com/influxdata/influxdb-observability/influx2otel v0.5.5 // indirect
	github.com/influxdata/line-protocol/v2 v2.2.1 // indirect
	github.com/ionos-cloud/sdk-go/v6 v6.1.6 // indirect
	github.com/jaegertracing/jaeger v1.41.0 // indirect
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
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/leodido/ragel-machinery v0.0.0-20181214104525-299bdde78165 // indirect
	github.com/leoluk/perflib_exporter v0.2.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/linode/linodego v1.16.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/microsoft/ApplicationInsights-Go v0.4.4 // indirect
	github.com/miekg/dns v1.1.53 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20220423185008-bf980b35cac4 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mongodb-forks/digest v1.0.4 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/mostynb/go-grpc-compression v1.2.0 // indirect
	github.com/mrunalp/fileutils v0.5.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nginxinc/nginx-prometheus-exporter v0.8.1-0.20201110005315-f5a5f8086c19 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/cwlogs v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/docker v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/loki v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/signalfx v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.83.0 // indirect; indi72.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/winperfcounters v0.83.0 // indirect; indi72.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc4 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.3 // indirect
	github.com/opencontainers/selinux v1.10.1 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.2 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.11.0 // indirect
	github.com/prometheus/prometheus v0.44.0 // indirect
	github.com/prometheus/statsd_exporter v0.22.7 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/cors v1.9.0 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.15 // indirect
	github.com/seccomp/libseccomp-golang v0.9.2-0.20220502022130-f33da4d89646 // indirect
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.3 // indirect
	github.com/signalfx/gohistogram v0.0.0-20160107210732-1ccfd2ff5083 // indirect
	github.com/signalfx/golib/v3 v3.3.47 // indirect
	github.com/signalfx/sapm-proto v0.13.0 // indirect
	github.com/sijms/go-ora/v2 v2.7.11 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/snowflakedb/gosnowflake v1.6.23 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tidwall/gjson v1.10.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20210330154013-f5de75959ad5 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/vmware/govmomi v0.30.7 // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.etcd.io/bbolt v1.3.7 // indirect
	go.mongodb.org/atlas v0.33.0 // indirect
	go.mongodb.org/mongo-driver v1.12.1 // indirect
	go.opentelemetry.io/collector/semconv v0.83.0 // indirect; indir7.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.42.1-0.20230612162650-64be7e574a17 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.42.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.17.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.42.0 // indirect
	go.opentelemetry.io/otel v1.17.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.17.0 // indirect
	go.opentelemetry.io/otel/sdk v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.17.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/goleak v1.2.1 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/exp v0.0.0-20230626212559-97b1e661b5df // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/oauth2 v0.11.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/term v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.12.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gonum.org/v1/gonum v0.13.0 // indirect
	google.golang.org/api v0.136.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.27.4 // indirect
	k8s.io/apimachinery v0.27.4 // indirect
	k8s.io/client-go v0.27.4 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230501164219-8b0f38b5fd1f // indirect
	k8s.io/kubelet v0.27.4 // indirect
	k8s.io/utils v0.0.0-20230308161112-d77c459e9343 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace github.com/observiq/bindplane-agent/processor/resourceattributetransposerprocessor => ./processor/resourceattributetransposerprocessor

replace github.com/observiq/bindplane-agent/receiver/pluginreceiver => ./receiver/pluginreceiver

replace github.com/observiq/bindplane-agent/receiver/m365receiver => ./receiver/m365receiver

replace github.com/observiq/bindplane-agent/receiver/routereceiver => ./receiver/routereceiver

replace github.com/observiq/bindplane-agent/receiver/sapnetweaverreceiver => ./receiver/sapnetweaverreceiver

replace github.com/observiq/bindplane-agent/exporter/googlecloudexporter => ./exporter/googlecloudexporter

replace github.com/observiq/bindplane-agent/packagestate => ./packagestate

replace github.com/observiq/bindplane-agent/processor/metricstatsprocessor => ./processor/metricstatsprocessor

replace github.com/observiq/bindplane-agent/processor/removeemptyvaluesprocessor => ./processor/removeemptyvaluesprocessor

replace github.com/observiq/bindplane-agent/processor/throughputmeasurementprocessor => ./processor/throughputmeasurementprocessor

replace github.com/observiq/bindplane-agent/processor/samplingprocessor => ./processor/samplingprocessor

replace github.com/observiq/bindplane-agent/processor/maskprocessor => ./processor/maskprocessor

replace github.com/observiq/bindplane-agent/processor/logcountprocessor => ./processor/logcountprocessor

replace github.com/observiq/bindplane-agent/processor/metricextractprocessor => ./processor/metricextractprocessor

replace github.com/observiq/bindplane-agent/processor/logdeduplicationprocessor => ./processor/logdeduplicationprocessor

replace github.com/observiq/bindplane-agent/processor/spancountprocessor => ./processor/spancountprocessor

replace github.com/observiq/bindplane-agent/processor/datapointcountprocessor => ./processor/datapointcountprocessor

replace github.com/observiq/bindplane-agent/expr => ./expr

replace github.com/observiq/bindplane-agent/counter => ./counter

replace github.com/observiq/bindplane-agent/exporter/googlemanagedprometheusexporter => ./exporter/googlemanagedprometheusexporter

// Does not build with windows and only used in configschema executable
// Relevant issue https://github.com/mattn/go-ieproxy/issues/45
replace github.com/mattn/go-ieproxy v0.0.9 => github.com/mattn/go-ieproxy v0.0.1

// Replaces below this are required by datadog exporter in v0.83.0 https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v0.83.0/exporter/datadogexporter/go.mod#L266-L275

// openshift removed all tags from their repo, use the pseudoversion from the release-3.9 branch HEAD
replace github.com/openshift/api v3.9.0+incompatible => github.com/openshift/api v0.0.0-20180801171038-322a19404e37

// v0.47.x and v0.48.x are incompatible, prefer to use v0.48.x
replace github.com/DataDog/datadog-agent/pkg/proto => github.com/DataDog/datadog-agent/pkg/proto v0.48.0-beta.1

replace github.com/DataDog/datadog-agent/pkg/trace => github.com/DataDog/datadog-agent/pkg/trace v0.48.0-beta.1
