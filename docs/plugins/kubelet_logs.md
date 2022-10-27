# Kubernetes Kubelet Plugin

Log parser for Kubelet journald logs

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| journald_directory | Directory containing journal files to read entries from | string |  | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| timezone | Timezone to use when parsing the timestamp | timezone | `UTC` | false |  |
| retain_raw_logs | When enabled will preserve the original log message on the body in a `raw_log` key | bool | `false` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/kubelet_logs.yaml
    parameters:
      start_at: end
      timezone: UTC
      retain_raw_logs: false
```
