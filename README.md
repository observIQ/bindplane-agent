<a href="https://observiq.com">
  <p align="center">
    <picture>
      <source media="(prefers-color-scheme: light)" srcset="https://res.cloudinary.com/du4nxa27k/image/upload/v1734023130/observiq-logo-dark_i3ycyh.svg" width="auto" height="50">
      <source media="(prefers-color-scheme: dark)" srcset="https://res.cloudinary.com/du4nxa27k/image/upload/v1734023130/observiq-logo-white_fdr6y8.svg" width="auto" height="80">
      <img alt="BindPlane Logo" src="https://res.cloudinary.com/du4nxa27k/image/upload/v1734023130/observiq-logo-dark_i3ycyh.svg" width="auto" height="50">
    </picture>
  </p>
</a>

<p align="center">
  The BindPlane Distro for OpenTelemetry Collector (BDOT Collector) is observIQ’s distribution of the upstream <a href="https://github.com/open-telemetry/opentelemetry-collector">OpenTelemetry Collector</a>. It’s the first distribution to implement the <a href="https://opentelemetry.io/docs/specs/opamp/">Open Agent Management Protocol</a> (OpAMP) and is designed to be fully managed with <a href="https://observiq.com/">BindPlane Telemetry Pipeline</a>.
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
<p align="center">
  <a href="https://github.com/observIQ/bindplane-otel-collector/actions">
    <img src="https://github.com/observIQ/bindplane-otel-collector/workflows/Build/badge.svg" alt="Action Status">
  </a>
  <a href="https://github.com/observIQ/bindplane-otel-collector/actions">
    <img src="https://github.com/observIQ/bindplane-otel-collector/workflows/Tests/badge.svg" alt="Action Test Status">
  </a>
  <a href="https://goreportcard.com/report/github.com/observIQ/bindplane-otel-collector">
    <img src="https://goreportcard.com/badge/github.com/observIQ/bindplane-otel-collector" alt="Go Report Card">
  </a>
  <a href="https://opensource.org/licenses/Apache-2.0">
    <img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License">
  </a>
</p>

<p align="center">
  <img src="https://res.cloudinary.com/du4nxa27k/image/upload/v1734000985/bindplane-overview_ke8xmq.webp" style="width:66%;height:auto">
</p>
<p align="center">
  <i>
    Learn how to connect BindPlane Distro for OpenTelemetry Collector to telemetry <a href="https://observiq.com/docs/resources/sources">sources</a> and <a href="https://observiq.com/docs/resources/destinations">destinations</a>, and use <a href="https://observiq.com/docs/resources/processors">processors</a> to transform data.
  </i>
</p>

## Why BindPlane Distro for OpenTelemetry Collector?

If you're managing telemetry at scale you'll run in to these problems sooner or later:

1. **Agent fatigue.** You'll manage endless proprietary agents and OpenTelemetry Collectors that collect and forward telemetry to different observability backends, leading to performance issues.
2. **Endless configuration files.** Even with GitOps practices you'll have to manage hundreds of configuration files for different sources, destinations, and processors written in either proprietary languages or YAML.
3. **High complexity.** OpenTelemetry's complexity and learning curve make it difficult to implement, manage, and re-point telemetry without a centralized management plane like BindPlane Telemetry Pipeline to standardize telemetry ingestion, processing, and shipping, with a unified, OpenTelemetry-native pipeline.

### An OpenTelemetry Collector you're used to

The BDOT Collector is observIQ’s distribution of the upstream [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector). It’s the first distribution to implement the [Open Agent Management Protocol](https://opentelemetry.io/docs/specs/opamp/) (OpAMP) and is designed to be fully managed with [BindPlane Telemetry Pipeline](https://observiq.com/solutions).

### Focused on usability

Increases the accessibility of OpenTelemetry by providing simplified installation scripts, tested example configurations, and end-to-end documentation making it easy to get started.

### All the best parts of OpenTelemetry and more

Bundled with all core OpenTelemetry receivers, processors, and exporters as well as additional capabilities for monitoring complex or enterprise technologies not yet available in upstream releases

### Always production-ready and fully-supported

Tested, verified, and supported by observIQ.

## Getting Started

Follow the [Getting Started](/docs/getting-started.md) guide for more detailed installation instructions, or view the list of [Supported Operating System Versions](https://observiq.com/docs/advanced-setup/installation/install-agent).

To continue with the quick start, follow along below.

### Linux

Install BDOT Collector using the installation script below.

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-otel-collector/releases/latest/download/install_unix.sh)" install_unix.sh
```

To install directly with the appropriate package manager, and how to configure OpAMP, see [installing on Linux](/docs/installation-linux.md).

### Windows

To install the BDOT Collector on Windows run the Powershell command below to install the MSI with no UI.

```pwsh
msiexec /i "https://github.com/observIQ/bindplane-otel-collector/releases/latest/download/observiq-otel-collector.msi" /quiet
```

Alternately, for an interactive installation [download the latest MSI](https://github.com/observIQ/bindplane-otel-collector/releases/latest).

After downloading the MSI, simply double click it to open the installation wizard. Follow the instructions to configure and install the BDOT Collector.

For more installation information, and how to configure OpAMP, see [installing on Windows](/docs/installation-windows.md).

### macOS

Install BDOT Collector using the installation script below.

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh
```

For more installation information, and how to configure OpAMP, see [installing on macOS](/docs/installation-mac.md).

## Next Steps

### BDOT Collector default `config.yaml`

With the BDOT Collector installed, it will start collecting basic metrics about the host machine printing them to the log. To further configure your collector edit the `config.yaml` file just like you would an OpenTelemetry Collector. To find your `config.yaml` file based on your operating system, reference the table below:

| OS      | Default Location                                              |
|:--------|:--------------------------------------------------------------|
| Linux   | `/opt/observiq-otel-collector/config.yaml`                      |
| Windows | `C:\Program Files\observIQ OpenTelemetry Collector\config.yaml` |
| macOS   | `/opt/observiq-otel-collector/config.yaml`                      |

For more information on configuration see the [Configuration section](#configuration).

### Manage BDOT Collector with BindPlane Telemetry Pipeline via OpAMP

Improving developer experience with OpenTelemetry is observIQ's primary focus. We're building BindPlane Telemetry Pipeline to help deploy and manage OpenTelemetry Collectors at scale, but retain core OpenTelemetry Standards for terminology and configuration, with the added benefit of enabling remote management with OpAMP.

The BDOT Collector can be configured as an OpenTelemetry Collector that is managed by the BindPlane Telemetry Pipeline via OpAMP. BindPlane is designed to be OpenTelemetry-first, with OpenTelemetry as its core framework. By providing a centralized management plane, it simplifies the development, implementation, management, and configuration of OpenTelemetry.

For more information on managing collectors via OpAMP see the [Connecting to BindPlane Telemetry Pipeline with OpAMP section](#connecting-to-bindplane-telemetry-pipeline-with-opamp).

## Configuration

The BDOT Collector uses OpenTelemetry Collector configuration.

For sample configs, see the [config](/config/) directory.
For general configuration help, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

For configuration options of a specific component, take a look at the README found in their respective module roots. For a list of currently supported components see [Included Components](#included-components).

For a list of possible command line arguments to use with the BDOT Collector, run the collector with the `--help` argument.

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

### Example `config.yaml`

Here's a sample setup for `hostmetrics` on Google Cloud. To make sure your environment is set up with required prerequisites, see the [Google Cloud Exporter Prerequisites](/config/google_cloud_exporter/README.md) page. Further details for this GCP example can be found [here](/config/google_cloud_exporter/hostmetrics).

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

## Connecting to BindPlane Telemetry Pipeline with OpAMP

BindPlane is designed to be OpenTelemetry-first, with OpenTelemetry as its core framework, to create a unified toolset with data ownership. By providing a centralized management plane, it simplifies the development, implementation, management, and configuration of OpenTelemetry.

To learn more about configuring OpAMP, see [OpAMP Configuration](/docs/opamp.md), or get started with BindPlane Telemetry Pipeline below.

### BindPlane Cloud

BindPlane Cloud is the quickest way to get started with OpenTelemetry-native telemetry pipelines. It offers managed infrastructure along with instant, free access for development projects and proofs of concept.

<a href="https://app.bindplane.com/signup"><img src="https://res.cloudinary.com/du4nxa27k/image/upload/v1734001746/sign-up-bindplane-cloud_tzhj8r.png" alt="Sign-up" width="200px"></a>

### BindPlane On Prem

You can also get started with BindPlane On Prem for free by hosting it yourself.

<a href="https://observiq.com/download"><img src="https://res.cloudinary.com/du4nxa27k/image/upload/v1734000970/download-bindplane-on-prem_rhdrme.png" alt="Download" width="200px"></a>

## Community

Have an idea to improve the BindPlane Distro for OpenTelemetry Collector? Here's how you can help:

- Star this repo ⭐️ and follow us on [Twitter](https://x.com/bindplane).
- Upvote issues with 👍 so we know what to prioritize in the road map.
- [Create issues](https://github.com/observIQ/bindplane-otel-collector/issues) when you feel something is missing or wrong.
- Join our [Slack Community](https://www.launchpass.com/bindplane), and ask us any questions there.

## Contributing

The BindPlane Distro for OpenTelemetry Collector is an open source project. If you'd like to contribute, take a look at our [contribution guidelines](/CONTRIBUTING.md) and [developer guide](/docs/development.md).

All sorts of contributions are **welcome and extremely helpful**. 🙌

## How can we help?

If you need any additional help feel free to reach out to us at support@observiq.com.
