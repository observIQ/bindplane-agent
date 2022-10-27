# JBoss Plugin

Log parser for JBoss

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_path | The absolute path to the JBoss logs | []string | `[/usr/local/JBoss/EAP-*/*/log/server.log]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| timezone | Timezone to use when parsing the timestamp | timezone | `UTC` | false |  |
| offset_storage_dir | The directory that the offset storage file will be created | string | `$OIQ_OTEL_COLLECTOR_HOME/storage` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/jboss_logs.yaml
    parameters:
      file_path: [/usr/local/JBoss/EAP-*/*/log/server.log]
      start_at: end
      timezone: UTC
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
