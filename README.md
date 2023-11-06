# BindPlane Agent

<center>

[![Action Status](https://github.com/observIQ/bindplane-agent/workflows/Build/badge.svg)](https://github.com/observIQ/bindplane-agent/actions)
[![Action Test Status](https://github.com/observIQ/bindplane-agent/workflows/Tests/badge.svg)](https://github.com/observIQ/bindplane-agent/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/observIQ/bindplane-agent)](https://goreportcard.com/report/github.com/observIQ/bindplane-agent)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Gosec](https://github.com/observIQ/bindplane-agent/actions/workflows/gosec.yml/badge.svg)](https://github.com/observIQ/bindplane-agent/actions/workflows/gosec.yml)

</center>

BindPlane Agent is observIQ’s distribution of the [OpenTelemetry collector](https://github.com/open-telemetry/opentelemetry-collector). It provides a simple and unified solution to collect, refine, and ship telemetry data anywhere.

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
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-agent/releases/latest/download/install_unix.sh)" install_unix.sh
```

To install directly with the appropriate package manager, see [installing on Linux](/docs/installation-linux.md).

#### Windows

To install the BindPlane Agent on Windows run the Powershell command below to install the MSI with no UI.
```pwsh
msiexec /i "https://github.com/observIQ/bindplane-agent/releases/latest/download/observiq-otel-collector.msi" /quiet
```

Alternately, for an interactive installation [download the latest MSI](https://github.com/observIQ/bindplane-agent/releases/latest).

After downloading the MSI, simply double click it to open the installation wizard. Follow the instructions to configure and install the agent.

For more installation information see [installing on Windows](/docs/installation-windows.md).

#### macOS

To install using the installation script, you may run:

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-agent/releases/latest/download/install_macos.sh)" install_macos.sh
```

For more installation information see [installing on macOS](/docs/installation-mac.md).

### Next Steps

Now that the agent is installed it is collecting basic metrics about the host machine printing them to the log. If you want to further configure your agent you may do so by editing the config file. To find your config file based on your OS reference the table below:

| OS | Default Location |
| :--- | :---- |
| Linux | /opt/observiq-otel-collector/config.yaml |
| Windows | C:\Program Files\observIQ OpenTelemetry Collector\config.yaml |
| macOS | /opt/observiq-otel-collector/config.yaml |

For more information on configuration see the [Configuration section](#configuration).

## Configuration

The BindPlane Agent uses OpenTelemetry configuration.

For sample configs, see the [config](/config/) directory.
For general configuration help, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

For configuration options of a specific component, take a look at the README found in their respective module roots. For a list of currently supported components see [Included Components](#included-components).

For a list of possible command line arguments to use with the agent, run the agent with the `--help` argument.

### Included Components

#### Receivers

For supported receivers and their documentation see [receivers](/docs/receivers.md).

#### Processors

For supported processors and their documentation see [processors](/docs/processors.md).

#### Exporters

For supported exporters and their documentation see [exporters](/docs/exporters.md).

#### Extensions

For supported extensions and their documentation see [extensions](/docs/extensions.md).

#### Connectors

For supported connectors and their documentation see [connectors](/docs/connectors.md).

## Example

Here is an example `config.yaml` setup for hostmetrics on Google Cloud. To make sure your environment is set up with required prerequisites, see our [Google Cloud Exporter Prerequisites](/config/google_cloud_exporter/README.md) page. Further details for this GCP example can be found [here](/config/google_cloud_exporter/hostmetrics).

```yaml
# Receivers collect metrics from a source. The hostmetrics receiver will get
# CPU load metrics about the machine the agent is running on every minute.
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

# Exporters send the data to a destination, in this case GCP.
exporters: 
  googlecloud:

# Service specifies how to construct the data pipelines using the configurations above.
service:
  pipelines:
    metrics:
      receivers: [hostmetrics]
      exporters: [googlecloud]
```

# Community

The BindPlane Agent is an open source project. If you'd like to contribute, take a look at our [contribution guidelines](/CONTRIBUTING.md) and [developer guide](/docs/development.md). We look forward to building with you.

# How can we help?

If you need any additional help feel free to file a GitHub issue or reach out to us at support@observiq.com.
