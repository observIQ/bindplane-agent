# observIQ OpenTelemetry Collector

The observIQ OpenTelemetry Collector is observIQ's distribution of the OpenTelemetry Collector.

## Installation

* [Installing on Windows](/docs/installation-windows.md)
* [Installing on Linux](/docs/installation-linux.md)
### Configs

For sample configs, see the `./config` directory.
For general configuration help, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

For configuration options of a specific component, take a look at the README found in their respective module roots.

For a list of possible command line arguments to use with the collector, run the collector with the `--help` argument.



# Included Components

### Receivers

For supported receivers and their documentation see [receivers](/docs/receivers.md).

### Processors

For supported processors and their documentation see [processors](/docs/processors.md).

### Exporters

For supported exporters and their documentation see [exporters](/docs/exporters.md).

### Extensions

For supported extensions and their documentation see [extensions](/docs/extensions.md).

## Development

### Initial Setup

Clone this repository, and run `make install-tools`

### Building

To create a build for your current machine, run `make collector`

To build for a specific architecture, see the [Makefile](./Makefile)

To build all targets, run `make build-all`

Build files will show up in the `./dist` directory

### Running CI checks locally

The CI runs the `ci-checks` make target, which includes linting, testing, and checking documentation for misspelling.
CI also does a build of all targets (`make build-all`)

## Releasing
To release the collector, see [RELEASING.md](RELEASING.md)
