# Log Count Processor
This processor is used to convert the number of logs received during an interval into a metric.

## Supported pipelines
- Logs

## How It Works
1. The user configures the processor in their logs pipeline and then a matching receiver with the same name in their metrics pipeline. This will bridge the two different telemetry pipelines.
2. If any incoming logs match the `match` expression, they are counted and dimensioned by their `attributes`. Regardless of match, all logs are sent to the next component in the pipeline.
3. After each configured interval, the observed log counts are converted to a gauge metric. This metric is sent to the matching receiver and the counter is reset.


## Configuration
| Field        | Type     | Default | Description |
| ---          | ---      | ---     | ---         |
| match        | string   | `true`  | A boolean [expression](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md) used to match which logs to count. By default, all logs are counted. |
| interval     | duration | `1m`    | The interval at which metrics are created. The counter will reset after each interval. |
| metric_name  | string   | `log.count` | The name of the metric created. |
| metric_unit  | string   | `{logs}`    | The unit of the metric created. |
| attributes   | map      | `{}`        | The mapped attributes of the metric created. Each key is an attribute name. Each value is an [expression](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md) that extracts data from the log. |

### Example Config
The following config is an example configuration of the log count processor using default values. In this example, logs are collected from a file, sent to the processor to be counted, and then consumed by the logging exporter. After each minute, the log counts are converted to metrics and sent to the matching receiver in the metrics pipeline, which then forwards to the Google Cloud exporter.
```yaml
receivers:
    filelog:
        include: [./example/apache.log]
    logcount:
processors:
    batch:
    logcount:
exporters:
    googlecloud:
    logging:

service:
    pipelines:
        logs:
            receivers: [filelog]
            processors: [batch, logcount]
            exporters: [logging]
        metrics:
            receivers: [logcount]
            processors: [batch]
            exporters: [googlecloud]
```

## Expression Language
In order to match or extract values from logs, the following `keys` are reserved and can be used to traverse the logs data model.

| Key          | Description |
| ---          | ---   |
| `body`       | Used to access the body of the log. |
| `attributes` | Used to access the attributes of the log. |
| `resource`   | Used to access the resource of the log. |
| `severity`   | Used to access the severity text of the log. |

In order to access embedded values, use JSON dot notation. For example, `body.example.field` can be used to access a field two levels deep on the log body. 

However, if a key already possesses a literal dot, users will need to use bracket notation to access that field. For example, when the field `service.name` exists on the log's resource, users will need to use `resource["service.name"]` to access this value.

For more information about syntax and available operators, see the [Expression Language Definition](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md).

## How To
### Match specific logs
The following configuration adds a match expression that will count only logs with a 4xx status code. In this example, incoming logs are expected to have a `status` attribute.
```yaml
processors:
    logcount:
        match: attributes.status startsWith "4"
```

### Extract metric attributes
The following configuration extracts the status and endpoint values from the body of a log. These values are used as metric attributes. For each unique combination observed, a unique metric count is created.
```yaml
processors:
    logcount:
        attributes:
            status_code: body.status
            endpoint: body.endpoint
```

### Use multiple processors
The following configuration uses multiple log count processors to send to different metric pipelines with different exporters. Each processor will only send metrics to the corresponding receiver with the same name and id. For example, `logcount/info` will create metrics for info level logs and send them only to the `logcount/info` receiver.
```yaml
receivers:
    filelog:
        include: [./example/apache.log]
    logcount/info:
    logcount/error:
processors:
    batch:
    logcount/info:
        match: severity contains "INFO"
    logcount/error:
        match: severity contains "ERROR"
exporters:
    googlecloud/info:
    googlecloud/error:
    logging:

service:
    pipelines:
        logs:
            receivers: [filelog]
            processors: [batch, logcount/info, logcount/error]
            exporters: [logging]
        metrics/info:
            receivers: [logcount/info]
            processors: [batch]
            exporters: [googlecloud/info]
        metrics/error:
            receivers: [logcount/error]
            processors: [batch]
            exporters: [googlecloud/error]
```
