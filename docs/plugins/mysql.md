# `mysql` plugin

The `mysql` plugin consumes [mysql](https://www.mysql.com/) log entries from the local filesystem and outputs parsed entries.

## Supported Platforms

- Linux
- Windows
- MacOS
- Kubernetes

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `enable_general_log` | `true` | Enable to collect MySQL general logs |
| `general_log_path` | `"/var/log/mysql/general.log"` | Path to general log file |
| `enable_slow_log` | `true`  | Enable to collect MySQL slow query logs |
| `slow_query_log_path` | `"/var/log/mysql/slow.log"` | Path to slow query log file |
| `enable_error_log` | `true` | Enable to collect MySQL error logs  |
| `error_log_path` | `"/var/log/mysql/mysqld.log"` | Path to mysqld log file |
| `enable_mariadb_audit_log` | `false` | Enable to collect MySQL audit logs provided by MariaDB Audit plugin |
| `mariadb_audit_log_path` | `"/var/log/mysql/audit.log"` | Path to audit log file created by MariaDB plugin |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Prerequisites

To enable certain log files, it may be necessary to edit the MySQL configuration file: mysqld.cnf.

### Audit Log
To enable MariaDB Audit logs, it is necessary to install the MariaDB Audit Plugin. More information can be found <a href="https://mariadb.com/kb/en/mariadb-audit-plugin-installation/" target="_blank">here</a>.

### Error Log
This is generally enabled by default. To change the file path, you can set or update "log_error" within mysqld.cnf.

For more details, see the error log documentation [here](https://dev.mysql.com/doc/refman/5.7/en/error-log.html ).

### Query Log
To enable the query log, set *general_log_file* to the desired log path and set *general_log = 1*.

For more details, see the query log documentation [here](https://dev.mysql.com/doc/refman/5.7/en/query-log.html).

### Slow Query Log
To enable the slow query log, set *slow_query_log_file* to the desired log path. Set *slow_query_log = 1* and optionally, configure *long_query_time*. 

For more details, see the slow query log documentation [here](https://dev.mysql.com/doc/refman/5.7/en/slow-query-log.html).

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: mysql
- type: stdout

```

With MariaDB Audit logs enabled and non-default log paths:

```yaml
pipeline:
- type: mysql
  enable_general_log: true
  general_log_path: "path/to/logs"
  enable_slow_log: true
  slow_query_log_path: "path/to/logs"
  enable_error_log: true
  error_log_path: "path/to/logs"
  enable_mariadb_audit_log: true
  mariadb_audit_log_path: "path/to/logs"
- type: stdout

```
