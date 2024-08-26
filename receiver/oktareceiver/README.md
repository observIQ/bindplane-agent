# Okta Receiver
This receiver is capable of collecting logs from an Okta domain.

## Minimum Agent Versions
- Introduced: [v1.59.0](https://github.com/observIQ/bindplane-agent/releases/tag/v1.59.0)

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
| domain               |  string   |                  | `true`   | The Okta domain the receiver should collect logs from.
| api_token            |  string   |                  | `true`   | An Okta API Token generated from the above Okta domain.
| poll_interval        |  string   | 1m               | `false`  | The rate at which this receiver will poll Okta for logs. This value must be at least 1 second and must be a string readable by Golang's [time.ParseDuration](https://pkg.go.dev/time#ParseDuration).
| start_time           |  string   | Now (UTC)        | `false`  | The UTC timestamp indicating the beginning of the range of logs this receiver will collect. Must be within the past 180 days and not in the future. Must be in the format "yyyy-mm-ddThh-mm-ssZ"

### Example Configuration
```yaml
receivers:
  okta:
    domain: example.okta.com
    api_token: 11Z-XDEwgRIf4p0-RqbSFoplFh_84EOtC_ka4J7ylx
    poll_interval: 2m
    start_time: "2024-08-12T00:00:00Z"
exporters:
  googlecloud:
    project: my-gcp-project

service:
  pipelines:
    logs:
      receivers: [okta]
      exporters: [googlecloud]
```
