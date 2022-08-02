# `jboss` plugin

The `jboss` plugin consumes [JBoss EAP](https://hadoop.apache.org/) logs entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- |--- | --- |
| `file_path` | `string` | `"/usr/local/JBoss/EAP-*/*/log/server.log"` | The absolute path to the JBoss logs |
| `start_at` | `enum` | `end` | Start reading file from 'beginning' or 'end' | 

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: jboss
- type: stdout

```

With non-standard log path:

```yaml
pipeline:
- type: jboss
  file_path: "/path/to/logs"
- type: stdout

```
