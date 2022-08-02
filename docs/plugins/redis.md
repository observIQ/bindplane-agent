# `redis` plugin

The `redis` plugin consumes [redis](https://redis.io/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `file_path` | `"/var/log/redis/redis.log*"` | The absolute path to the Redis logs |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using default log path:

```yaml
pipeline:
- type: redis
- type: stdout

```

With non-standard log path:

```yaml
pipeline:
- type: redis
  file_path: "/path/to/logs"
- type: stdout

```
