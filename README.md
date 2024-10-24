# BindPlane Agent

<center>

[![Action Status](https://github.com/observIQ/bindplane-agent/workflows/Build/badge.svg)](https://github.com/observIQ/bindplane-agent/actions)
[![Action Test Status](https://github.com/observIQ/bindplane-agent/workflows/Tests/badge.svg)](https://github.com/observIQ/bindplane-agent/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/observIQ/bindplane-agent)](https://goreportcard.com/report/github.com/observIQ/bindplane-agent)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

</center>

The BindPlane Agent is observIQ’s custom distribution of the [OpenTelemetry collector](https://github.com/open-telemetry/opentelemetry-collector) built using the [OpenTelemetry collector builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder). The [OpenTelemetry supervisor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/cmd/opampsupervisor) is used to manage the collector (start/stop, apply configurations, etc) and communicate between the collector and [BindPlane OP](https://observiq.com/). To get started, follow our [Quickstart Guide](https://observiq.com/docs/getting-started/quickstart-guide).

## Benefits

### Focused on usability
Increases the accessibility of OpenTelemetry by providing simplified installation scripts and end-to-end documentation making it easy to get started

### All the best parts of OpenTelemetry and more
Bundled with all core OpenTelemetry receivers, processors, and exporters as well as additional capabilities for monitoring complex or enterprise technologies not yet available in upstream releases
 
### Always Production-ready and fully-supported
Tested, verified, and supported by observIQ

## Quick Start

You'll need a BindPlane server in order to run the collector with the supervisor. Namely you'll need the server's OpAMP endpoint and a secret key the agent can use to authenticate with. For more information, see the [supervisor documentation](./docs/supervisor.md).

### Installation

#### Linux

To install using the installation script, you may run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-agent/releases/latest/download/install_unix.sh)" install_unix.sh
```

To install directly with the appropriate package manager, see [installing on Linux](/docs/installation-linux.md).

#### Windows

To install the BindPlane Agent on Windows, run the Powershell command below to install the MSI with no UI.
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

With the agent installed, you can use BindPlane to create a configuration and begin monitoring.

You can edit the supervisor config as needed for communicating with BindPlane and managing the agent. To find your config file based on your OS reference the table below:

| OS      | Default Location                                                  |
|:--------|:------------------------------------------------------------------|
| Linux   | /opt/observiq-otel-collector/supervisor.yaml                      |
| Windows | C:\Program Files\observIQ OpenTelemetry Collector\supervisor.yaml |
| macOS   | /opt/observiq-otel-collector/supervisor.yaml                      |

For more information on supervisor configuration see the [supervisor documentation](./docs/supervisor.md).

## Configuration

### Agent

The BindPlane Agent uses OpenTelemetry configuration. Running the agent with the supervisor requires receiving the agent's config from an OpAMP management server, namely BindPlane in this context.

Specific information on managing agents, creating a configuration, and rolling out configs can be found in [BindPlane documentation](https://observiq.com/docs/getting-started/quickstart-guide).

For general configuration help, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

For configuration options of a specific component, take a look at the README found in their respective module roots. For a list of currently supported components see [Included Components](#included-components).

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

# Community

The BindPlane Agent is an open source project. If you'd like to contribute, take a look at our [contribution guidelines](/CONTRIBUTING.md) and [developer guide](/docs/development.md). We look forward to building with you.

# How can we help?

If you need any additional help feel free to file a GitHub issue or reach out to us at support@observiq.com.
