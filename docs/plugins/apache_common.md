# `apache_common` plugin

The `apache_common` plugin consumes [Apache Common](https://httpd.apache.org/docs/2.4/logs.html) log format entries from the local filesystem and outputs parsed entries.

## Supported Platforms

- Linux
- Windows
- MacOS
- Kubernetes

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `log_path` | `string` | `"/var/log/apache2/access.log"` | Path to apache common formatted log file |
| `start_at` | `enum` | `end` | Start reading file from 'beginning' or 'end' |

## Prerequisites

No prerequisite actions required.

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: apache_common
- type: stdout

```

With non-standard log paths:

```yaml
pipeline:
- type: apache_common
  log_path: "/path/to/logs"
- type: stdout

```
