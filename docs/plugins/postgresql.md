# `postgresql` plugin

The `postgresql` plugin consumes [PostgreSQL](https://www.postgresql.org/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `postgresql_log_path` | `"/var/lib/pgsql/*/data/pg_log/postgresql-*.log"` | Path to the PostgreSQL log file |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Prerequisites

Locate the local PostgreSQL postgresql.conf configuration file. This file is typically located in the database data_directory. For more information, see the [PostgreSQL documentation](https://www.postgresql.org/docs/9.1/runtime-config-file-locations.html).

Modify the postgresql.conf file with the following logging parameters. Changes should go under the ERROR REPORTING AND LOGGING section of the file. For more details on the log parameters, see the [PostgreSQL documentation](https://www.postgresql.org/docs/11/runtime-config-logging.html).

```text
log_destination = 'stderr'
logging_collector = on
log_directory = 'pg_log'        
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_truncate_on_rotation = off
log_rotation_age = 1d
log_min_duration_statement = 0
log_connections = on
log_duration = on
log_hostname = on
log_timezone = 'UTC'
log_line_prefix = 't=%t p=%p s=%c l=%l u=%u db=%d r=%r '
```

After editing the postgresql.conf file, you will need to restart the postgres server:

```shell
sudo service postgresql restart
```

## Example usage

### Configuration

Using default log path:

```yaml
pipeline:
- type: postgresql
- type: stdout

```

With non-standard log path:

```yaml
pipeline:
- type: postgresql
  postgresql_log_path: "/path/to/logs"
- type: stdout

```
