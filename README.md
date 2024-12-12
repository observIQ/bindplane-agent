<a href="https://observiq.com">
  <p align="center">
    <picture>
      <source media="(prefers-color-scheme: light)" srcset="docs/assets/bindplane-logo.svg" width="auto" height="50">
      <source media="(prefers-color-scheme: dark)" srcset="docs/assets/bindplane-logo-dark.svg" width="auto" height="80">
      <img alt="BindPlane Logo" src="docs/assets/bindplane-logo.svg" width="auto" height="50">
    </picture>
  </p>
</a>

<p align="center">
  BindPlane standardizes your telemetry ingestion, processing, and shipping, by providing a unified, OTel-native pipeline.
</p>

<b>
  <p align="center">
    <a href="https://observiq.com/docs/getting-started/quickstart-guide">
      Get Started! &nbsp;👉&nbsp;
    </a>
  </p>
</b>

<b>
  <p align="center">
    <a href="https://observiq.com">Website</a>&nbsp;|&nbsp;
    <a href="https://observiq.com/docs/advanced-setup/installation">Docs</a>&nbsp;|&nbsp;
    <a href="https://observiq.com/docs/how-to-guides/routing-telemetry">How-to Guides</a>&nbsp;|&nbsp;
    <a href="https://observiq.com/docs/feature-guides/processors">Feature Guides</a>&nbsp;|&nbsp;
    <a href="https://observiq.com/blog">Blog</a>&nbsp;|&nbsp;
    <a href="https://observiq.com/mastering-opentelemetry">OTel Hub</a>&nbsp;|&nbsp;
    <a href="https://www.launchpass.com/bindplane">Slack</a>
  </p>
</b>

<!-- badges -->
<center>

[![Action Status](https://github.com/observIQ/bindplane-agent/workflows/Build/badge.svg)](https://github.com/observIQ/bindplane-agent/actions)
[![Action Test Status](https://github.com/observIQ/bindplane-agent/workflows/Tests/badge.svg)](https://github.com/observIQ/bindplane-agent/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/observIQ/bindplane-agent)](https://goreportcard.com/report/github.com/observIQ/bindplane-agent)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

</center>

<p align="center">
  <img src="docs/assets/bindplane-overview.webp" style="width:66%;height:auto">
</p>
<p align="center">
  <i>
    Learn how to connect BindPlane to telemetry <a href="https://observiq.com/docs/resources/sources">sources</a> and <a href="https://observiq.com/docs/resources/destinations">destinations</a>, and use <a href="https://observiq.com/docs/resources/processors">processors</a> to transform data.
  </i>
</p>

BindPlane is designed to be OpenTelemetry-first, with OpenTelemetry as its core framework, to create a unified toolset with data ownership. By providing a centralized management plane, BindPlane simplifies the development, implementation, management, and configuration of OpenTelemetry.

## Why BindPlane?

If you're managing telemetry at scale you'll run in to these problems sooner or later:

1. **Agent fatigue.** You'll end up managing dozens of proprietary agents all collecting and forwarding telemetry to different observability backends, which leads to performance issues and...
2. **Endless configuration files.** Even with GitOps practices you'll end up managing hundreds of configuration files for different sources, destinations, and processors that are written in proprietary languages instead leading to...
3. **Vendor lock-in.** BindPlane's primary focus is OpenTelemetry. It deploys and manages OpenTelemetry Collectors, uses OpenTelemetry Standards for terminology and configuration, and enables remote management with OpAMP.

That's why BindPlane will always be committed to these 4 core tenets.

### A collector you're used to

The BindPlane Agent is observIQ’s distribution of the [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector). It’s the first distribution to implement the [Open Agent Management Protocol](https://opentelemetry.io/docs/specs/opamp/) (OpAMP) and is designed to be fully managed with [BindPlane OP](https://observiq.com/solutions).

### Focus on usability

Increases the accessibility of OpenTelemetry by providing simplified installation scripts, tested example configurations, and end-to-end documentation making it easy to get started.

### All the best parts of OpenTelemetry and more

Bundled with all core OpenTelemetry receivers, processors, and exporters as well as additional capabilities for monitoring complex or enterprise technologies not yet available in upstream releases

### Always production-ready and fully-supported

Tested, verified, and supported by observIQ

## Getting Started

### BindPlane Cloud

BindPlane Cloud is the quickest way to get started. It offers managed infrastructure along with instant, free access for development projects and proofs of concept.

<a href="https://app.bindplane.com/signup"><img src="docs/assets/sign-up-bindplane-cloud.png" alt="Sign-up" width="200px"></a>

### BindPlane On Prem

You can also get started with BindPlane On Prem by hosting it yourself.

<a href="https://observiq.com/download"><img src="docs/assets/download-bindplane-on-prem.png" alt="Sign-up" width="200px"></a>

## BindPlane Agent Installation

BindPlane Agent uses OpAMP to connect to BindPlane OP.

### Linux

To install using the installation script, you may run:

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-agent/releases/latest/download/install_unix.sh)" install_unix.sh
```

To install directly with the appropriate package manager, see [installing on Linux](/docs/installation-linux.md).

### Windows

To install the BindPlane Agent on Windows run the Powershell command below to install the MSI with no UI.

```pwsh
msiexec /i "https://github.com/observIQ/bindplane-agent/releases/latest/download/observiq-otel-collector.msi" /quiet
```

Alternately, for an interactive installation [download the latest MSI](https://github.com/observIQ/bindplane-agent/releases/latest).

After downloading the MSI, simply double click it to open the installation wizard. Follow the instructions to configure and install the agent.

For more installation information see [installing on Windows](/docs/installation-windows.md).

### macOS

To install using the installation script, you may run:

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-agent/releases/latest/download/install_macos.sh)" install_macos.sh
```

For more installation information see [installing on macOS](/docs/installation-mac.md).

### Next Steps

Now that the agent is installed it is collecting basic metrics about the host machine printing them to the log. If you want to further configure your agent you may do so by editing the config file. To find your config file based on your OS reference the table below:

| OS      | Default Location                                              |
|:--------|:--------------------------------------------------------------|
| Linux   | /opt/observiq-otel-collector/config.yaml                      |
| Windows | C:\Program Files\observIQ OpenTelemetry Collector\config.yaml |
| macOS   | /opt/observiq-otel-collector/config.yaml                      |

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

## Example `config.yaml`

Here's a sample setup for `hostmetrics` on Google Cloud. To make sure your environment is set up with required prerequisites, see the [Google Cloud Exporter Prerequisites](/config/google_cloud_exporter/README.md) page. Further details for this GCP example can be found [here](/config/google_cloud_exporter/hostmetrics).

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

## Community

Have an idea to improve BindPlane Agent or BindPlane OP? Here's how you can help:

- Star this repo ⭐️ and follow us on [Twitter](https://x.com/bindplane).
- Upvote issues with 👍 so we know what to prioritize in the road map.
- [Create issues](https://github.com/observIQ/bindplane-agent/issues) when you feel something is missing or wrong.
- Join our [Slack Community](https://www.launchpass.com/bindplane), and ask us any questions there.

## Contributing

The BindPlane Agent is an open source project. If you'd like to contribute, take a look at our [contribution guidelines](/CONTRIBUTING.md) and [developer guide](/docs/development.md).

All sorts of contributions are **welcome and extremely helpful**. 🙌

## How can we help?

If you need any additional help feel free to reach out to us at support@observiq.com.
