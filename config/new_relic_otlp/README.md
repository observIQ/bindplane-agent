# New Relic OTLP Exporter Prerequisites

New Relic supports native ingestion of OTLP telemetry. The [Open Telemetry OTLP gRPC exporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlpexporter) can be used to send metrics, logs, and traces to New Relic.

## Prerequisites

**Endpoint**

Update the exporter configuration's `endpoint` value with the New Relic OTLP endpoint you wish to send metrics to.
- https://otlp.nr-data.net
- https://otlp.eu01.nr-data.net
- https://gov-otlp.nr-data.net

You can read more [here](https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/opentelemetry-setup/).

**API Key**

You must have an API Key capable of writing telemetry to the OTLP endpoint. You can read more [here](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/).

## Example Config

Using the `endpoint` and `api-key` values from the Prerequisites step, configure the exporter.

```yaml
exporters:
  otlp:
    endpoint: https://otlp.nr-data.net
    headers:
      api-key: 00000-00000-00000
    tls:
      insecure: false
```