# Sampling Processor

Supported pipelines: logs, metrics, traces

This processor samples OTLP payloads based on the configured `drop_ratio`.

## Configuration

The following options may be configured:
- `drop_ratio` (default: 0.5): The ratio of payloads that are dropped. Values between `0.0` and `1.0`. Values closer to `1.0` mean any individual payload is more likely to be dropped.

### Example configuration

```yaml
processors:
  sampling:
    drop_ratio: 0.5
```

