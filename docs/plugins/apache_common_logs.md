# Apache Common Plugin

Log parser for Apache common format

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_path | Path to apache common formatted log file | []string | `[/var/log/apache2/access.log]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/apache_common_logs.yaml
    parameters:
      file_path: [/var/log/apache2/access.log]
      start_at: end
```
