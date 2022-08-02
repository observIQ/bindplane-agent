# `hbase` plugin

The `hbase` plugin consumes [Apache HBase](https://hbase.apache.org/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `enable_master_log` | `true` | Enable to collect HBase master logs |
| `master_log_path` | `"/usr/local/hbase*/logs/hbase*-master-*.log"`  | The absolute path to the HBase master logs |
| `enable_region_log` | `true` | Enable to collect HBase region logs |
| `region_log_path` | `"/usr/local/hbase*/logs/hbase*-regionserver-*.log"` | The absolute path to the HBase region logs |
| `enable_zookeeper_log` | `false` | Enable to collect HBase zookeeper logs |
| `zookeeper_log_path` | `"/usr/local/hbase*/logs/hbase*-zookeeper-*.log"` | The absolute path to the HBase zookeeper logs |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: hbase
- type: stdout

```

With non-standard log paths:

```yaml
pipeline:
- type: hbase
  master_log_path: "/path/to/logs"
- type: stdout

```
