# Microsoft SQL Server Plugin

Log Parser for Microsoft SQL Server Event Logs

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| poll_interval | The interval at which a channel is checked for new log entries | string | `1s` | false |  |
| max_reads | The maximum number of events read into memory at one time | int | `1000` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/sql_server_logs.yaml
    parameters:
      poll_interval: 1s
      max_reads: 1000
      start_at: end
```
