# Apache Combined Plugin

Log parser for Apache combined format

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_path | Paths to Apache combined formatted log files | []string | `[/var/log/apache_combined.log]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| retain_raw_logs | When enabled will preserve the original log message on the body in a `raw_log` key | bool | `false` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/apache_combined_logs.yaml
    parameters:
      file_path: [/var/log/apache_combined.log]
      start_at: end
      retain_raw_logs: false
```
