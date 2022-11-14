# Aerospike Plugin

Log parser for Aerospike

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| journald_directory | The absolute path to the general Aerospike logs | string |  | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/aerospike_logs.yaml
    parameters:
      start_at: end
```
