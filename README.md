# observIQ Collector

The observIQ Collector is observIQ's distribution of the OpenTelemetry Collector. It provides everything you need to get setup and sending logs to 
[observIQ](https://observiq.com/).

### Configs

For sample configs, see the `./config` directory.
For general configuration help, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

For configuration options of a specific component, take a look at the README found in their respective module roots.

For a list of possible command line arguments to use with the collector, run the collector with the `--help` argument.
# Included Components
## Upstream components

### Receivers
* [filelogreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filelogreceiver)
* [otlpreceiver](https://github.com/open-telemetry/opentelemetry-collector/tree/main/receiver/otlpreceiver)
* [syslogreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/syslogreceiver)
* [tcplogreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/tcplogreceiver)
* [udpreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/udplogreceiver)
### Processors
* [groupbyattrsprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/groupbyattrsprocessor)
* [k8sprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/k8sprocessor)
* [attributesprocessor](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/attributesprocessor)
* [resourceprocessor](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/resourceprocessor)
* [batchprocessor](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/batchprocessor)
* [memorylimiter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/memorylimiter)
* [probabilisticsamplerprocessor](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/probabilisticsamplerprocessor)

### Exporters
* [fileexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/fileexporter)
* [otlpexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlpexporter)
* [otlphttpexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlphttpexporter)
* [observiqexporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/observiqexporter)
* [loggingexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/loggingexporter)
### Extensions
* [bearertokenauthextension](https://github.com/open-telemetry/opentelemetry-collector/tree/main/extension/bearertokenauthextension)
* [healthcheckextension](https://github.com/open-telemetry/opentelemetry-collector/tree/main/extension/healthcheckextension)
* [oidcauthextension](https://github.com/open-telemetry/opentelemetry-collector/tree/main/extension/oidcauthextension)
* [pprofextension](https://github.com/open-telemetry/opentelemetry-collector/tree/main/extension/pprofextension)
* [zpagesextension](https://github.com/open-telemetry/opentelemetry-collector/tree/main/extension/zpagesextension)
* [filestorage](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/storage/filestorage)

## Additional Components
These components are based in this repository:

### Receivers
* [logsreceiver](./receiver/logsreceiver)

### Extensions
* [orphandetectorextension](./extension/orphandetectorextension)

## Development

### Initial Setup

Clone this repository, and run `make install-tools`

### Building

To create a build for your current machine, run `make observiqcol`

To build for a specific architecture, see the [Makefile](./Makefile)

To build all targets, run `make build-all`

Build files will show up in the `./build` directory

### Running CI checks locally

The CI runs the `ci-checks` make target, which includes linting, testing, and checking documentation for misspelling.
CI also does a build of all targets (`make build-all`)

## Releasing
To release the collector, see [RELEASING.md](RELEASING.md)