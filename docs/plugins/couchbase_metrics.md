# Couchbase Metrics Plugin

Metrics receiver for Couchbase

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| endpoint | Address to scrape metrics from | string | `localhost:8091` | false |  |
| username | Username to use as header of every scrape request | string |  | true |  |
| password | Password to use as header of every scrape request | string |  | true |  |
| scrape_interval | Time in between every scrape request | string | `60s` | false |  |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/couchbase_metrics.yaml
    parameters:
      endpoint: localhost:8091
      username: $USERNAME
      password: $PASSWORD
      scrape_interval: 60s
```
