# `elasticsearch` plugin

The `elasticsearch` plugin consumes [Elasticsearch](https://www.elastic.co/elasticsearch/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `enable_general_logs` | `true` | Enable to collect Elasticsearch general logs |
| `general_log_path` | `"/var/log/elasticsearch/elasticsearch.log*"`  | The absolute path to the Elasticsearch general logs |
| `enable_deprecation_logs` | `true` | Enable to collect Elasticsearch deprecation logs |
| `deprecation_log_path` | `"/var/log/elasticsearch/elasticsearch_deprecation.log*"` | The absolute path to the Elasticsearch deprecation logs |
| `enable_gc_logs` | `false` | Enable to collect Elasticsearch garbage collection logs |
| `gc_log_path` | `"/var/log/elasticsearch/gc.log*"` | The absolute path to the Elasticsearch garbage collection logs |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: elasticsearch
- type: stdout

```

With non-standard log paths:

```yaml
pipeline:
- type: elasticsearch
  general_log_path: "/path/to/logs"
  deprecation_log_path: "/path/to/logs"
  enable_gc_logs: true
  gc_log_path: "/path/to/logs"
- type: stdout

```
