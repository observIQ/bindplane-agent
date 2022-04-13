# observIQ OpenTelemetry Collector

<center>

[![Action Status](https://github.com/observIQ/observiq-otel-collector/workflows/Build/badge.svg)](https://github.com/observIQ/observiq-otel-collector/actions)
[![Action Test Status](https://github.com/observIQ/observiq-otel-collector/workflows/Tests/badge.svg)](https://github.com/observIQ/observiq-otel-collector/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/observIQ/observiq-otel-collector)](https://goreportcard.com/report/github.com/observIQ/observiq-otel-collector)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Gosec](https://github.com/observIQ/observiq-otel-collector/actions/workflows/gosec.yml/badge.svg)](https://github.com/observIQ/observiq-otel-collector/actions/workflows/gosec.yml)

</center>

## About

The observIQ OpenTelemetry Collector is observIQ's distribution of the OpenTelemetry Collector.

## Benefits

The observIQ OpenTelemetry Collector provides a customer centric experience based around the OpenTelemetry Collector. It provides customer friendly installation as well as curated configurations for common workflows.

## Quick Start

### Installation

#### Linux

To install using the installation script, you may run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_unix.sh)" install_unix.sh
```

To install directly with the appropriate package manager, see [installing on Linux](/docs/installation-linux.md).

#### Windows

To install the collector on Windows run the Powershell command below to install the MSI with no UI.
```pwsh
msiexec /i "https://github.com/observIQ/observiq-otel-collector/releases/latest/download/observiq-otel-collector.msi" /quiet
```

Alternately, for an interactive installation [download the latest MSI](https://github.com/observIQ/observiq-otel-collector/releases/latest).

After downloading the MSI, simply double click it to open the installation wizard. Follow the instructions to configure and install the collector.

For more installation information see [installing on Windows](/docs/installation-windows.md).

#### macOS

To install the collector on macOS use the following brew commands:

```sh
brew tap observiq/homebrew-observiq-otel-collector
brew update
brew install observiq/observiq-otel-collector/observiq-otel-collector
```

For more installation information see [installing on macOS](/docs/installation-mac.md).

### Next Steps

Now that the collector is installed it is collecting basic metrics about the host machine printing them to the log. If you want to further configure your collector you may do so by editing the config file. To find your config file based on your OS reference the table below:

| OS | Default Location |
| :--- | :---- |
| Linux | /opt/observiq-otel-collector/config.yaml |
| Windows | C:\Program Files\observIQ OpenTelemetry Collector\config.yaml |
| macOS | $(brew --prefix observiq/observiq-otel-collector/observiq-otel-collector)/config.yaml |

For more information on configuration see the [Configuration section](#configuration).

## Configuration

The observIQ OpenTelemetry Collector uses OpenTelemetry configuration.

For sample configs, see the [config](/config/) directory.
For general configuration help, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

For configuration options of a specific component, take a look at the README found in their respective module roots. For a list of currently supported components see [Included Components](#included-components).

For a list of possible command line arguments to use with the collector, run the collector with the `--help` argument.

### Included Components

#### Receivers

For supported receivers and their documentation see [receivers](/docs/receivers.md).

#### Processors

For supported processors and their documentation see [processors](/docs/processors.md).

#### Exporters

For supported exporters and their documentation see [exporters](/docs/exporters.md).

#### Extensions

For supported extensions and their documentation see [extensions](/docs/extensions.md).

## Example

Here is an example `config.yaml` setup for hostmetrics on Google Cloud. To make sure your environment is set up with required prerequisites, see our [Google Cloud Exporter Prerequisites](/config/google_cloud_exporter/README.md) page.

```yaml
receivers:
  hostmetrics:
    collection_interval: 60s
    scrapers:
      cpu:
      disk:
      load:
      filesystem:
      memory:
      network:
      paging:
      processes:


processors:
  # Resourcedetection is used to add a unique (host.name)
  # to the metric resource(s), allowing users to filter
  # between multiple systems.
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]

  resourceattributetransposer:
    operations:
      # Process metrics require unique metric labels, otherwise the Google
      # API will reject some metrics as "out of order" / duplicates.
      - from: "host.name"
        to: "hostname"
      - from: "process.pid"
        to: "pid"
      - from: "process.executable.name"
        to: "binary"

  normalizesums:

  batch:

exporters: 
  googlecloud:
    retry_on_failure:
      enabled: false
    metric:
      prefix: custom.googleapis.com

service:
  pipelines:
    metrics:
      receivers:
      - hostmetrics
      processors:
      - resourcedetection
      - resourceattributetransposer
      - normalizesums
      - batch
      exporters:
      - googlecloud
```

Further details for this example can be found [here](/config/google_cloud_exporter/hostmetrics).

# Community

The observIQ OpenTelemetry Collector is an open source project. If you'd like to contribute, take a look at our [contribution guidelines](/CONTRIBUTING.md) and [developer guide](/docs/development.md). We look forward to building with you.

# How can we help?

If you need any additional help feel free to file a GitHub issue or reach out to us at support@observiq.com.
