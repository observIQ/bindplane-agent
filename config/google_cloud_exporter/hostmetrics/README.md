# Host Metrics with Google Cloud

The Host Metrics Receiver can be used to send system metrics from an agent host to Google Cloud Monitoring. See the [documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/hostmetricsreceiver) for more details.

## Prerequisites

See the [prerequisites](../prerequisites.md) doc for Google Cloud prerequisites.

## Configuration

An example configuration is located [here](./config.yaml).

1. Copy [config.yaml](./config.yaml) to `/opt/observiq-otel-collector/config.yaml`
2. Restart the collector: `sudo systemctl restart observiq-otel-collector`

## Setup

This example assumes that [node exporter](https://github.com/prometheus/node_exporter) is running on the host system, on port 9100. Generally, all prometheus style exporters are supported, such as the [aerospike exporter](https://github.com/aerospike/aerospike-prometheus-exporter).

If you wish to scrape metrics from external systems, update `config.yaml`'s prometheus scrape config with a list of IP addresses:
```yaml
- targets:
    - '127.0.0.1:9100'
    - '10.128.1.30:9100'
    - '10.128.1.31:9100'
    - '10.128.1.32:9100'
    - '10.128.1.33:9100'
```

## Process metrics

The host metrics receiver supports per process cpu, memory, disk, and network metrics. This feature requires elevetated privileges. On Linux, you must update the service configuration to run the collector as the root user. See the [installation documentation](https://github.com/observIQ/observiq-otel-collector/blob/main/docs/installation-linux.md#configuring-the-collector) for instructions.

## Metrics

Metrics can be found with the `custom.googleapis.com/system` prefix.

Example MQL query for `cpu.time`:
```
fetch global
| metric 'custom.googleapis.com/system.cpu.time'
| align rate(1m)
| every 1m
```
