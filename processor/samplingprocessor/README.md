# Sampling Processor

This processor samples incoming OTLP objects and drops those objects based on a configured `drop_ratio`.

## Supported pipelines

- Logs
- Metrics
- Traces

## How it works

1. The user configures the processor in their pipeline with a `drop_ratio` and condition that is the desired.
2. If an incoming log matches the `condition` expression, the remaining steps are performed on it. Any log record that does not match the `condition` gets forwarded through the pipeline regardless of the `drop_ratio`.
3. A number between 0 and 1 will be randomly generated for each piece incoming telemetry data.
4. If the generated number is less than or equal to the `drop_ratio`, then the telemetry data is dropped.
5. If the generated number is greater than the `drop_ratio`, then the telemetry data makes it further in the pipeline.

## Configuration

The following options may be configured:

| Field      | Type   | Default | Description                                                                                                                                                                 |
| --         | --     | --      | --                                                                                                                                                                          |
| drop_ratio | float  | 0.5     | The ratio of payload objects that are dropped. Values between `0.0` and `1.0`. Values closer to `1.0` mean any individual object in a payload is more likely to be dropped. |
| condition  | string | `true`  | An [OTTL] expression used to match which log records to sample from. All paths in the [log context] are available to reference. All [converters] are available to use.      |

### Example Configuration

The following config is an example configuration of the `sampling` processor with defaults in a logs pipeline sending to the `logging` exporter.

```yaml
receivers:
  filelog:
    include: [/tmp/example/apache.log]
processors:
  sampling:
    drop_ratio: 0.5
exporters:
  logging:

service:
  pipelines:
    logs:
      receivers: [filelog]
      processors: [sampling]
      exporters: [logging]
```

## How to

### Sample 75% of incoming telemetry

The following configuration will drop 75% of incoming `metrics`, `logs` or `traces`.

```yaml
processors:
  sampling:
    drop_ratio: 0.75
```

### Drop all incoming telemetry

The following configuration will drop 100% of incoming `metrics`, `logs`, or `traces`. Essentially dropping all data.

```yaml
processors:
  sampling:
    drop_ratio: 1.0
```

### Sample 50% of incoming telemetry where body field "ID" equals 1

The following configuration will drop 50% of incoming telemetry where the body field "ID" equals 1.

```yaml
processors:
  sampling:
    drop_ratio: 0.5
    condition: (body["ID"] == 1)
```
