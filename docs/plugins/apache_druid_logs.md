# Apache Druid Plugin

Log parser for Apache Druid

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| log_directory | Absolute filepath of the folder which contains the Apache Druid logs | string |  | true |  |
| enable_broker_logs | Enable to collect Apache Druid broker logs | bool | `true` | false |  |
| broker_log_filename | Filename containing Apache druid broker logs | string | `broker.log` | false |  |
| enable_coordinator_overlord_logs | Enable to collect Apache Druid coordinator and overlord logs | bool | `true` | false |  |
| coordinator_overlord_log_filename | Filename containing Apache Druid coordinator and overlord logs | string | `coordinator-overlord.log` | false |  |
| enable_historical_logs | Enable to collect Apache Druid historical logs | bool | `true` | false |  |
| historical_log_filename | Filename containing Apache Druid historical logs | string | `historical.log` | false |  |
| enable_middle_manager_logs | Enable to collect Apache Druid middle manager logs | bool | `true` | false |  |
| middle_manager_log_filename | Filename containing Apache Druid middle manager logs | string | `middleManager.log` | false |  |
| enable_router_logs | Enable to collect Apache Druid router logs | bool | `true` | false |  |
| router_log_filename | Filename containing Apache Druid router logs | string | `router.log` | false |  |
| enable_zookeeper_logs | Enable to collect Apache Druid ZooKeeper logs | bool | `true` | false |  |
| zookeeper_log_filename | Filename containing Apache Druid ZooKeeper logs | string | `zookeeper.log` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false |  |
| timezone | Timezone to use when parsing the timestamp | timezone | `UTC` | false |  |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/apache_druid_logs.yaml
    parameters:
      log_directory: $LOG_DIRECTORY
      enable_broker_logs: true
      broker_log_filename: broker.log
      enable_coordinator_overlord_logs: true
      coordinator_overlord_log_filename: coordinator-overlord.log
      enable_historical_logs: true
      historical_log_filename: historical.log
      enable_middle_manager_logs: true
      middle_manager_log_filename: middleManager.log
      enable_router_logs: true
      router_log_filename: router.log
      enable_zookeeper_logs: true
      zookeeper_log_filename: zookeeper.log
      start_at: end
      timezone: UTC
```
