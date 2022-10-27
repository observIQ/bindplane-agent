# HAProxy Plugin

Log parser for HAProxy

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_path | Specify a single path or multiple paths to read one or many files. You may also use a wildcard (*) to read multiple files within a directory | []string | `[/var/log/haproxy/haproxy.log]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| timezone | Timezone to use when parsing the timestamp | timezone | `UTC` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/haproxy_logs.yaml
    parameters:
      file_path: [/var/log/haproxy/haproxy.log]
      start_at: end
      timezone: UTC
```
