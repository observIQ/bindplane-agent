# `couchdb` plugin

The `couchdb` plugin consumes [Apache CouchDB](https://couchdb.apache.org/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `file_path` | `"/var/log/couchdb/couchdb.log"` | The absolute path to the CouchDB logs |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using default log path:

```yaml
pipeline:
- type: couchdb
- type: stdout

```

With non-standard log path:

```yaml
pipeline:
- type: couchdb
  log_path: "/path/to/logs"
- type: stdout

```
