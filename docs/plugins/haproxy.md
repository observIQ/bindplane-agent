# `haproxy` plugin

The `haproxy` plugin consumes [HAProxy](http://www.haproxy.org/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `file_log_path` |  | Specify a single path or multiple paths to read one or many files. You may also use a wildcard (*) to read multiple files within a directory. This field is required. |
| `log_format` | `default`  | When choosing the 'default' option, the agent will expect and parse logs in a format of HTTP or TCP as well as any log entries that matches the default or error logging configuration. HAProxy uses default logging format when no specific option is set. When choosing the 'observIQ' option, the agent will expect and parse logs in an optimized JSON format that adheres to the observIQ specification, requiring an update to the Log-Format for each mode. See the HAProxy source page for more information. |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using default log format:

```yaml
pipeline:
- type: haproxy
  file_log_path:
    - "/path/to/logs"
- type: stdout

```

With `observiq` log format:

```yaml
pipeline:
- type: haproxy
  file_log_path:
    - "/path/to/logs"
  log_format: observiq
- type: stdout

```
