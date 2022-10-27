# OracleDB Metrics Plugin

Metrics receiver for OracleDB

## Supported Versions

- 12.2
- 18c
- 19c
- 21c

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| endpoint | Address to scrape metrics from | string | `localhost:1521` | false |  |
| username | Database user to run metric queries with | string |  | true |  |
| password | Password for user | string |  | false |  |
| sid | Site Identifier. One or both of sid or service_name must be specified. | string |  | false |  |
| service_name | OracleDB Service Name. One or both of sid or service_name must be specified. | string |  | false |  |
| wallet | OracleDB Wallet file location (must be URL encoded) | string |  | false |  |
| scrape_interval | Time in between every scrape request | string | `60s` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/oracledb_metrics.yaml
    parameters:
      endpoint: localhost:1521
      username: $USERNAME
      scrape_interval: 60s 
      sid: $SID
```

## Least Privileged User

To collect metrics, a monitoring user requires `SELECT` access to the relevant views. The following sql
script should create a monitoring user and give it the appropriate permissions if executed by
a user with sufficient permissions connected to the Oracle DB instance as SYSDBA or SYSOPER.

```sql
-- Create the monitoring user "otel"
CREATE USER otel IDENTIFIED BY <authentication password>;

-- Grant the "otel" user the required permissions
GRANT CONNECT TO otel;
GRANT SELECT ON SYS.GV_$DATABASE to otel;
GRANT SELECT ON SYS.GV_$INSTANCE to otel;
GRANT SELECT ON SYS.GV_$PROCESS to otel;
GRANT SELECT ON SYS.GV_$RESOURCE_LIMIT to otel;
GRANT SELECT ON SYS.GV_$SYSMETRIC to otel;
GRANT SELECT ON SYS.GV_$SYSSTAT to otel;
GRANT SELECT ON SYS.GV_$SYSTEM_EVENT to otel;
GRANT SELECT ON SYS.V_$RMAN_BACKUP_JOB_DETAILS to otel;
GRANT SELECT ON SYS.V_$SORT_SEGMENT to otel;
GRANT SELECT ON SYS.V_$TABLESPACE to otel;
GRANT SELECT ON SYS.V_$TEMPFILE to otel;
GRANT SELECT ON SYS.DBA_DATA_FILES to otel;
GRANT SELECT ON SYS.DBA_FREE_SPACE to otel;
GRANT SELECT ON SYS.DBA_TABLESPACE_USAGE_METRICS to otel;
GRANT SELECT ON SYS.DBA_TABLESPACES to otel;
GRANT SELECT ON SYS.GLOBAL_NAME to otel;
```
