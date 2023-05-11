# Cockroach Database Metrics Plugin

Metrics receiver for Cockroach Database

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| endpoint | Endpoint used for HTTP requests from the DB Console | string | `localhost:8080` | false |  |
| username | Username to access sql Database (only needed if database is secure) | string | `` | false |  |
| password | Password to access sql Database (only needed if database is secure) | string | `` | false |  |
| scrape_interval | Time in between every scrape request | string | `60s` | false |  |
| timezone | Timezone to use when parsing the timestamp | timezone | `UTC` | false |  |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/cockroachdb_metrics.yaml
    parameters:
      endpoint: localhost:8080
      scrape_interval: 60s
      timezone: UTC
```
