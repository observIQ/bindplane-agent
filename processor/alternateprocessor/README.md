# Alternate Processor

This processor is used in high throughput scenarios to divert telemetry to another specified pipeline.

## Supported pipelines

- Metrics
- Logs
- Traces

## How It Works

1. The user configures the `alternate` processor in pipeline and a route receiver in their desired pipeline.
2. The processor keeps track of the byte size of incoming telemetry over an rate and interval.
3. If the rate of ingestion into the processor is less than or equal to the rate, then the telemetry will go further in the original pipeline
4. If the rate of ingestion is greater than the rate, then the telemetry will go into the specified route.

## Configuration

| Field        | Type     | Default | Description |
| ---          | ---      | ---     | ---         |
| metrics |  | | |
| metrics.enabled | bool | false | Whether or not to enable metrics alternate flow |
| metrics.rate | string | | Required if enabled to specify the desired rate to compare ingestion against. In the format of `value unit/duration` i.e. `30 MB/sec` |
| metrics.route | string | | The name of the route to send data to in the case that ingestion exceeds the specified rate |
| metrics.aggregation_interval| duration| 10s | How often it is desired to check the ingestion average over the interval |
| logs | | | |
| logs.enabled | bool | false | Whether or not to enable logs alternate flow |
| logs.rate | string | | Required if enabled to specify the desired rate to compare ingestion against. In the format of `value unit/duration` i.e. `30 MB/sec` |
| logs.route | string | | The name of the route to send data to in the case that ingestion exceeds the specified rate |
| logs.aggregation_interval| duration| 10s | How often it is desired to check the ingestion average over the interval |
| traces| | | |
| traces.enabled | bool | false | Whether or not to enable traces alternate flow |
| traces.rate | string | | Required if enabled to specify the desired rate to compare ingestion against. In the format of `value unit/duration` i.e. `30 MB/sec` |
| traces.route | string | | The name of the route to send data to in the case that ingestion exceeds the specified rate |
| traces.aggregation_interval| duration| 10s | How often it is desired to check the ingestion average over the interval |

### Example Config

The following config is ingesting logs from the path `/var/log/apache.log` and is exported via the `googlecloud` exporter. The `alternate` processor will divert logs to a file located at `/var/log/hydrate.log` if the rate of ingestion for logs exceeds `100 MiB/sec`.

```yaml
receivers:
  filelog:
    include: [/var/log/apache.log]
    start_at: end
  route/alternate-route:

exporters:
  file/1:
    path: /var/log/hydrate.log
  googlecloud:

processors:
  alternate:
    logs:
      enabled: true
      rate:  100 MiB/sec
      route: alternate-route

service:
  pipelines:
    logs/original:
      receivers:
      - filelog
      processors:
      - alternate
      exporters:
      - googlecloud
    logs/alternate:
      receivers:
      - route/logs-route
      exporters:
      - file/1
```
