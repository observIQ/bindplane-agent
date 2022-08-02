# `couchbase` plugin

The `couchbase` plugin consumes [Couchbase](https://www.couchbase.com/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- |--- | --- |
| `enable_error_log` | `bool` | `true` | The absolute path to the Couchbase error logs messages |
| `error_log_path` | `string` | `"/opt/couchbase/var/lib/couchbase/logs/error.log"` | The absolute path to the Couchbase error logs |
| `enable_debug_log` | `bool` | `false`  | Enable to collect Couchbase debug logs |
| `debug_log_path` | `string` | `"/opt/couchbase/var/lib/couchbase/logs/debug.log"` | The absolute path to the Couchbase debug logs | 
| `enable_info_log` | `bool` | `false` | Enable to collect Couchbase info logs | 
| `info_log_path` | `string` | `"/opt/couchbase/var/lib/couchbase/logs/info.log"` | The absolute path to the Couchbase information logs | 
| `enable_access_log` | `bool` | `true` | Enable to collect Couchbase http access logs | 
| `http_access_log_path` | `string` | `"/opt/couchbase/var/lib/couchbase/logs/http_access.log"` | The absolute path to the Couchbase http access logs | 
| `enable_internal_access_log` | `bool` | `false` | Enable to collect Couchbase internal http access logs | 
| `http_internal_access_log_path` | `string` | `"/opt/couchbase/var/lib/couchbase/logs/http_access_internal.log"` | The absolute path to the Couchbase internal http access logs | 
| `enable_babysitter_log` | `bool` | `true` | Enable to collect Couchbase babysitter logs | 
| `babysitter_log_path` | `string` | `"/opt/couchbase/var/lib/couchbase/logs/babysitter.log"` | The absolute path to the Couchbase babysitter logs | 
| `enable_xdcr_log` | `bool` | `false` | Enable to collect Couchbase Cross Datacenter Replication logs | 
| `xdcr_log_path` | `string` | `"/opt/couchbase/var/lib/couchbase/logs/goxdcr.log"` | The absolute path to the Couchbase cross datacenter replication logs | 
| `start_at` | `enum` | `end` | Start reading file from 'beginning' or 'end' | 

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: couchbase
- type: stdout

```

With internal http access logs enabled and a non-standard error log path:

```yaml
pipeline:
- type: couchbase
  enable_internal_access_log: true
  error_log_path: "/path/to/logs"
- type: stdout

```
