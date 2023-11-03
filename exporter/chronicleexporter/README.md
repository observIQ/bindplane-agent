# Chronicle Exporter

This exporter facilitates the sending of logs to Chronicle, which is a security analytics platform provided by Google. It is designed to integrate with OpenTelemetry collectors to export telemetry data such as logs to a Chronicle account.

## Minimum Collector Versions

- The minimum version of the OpenTelemetry Collector required for this exporter is not specified in the provided documentation.

## Supported Pipelines

- Logs

## How It Works

1. The exporter uses the configured credentials to authenticate with the Google Cloud services.
2. It marshals logs into the format expected by Chronicle.
3. It sends the logs to the appropriate regional Chronicle endpoint.

## Configuration

The exporter can be configured using the following fields:

| Field              | Type   | Default | Required | Description                                                                       |
| ------------------ | ------ | ------- | -------- | --------------------------------------------------------------------------------- |
| `region`           | string |         | `true`   | The region where the data will be sent, it must be one of the predefined regions. |
| `creds_file_path`  | string |         | `true`   | The file path to the Google credentials JSON file.                                |
| `creds`            | string |         | `true`   | The Google credentials JSON.                                                      |
| `log_type`         | string |         | `true`   | The type of log that will be sent.                                                |
| `raw_log_field`    | string |         | `false`  | The field name for raw logs.                                                      |
| `customer_id`      | string |         | `false`  | The customer ID used for sending logs.                                            |
| `sending_queue`    | struct |         | `false`  | Configuration for the sending queue.                                              |
| `retry_on_failure` | struct |         | `false`  | Configuration for retry logic on failure.                                         |
| `timeout_settings` | struct |         | `false`  | Configuration for timeout settings.                                               |

### Regions

Predefined regions include multiple global locations such as `Europe Multi-Region`, `Frankfurt`, `London`, `Singapore`, `Sydney`, `Tel Aviv`, `United States Multi-Region`, and `Zurich`. Each region has a specific endpoint URL.

## Example Configuration

```yaml
chronicle:
  region: "Europe Multi-Region"
  creds_file_path: "/path/to/google/creds.json"
  log_type: "threat_detection"
  customer_id: "customer-123"
```
