# `kafka` plugin

The `kafka` plugin consumes [Kafka](https://kafka.apache.org/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `enable_server_log` | `true` | Enable to collect Apache Kafka server logs |
| `server_log_path` | `"/home/kafka/kafka/logs/server.log*"`  | Apache Kafka server log path |
| `enable_controller_log` | `true` | Enable to collect Apache Kafka controller logs |
| `controller_log_path` | `"/home/kafka/kafka/logs/controller.log*"` | Apache Kafka controller log path |
| `enable_state_change_log` | `true` | Enable to collect Apache Kafka state change logs |
| `state_change_log_path` | `"/home/kafka/kafka/logs/state-change.log*"` | Apache Kafka state-change log path |
| `enable_log_cleaner_log` | `true` | Enable to collect Apache Kafka log cleaner logs |
| `log_cleaner_log_path` | `"/home/kafka/kafka/logs/log-cleaner.log*"` | Apache Kafka log-cleaner log path |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: kafka
- type: stdout

```

With non-standard log paths:

```yaml
pipeline:
- type: kafka
  server_log_path: "/path/to/logs"
- type: stdout

```
