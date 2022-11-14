# MongoDB Plugin

Log parser for MongoDB

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| log_paths | The path of the log files | []string | `[/var/log/mongodb/mongodb.log* /var/log/mongodb/mongod.log*]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| retain_raw_logs | When enabled will preserve the original log message on the body in a `raw_log` key | bool | `false` | false |  |
| offset_storage_dir | The directory that the offset storage file will be created | string | `$OIQ_OTEL_COLLECTOR_HOME/storage` | false |  |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/mongodb_logs.yaml
    parameters:
      log_paths: [/var/log/mongodb/mongodb.log* /var/log/mongodb/mongod.log*]
      start_at: end
      retain_raw_logs: false
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
