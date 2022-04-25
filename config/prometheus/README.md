# Prometheus Exporter Prerequisites

## Configure Prometheus to scrape the collector

Using the basic configuration provided in [these docs](https://prometheus.io/docs/prometheus/latest/getting_started/#configuring-prometheus-to-monitor-itself) instead of:

```
static_configs:
    - targets: ['localhost:<PROMETHEUS_PORT>']
```
edit the config target the collector:
```
static_configs:
    - targets: ['localhost:<COLLECTOR_PORT>']
```
