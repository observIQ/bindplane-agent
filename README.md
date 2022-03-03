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

For supported exports and their documentation see [exporters](/docs/exporters.md).

### Extensions

For supported extensions and their documentation see [extensions](/docs/extensions.md).

## Additional Components
These components are based in this repository:

### Receivers
* [logsreceiver](./receiver/logsreceiver)

## Development

### Initial Setup

Clone this repository, and run `make install-tools`

### Building

To create a build for your current machine, run `make agent_manager`

To build for a specific architecture, see the [Makefile](./Makefile)

To build all targets, run `make build-all`

Build files will show up in the `./build` directory

### Running CI checks locally

The CI runs the `ci-checks` make target, which includes linting, testing, and checking documentation for misspelling.
CI also does a build of all targets (`make build-all`)

## Releasing
To release the collector, see [RELEASING.md](RELEASING.md)
