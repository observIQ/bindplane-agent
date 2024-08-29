# Proofpoint Receiver
This receiver is capable of collecting logs from Proofpoint Targeted Attack Protection (TAP)

## Minimum Agent Versions
- Introduced: [v1.59.0](https://github.com/observIQ/bindplane-agent/releases/tag/v1.59.0)

## Supported Pipelines
- Logs

## How It Works
1. The user configures this receiver in a pipeline.
2. The user configures a supported component to route telemetry from this receiver.

## Prerequisites
- A Proofpoint API Token will be needed to authorize the receiver with your Proofpoint account.

## Configuration
| Field                | Type      | Default          | Required | Description                                                                                                                                                                             |
|----------------------|-----------|------------------|----------| ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| principal            |  string   |                  | `true`   |                                                                                                                      |
| secret               |  string   |                  | `true`   |                                                                                                                      |
| poll_interval        |  string   | 5m               | `false`  |                                                                                                                      |

### Example Configuration
```yaml
receivers:
  proofpoint:
    principal: example.okta.com
    secret: 11Z-XDEwgRIf4p0-RqbSFoplFh_84EOtC_ka4J7ylx
    poll_interval: 2m
exporters:
  googlecloud:
    project: my-gcp-project

service:
  pipelines:
    logs:
      receivers: [proofpoint]
      exporters: [googlecloud]
```
