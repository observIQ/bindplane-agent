# Router Processor
This processor is used to route logs to different pipelines based on the log's structure.

## Supported pipelines
- Logs

## How It Works
1. The user configures the router processor in their logs pipeline and one or more route receivers in other log pipelines.
2. If any incoming logs match a routes `match` expression, the router processor will send the matching logs to the corresponding route receiver.
3. If a log does not match any route then it will be sent along to the next component after the router processor in the main pipeline.

## Configuration
| Field | Type | Default | Description |
| --- | --- | --- | --- |
| routes | [][Route](#route) | [] | A list of routes to match incoming logs against. **Required field** |

### Route
| Field | Type | Default | Description |
| --- | --- | --- | --- |
| match  | string   | `""`  | A boolean [expression](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md) used to match which logs to route. By default, all logs are counted. **Required field** |
| route        | string   | ` `      | The name of the [route receiver](../../receiver/routereceiver/README.md) to send metrics to. |


### Example Config
The following config is an example of routing a single log file into several different log files based on the log's severity. If any logs have a `severity_enum` other than those in the routes it'll pass through to the `file/other` exporter.

```yaml
receivers:
    filelog:
        include: [ /var/log/myapp.log ]
    routereceiver/debug:
    routereceiver/info:
    routereceiver/warn:
    routereceiver/error:

processors:
  router:
    routes:
      - match: 'severity_enum == "debug"'
        route: debug
      - match: 'severity_enum == "info"'
        route: info
      - match: 'severity_enum == "warn"'
        route: warn
      - match: 'severity_enum == "error"'
        route: error

exporters:
    file/debug:
        path: /var/log/myapp_debug.log
    file/info:
        path: /var/log/myapp_info.log
    file/warn:
        path: /var/log/myapp_warn.log
    file/error:
        path: /var/log/myapp_error.log
    file/other:
        path: /var/log/myapp_other.log

service:
  pipelines:
    logs/main:
        receivers: [filelog]
        processors: [router]
        exporters: [file/other]
    logs/debug:
        receivers: [routereceiver/debug]
        exporters: [file/debug]
    logs/info:
        receivers: [routereceiver/info]
        exporters: [file/info]
    logs/warn:
        receivers: [routereceiver/warn]
        exporters: [file/warn]
    logs/error:
        receivers: [routereceiver/error]
        exporters: [file/error]
```

## Expression Language
In order to match or extract values from logs, the following `keys` are reserved and can be used to traverse the logs data model.

| Key               | Description |
| ---               | ---   |
| `body`            | Used to access the body of the log. |
| `attributes`      | Used to access the attributes of the log. |
| `resource`        | Used to access the resource of the log. |
| `severity_enum`   | Used to access the severity enum of the log. |
| `severity_number` | Used to access the severity number of the log. |

In order to access embedded values, use JSON dot notation. For example, `body.example.field` can be used to access a field two levels deep on the log body. 

However, if a key already possesses a literal dot, users will need to use bracket notation to access that field. For example, when the field `service.name` exists on the log's resource, users will need to use `resource["service.name"]` to access this value.

For more information about syntax and available operators, see the [Expression Language Definition](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md).

## How To
### Route on Severity
The following configuration adds a match expression that will route on severity.
```yaml
processors:
  router:
    routes:
      - match: 'severity_enum == "debug"'
        route: debug
      - match: 'severity_enum == "info"'
        route: info
      - match: 'severity_enum == "warn"'
        route: warn
      - match: 'severity_enum == "error"'
        route: error
```

### Route on attributes
The following configuration routes on an attribute `status` that as parsed out of the log.
```yaml
processors:
  router:
    routes:
      - match: attributes.status startsWith "2"'
        route: httpOK
      - match: attributes.status startsWith "3"'
        route: httpRedirect
      - match: attributes.status startsWith "4"'
        route: httpClientError
      - match: attributes.status startsWith "5"'
        route: httpServerError
```
