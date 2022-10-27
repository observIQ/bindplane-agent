# Redis Plugin

Log parser for Redis

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_path | The absolute path to the Redis logs | []string | `[/var/log/redis/redis-server.log /var/log/redis_6379.log /var/log/redis/redis.log /var/log/redis/default.log /var/log/redis/redis_6379.log]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| offset_storage_dir | The directory that the offset storage file will be created | string | `$OIQ_OTEL_COLLECTOR_HOME/storage` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/redis_logs.yaml
    parameters:
      file_path: [/var/log/redis/redis-server.log /var/log/redis_6379.log /var/log/redis/redis.log /var/log/redis/default.log /var/log/redis/redis_6379.log]
      start_at: end
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
