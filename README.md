# Disclaimer
While the `main` branch is stable, it is still under development and documentation may not fully reflect the current feature set. Please refer to documentation for your specific release.

# observIQ OpenTelemetry Collector

<center>

[![Action Status](https://github.com/observIQ/observiq-otel-collector/workflows/Build/badge.svg)](https://github.com/observIQ/observiq-otel-collector/actions)
[![Action Test Status](https://github.com/observIQ/observiq-otel-collector/workflows/Tests/badge.svg)](https://github.com/observIQ/observiq-otel-collector/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/observIQ/observiq-otel-collector)](https://goreportcard.com/report/github.com/observIQ/observiq-otel-collector)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Gosec](https://github.com/observIQ/observiq-otel-collector/actions/workflows/gosec.yml/badge.svg)](https://github.com/observIQ/observiq-otel-collector/actions/workflows/gosec.yml)

</center>

observIQ OpenTelemetry Collector is observIQ’s distribution of the [OpenTelemetry collector](https://github.com/open-telemetry/opentelemetry-collector). It provides a simple and unified solution to collect, refine, and ship telemetry data anywhere.

## Benefits

### Focused on usability
Increases the accessibility of OpenTelemetry by providing simplified installation scripts, tested example configurations, and end-to-end documentation making it easy to get started

### All the best parts of OpenTelemetry and more
Bundled with all core OpenTelemetry receivers, processors, and exporters as well as additional capabilities for monitoring complex or enterprise technologies not yet available in upstream releases
 
### Always Production-ready and fully-supported
Tested, verified, and supported by observIQ

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

To install using the installation script, you may run:

```sh
sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh
```

For more installation information see [installing on macOS](/docs/installation-mac.md).

#### Kubernetes

To deploy the collector on Kubernetes, further documentation can be found at our [observiq-otel-collector-k8s](https://github.com/observIQ/observiq-otel-collector-k8s) repository.

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

Here is an example `config.yaml` setup for hostmetrics on Google Cloud. To make sure your environment is set up with required prerequisites, see our [Google Cloud Exporter Prerequisites](/config/google_cloud_exporter/README.md) page. Further details for this GCP example can be found [here](/config/google_cloud_exporter/hostmetrics).

```yaml
# Receivers collect metrics from a source. The hostmetrics receiver will get
# CPU load metrics about the machine the collector is running on every minute.
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

# Processors are run on data between being received and being exported.
processors:
  # Resourcedetection is used to add a unique (host.name)
  # to the metric resource(s), allowing users to filter
  # between multiple systems.
  resourcedetection:
    detectors: ["system"]
    system:
      hostname_sources: ["os"]

  # Resourceattributetransposer is used to add labels to metrics.
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
  
  # Normalizesums smoothes out data points for more comprehensive visualizations.
  normalizesums:

  # The batch processor aggregates incoming metrics into a batch, releasing them if
  # a certain time has passed or if a certain number of entries have been aggregated.
  batch:

# Exporters send the data to a destination, in this case GCP.
exporters: 
  googlecloud:
    retry_on_failure:
      enabled: false
    metric:
      prefix: custom.googleapis.com

# Service specifies how to construct the data pipelines using the configurations above.
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

# Community

The observIQ OpenTelemetry Collector is an open source project. If you'd like to contribute, take a look at our [contribution guidelines](/CONTRIBUTING.md) and [developer guide](/docs/development.md). We look forward to building with you.

# How can we help?

If you need any additional help feel free to file a GitHub issue or reach out to us at support@observiq.com.
