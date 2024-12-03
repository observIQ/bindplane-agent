# Topology Processor
This processor utilizes request headers to provide extended topology functionality in BindPlane.

## Minimum agent versions
- Introduced: [v1.6.6](https://github.com/observIQ/bindplane-agent/releases/tag/v1.6.6)

## Supported pipelines:
- Logs
- Metrics
- Traces

## Configuration
| Field               | Type      | Default | Description                                                                                                                                                               |
|---------------------|-----------|---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `enabled`           | bool      | `false` | When `true`, this processor will look for incoming topology headers and track the relevant connections accordingly. If false this processor acts as a no-op.              |
| `interval`          | duration  | `1m`    | The interval at which topology data is sent to Bindplane via OpAMP.                                                                                                       |
| `configName`        | string    |         | The name of the Bindplane configuration this processor is running on.                                                                                                     |
| `orgID`             | string    |         | The Organization ID of the Bindplane configuration where this processor is running.                                                                                       |
| `accountID`         | string    |         | The Account ID of the Bindplane configuration where this processor is running.                                                                                            |


### Example configuration

The example configuration below shows ingesting logs and sampling the size of 50% of the OTLP log objects.

```yaml
receivers:
  filelog:
    inclucde: ["/var/log/*.log"]

processors:
  topology:
    enabled: true
    interval: 1m
    configName: "myConfigName"
    orgID: "myOrgID"
    accountID: "myAccountID"

exporters:
  googlecloud:

service:
  pipelines:
    logs:
      receivers:
        - filelog
      processors:
        - topology
      exporters:
        - googlecloud
```
