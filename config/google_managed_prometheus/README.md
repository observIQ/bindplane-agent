# Google Managed Prometheus

The Prometheus Exporter can be used to send metrics to [Google Managed Prometheus](https://cloud.google.com/stackdriver/docs/managed-prometheus) (GMP). This directory contains sub directories with usercase specific configurations, all of which are compatible with GMP.

## Google Cloud APIs

Enable the following APIs.
- Cloud Metrics

To learn more about enabling APIs, see the [documentation](https://cloud.google.com/endpoints/docs/openapi/enable-api).

## Prometheus

It is assumed that you have Google's [Prometheus fork](https://github.com/GoogleCloudPlatform/prometheus/tree/main) running within your
environment.

An example Prometheus configuration:

```yaml
global:
  external_labels:
    project_id: otel-managed-prometheus
    location: us-east1-b
    namespace: otelprom

scrape_configs:
  # Collector metrics
  - job_name: "otelcol"
    static_configs:
      - targets:
        # Prometheus exporter metrics
        - "localhost:9000"
```
