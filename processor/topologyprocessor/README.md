# Topology Processor
This processor utilizes request headers to provide extended topology functionality in BindPlane.

## Minimum agent versions
- Introduced: [v1.6.7](https://github.com/observIQ/bindplane-agent/releases/tag/v1.6.7)

## Supported pipelines:
- Logs
- Metrics
- Traces

## Configuration
| Field                | Type      | Default | Description                                                                                                                                                               |
|----------------------|-----------|---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `interval`           | duration  | `1m`    | The interval at which topology data is sent to Bindplane via OpAMP.                                                                                                       |
| `organizationID`     | string    |         | The Organization ID of the Bindplane configuration where this processor is running.                                                                                       |
| `accountID`          | string    |         | The Account ID of the Bindplane configuration where this processor is running.                                                                                            |
| `configuration`      | string    |         | The name of the Bindplane configuration this processor is running on.                                                                                                     |


### Example configuration

The example configuration below shows ingesting logs and sampling the size of 50% of the OTLP log objects.

```yaml
receivers:
  filelog:
    inclucde: ["/var/log/*.log"]

processors:
  topology:
    interval: 1m
    organizationID: "myOrganizationID"
    accountID: "myAccountID"
    configuration: "myConfiguration"


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
