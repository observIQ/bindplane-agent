# `oracle_database` plugin

The `oracle_database` plugin consumes [OracleDB](https://www.oracle.com/database/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- |--- | --- |
| `enable_audit_log` | `bool` | `true` |  Enable to collect OracleDB audit logs. |
| `audit_log_path` | `string` | `"/u01/app/oracle/product/*/dbhome_1/admin/*/adump/*.aud"` | Path to the audit log file. | 
| `enable_truncate_audit_action` | `bool` | `true` | Whether or not to truncate the audit log action field. |
| `enable_alert_log` | `bool` | `true` | Enable to collect OracleDB alert logs. |
| `alert_log_path` | `string` | `"/u01/app/oracle/product/*/dbhome_1/diag/rdbms/*/*/trace/alert_*.log"` | Path to the alert log file. | 
| `enable_listener_log` | `bool` | `true` |  Enable to collect OracleDB listener logs. |
| `listener_log_path` | `string` | `"/u01/app/oracle/product/*/dbhome_1/diag/tnslsnr/*/listener/alert/log.xml"` | Path to the listener log file. | 
| `start_at` | `enum` | `end` | Start reading file from 'beginning' or 'end' | 

## Prerequisites

### Audit Log Location

You can find the Audit log location by running the following query while connected to the database.

```sql
select
   name 
from
   v$parameter
where
   name = 'audit_file_dest';
```

### Alert Log location
You can find the Alert log location by running the following query while connected to the database.
```sql
select
   name 
from
   v$parameter
where
   name = 'background_dump_destination';
```
### Listener Log Location
You can find the Listener log location by running the following command on the database host.
```shell
$ lsnrctl status
```

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: oracle_database
- type: stdout

```

With audit logs disabled:

```yaml
pipeline:
- type: oracle_database
  enable_audit_log: false
- type: stdout

```
