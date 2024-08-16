# Okta Receiver
This receiver is capable of collecting logs from an Okta domain.

## Minimum Agent Versions
- TODO

## Supported Pipelines
- Logs

## How It Works
1. The user configures this receiver in a pipeline.
2. The user configures a supported component to route telemetry from this receiver.

## Prerequisites
- An Okta API Token will be needed to authorize the receiver with your Okta Domain.

## Configuration
| Field                | Type      | Default          | Required | Description                                                                                                                                                                            |
|----------------------|-----------|------------------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| domain             |  string   |                  | `true`   | The Okta domain the receiver should collect logs from.                                                                                         |
| api_token                 |  string   |                  | `true`  | An Okta API Token generated from the above Okta domain.

### Example Configuration
```yaml
receivers:
  http:
    domain: example.okta.com
    api_token: 11Z-XDEwgRIf4p0-RqbSFoplFh_84EOtC_ka4J7ylx
exporters:
  googlecloud:
    project: my-gcp-project

service:
  pipelines:
    logs:
      receivers: [okta]
      exporters: [googlecloud]
```
