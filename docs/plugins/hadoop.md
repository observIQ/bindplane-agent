# `hadoop` plugin

The `hadoop` plugin consumes [Apache Hadoop](https://hadoop.apache.org/) logs entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- |--- | --- |
| `enable_datanode_logs` | `bool` | `true` | Enable collection of Hadoop data node logs |
| `datanode_log_path` | `string` | `"/usr/local/hadoop/logs/hadoop-*-datanode-*.log"`  | The absolute path to the data node logs |
| `enable_resourcemgr_logs` | `bool` | `true` | Enable the collection of ResourceManager logs |
| `resourcemgr_log_path` | `string` | `"/usr/local/hadoop/logs/hadoop-*-resourcemgr-*.log"`  | The absolute path to the ResourceManager logs |
| `enable_namenode_logs` | `bool` | `true` | Enable collection of Hadoop NameNode logs |
| `namenode_log_path` | `string` | `"/usr/local/hadoop/logs/hadoop-*-namenode-*.log"`  | The absolute path to the NameNode logs |
| `start_at` | `enum` | `end` | Start reading file from 'beginning' or 'end' | 

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: hadoop
- type: stdout

```

With non-standard log path for data node logs:

```yaml
pipeline:
- type: hadoop
  datanode_log_path: "/path/to/logs"
- type: stdout

```
