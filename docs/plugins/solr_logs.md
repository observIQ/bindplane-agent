# Solr Plugin

Log parser for Solr

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_log_path | The absolute path to the Solr logs | []string | `[/var/solr/logs/solr.log]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| offset_storage_dir | The directory that the offset storage file will be created | string | `$OIQ_OTEL_COLLECTOR_HOME/storage` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/solr_logs.yaml
    parameters:
      file_log_path: [/var/solr/logs/solr.log]
      start_at: end
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
