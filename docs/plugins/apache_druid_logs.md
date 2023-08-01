# Apache Druid Plugin

Log parser for Apache Druid

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| enable_broker_logs | Enable to collect Apache Druid broker logs | bool | `true` | false |  |
| broker_log_path | Absolute filepath containing Apache druid broker logs. Required if enable_broker_logs is true | string |  | false |  |
| enable_coordinator_overlord_logs | Enable to collect Apache Druid coordinator and overlord logs | bool | `true` | false |  |
| coordinator_overlord_log_path | Absolute filepath containing Apache Druid coordinator and overlord logs. Required if enable_coordinator_overlord_logs is true | string |  | false |  |
| enable_historical_logs | Enable to collect Apache Druid historical logs | bool | `true` | false |  |
| historical_log_path | Absolute filepath containing Apache Druid historical logs. Required if enable_historical_logs is true | string |  | false |  |
| enable_middle_manager_logs | Enable to collect Apache Druid middle manager logs | bool | `true` | false |  |
| middle_manager_log_path | Absolute filepath containing Apache Druid middle manager logs. Required if enable_middle_manager_logs is true | string |  | false |  |
| enable_router_logs | Enable to collect Apache Druid router logs | bool | `true` | false |  |
| router_log_path | Absolute filepath containing Apache Druid router logs. Required if enable_router_logs is true | string |  | false |  |
| enable_zookeeper_logs | Enable to collect Apache Druid ZooKeeper logs | bool | `true` | false |  |
| zookeeper_log_path | Absolute filepath containing Apache Druid ZooKeeper logs. Required if enable_zookeeper_logs is true | string |  | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false |  |
| timezone | Timezone to use when parsing the timestamp, since the log timestamp does not contain a timezone. | timezone | `UTC` | false |  |
| retain_raw_logs | When enabled will preserve the original log message in a `raw_log` key. This will either be in the `body` or `attributes` depending on how `parse_to` is configured. | bool | `false` | false |  |
| parse_to | Where to parse structured log parts | string | `body` | false | `body`, `attributes` |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/apache_druid_logs.yaml
    parameters:
      enable_broker_logs: true
      enable_coordinator_overlord_logs: true
      enable_historical_logs: true
      enable_middle_manager_logs: true
      enable_router_logs: true
      enable_zookeeper_logs: true
      start_at: end
      timezone: UTC
      retain_raw_logs: false
      parse_to: body
```
