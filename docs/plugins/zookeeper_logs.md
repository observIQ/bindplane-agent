# Apache Zookeeper Plugin

Log parser for Apache Zookeeper

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_path | The absolute path to the Zookeeper logs | []string | `[/home/kafka/kafka/logs/zookeeper.log]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| offset_storage_dir | The directory that the offset storage file will be created | string | `$OIQ_OTEL_COLLECTOR_HOME/storage` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/zookeeper_logs.yaml
    parameters:
      file_path: [/home/kafka/kafka/logs/zookeeper.log]
      start_at: end
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
