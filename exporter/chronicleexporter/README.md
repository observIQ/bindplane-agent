# Chronicle Exporter

**Currently only v2 of the ingestion API is supported**

This exporter facilitates the sending of logs to Chronicle, which is a security analytics platform provided by Google. It is designed to integrate with OpenTelemetry collectors to export telemetry data such as logs to a Chronicle account.

## Minimum Collector Versions

- Introduced: [v1.39.0](https://github.com/observIQ/bindplane-agent/releases/tag/v1.39.0)

## Supported Pipelines

- Logs

## How It Works

1. The exporter uses the configured credentials to authenticate with the Google Cloud services.
2. It marshals logs into the format expected by Chronicle.
3. It sends the logs to the appropriate Chronicle endpoint.

## Configuration

The exporter can be configured using the following fields:

| Field                           | Type              | Default                                | Required | Description                                                                                 |
| ------------------------------- | ----------------- | -------------------------------------- | -------- | ------------------------------------------------------------------------------------------- |
| `endpoint`                      | string            | `malachiteingestion-pa.googleapis.com` | `false`  | The Endpoint for sending to chronicle.                                                      |
| `creds_file_path`               | string            |                                        | `true`   | The file path to the Google credentials JSON file.                                          |
| `creds`                         | string            |                                        | `true`   | The Google credentials JSON.                                                                |
| `log_type`                      | string            |                                        | `false`  | The type of log that will be sent.                                                          |
| `raw_log_field`                 | string            |                                        | `false`  | The field name for raw logs.                                                                |
| `customer_id`                   | string            |                                        | `false`  | The customer ID used for sending logs.                                                      |
| `override_log_type`             | bool              | `true`                                 | `false`  | Whether or not to override the `log_type` in the config with `attributes["log_type"]`       |
| `namespace`                     | string            |                                        | `false`  | User-configured environment namespace to identify the data domain the logs originated from. |
| `compression`                   | string            | `none`                                 | `false`  | The compression type to use when sending logs. valid values are `none` and `gzip`           |
| `ingestion_labels`              | map[string]string |                                        | `false`  | Key-value pairs of labels to be applied to the logs when sent to chronicle.                 |
| `collect_agent_metrics`         | bool              | `true`                                 | `false`  | Enables collecting metrics about the agent's process and log ingestion metrics              |
| `batch_log_count_limit_grpc`    | int               | `1000`                                 | `false`  | The maximum number of logs allowed in a gRPC batch creation request.                        |
| `batch_request_size_limit_grpc` | int               | `1048576`                              | `false`  | The maximum size, in bytes, allowed for a gRPC batch creation request.                      |
| `batch_log_count_limit_http`    | int               | `1000`                                 | `false`  | The maximum number of logs allowed in a HTTP batch creation request.                        |
| `batch_request_size_limit_http` | int               | `1048576`                              | `false`  | The maximum size, in bytes, allowed for a HTTP batch creation request.                      |

### Log Type

If the `attributes["log_type"]` field is present in the log, and maps to a known Chronicle `log_type` the exporter will use the value of that field as the log type. If the `attributes["log_type"]` field is not present, the exporter will use the value of the `log_type` configuration field as the log type.

currently supported log types are:

- windows_event.security
- windows_event.custom
- windows_event.application
- windows_event.system
- sql_server

If the `attributes["chronicle_log_type"]` field is present in the log, we will use its value in the payload instead of the automatic detection or the `log_type` in the config.

### Namespace and Ingestion Labels

If the `attributes["chronicle_namespace"]` field is present in the log, we will use its value in the payload instead of the `namespace` in the config.

If there are nested fields in `attributes["chronicle_ingestion_label"]`, we will use the values in the payload instead of the `ingestion_labels` in the config.

## Credentials

This exporter requires a Google Cloud service account with access to the Chronicle API. The service account must have access to the endpoint specfied in the config.
Besides the default endpoint, there are also regional endpoints that can be used [here](https://cloud.google.com/chronicle/docs/reference/ingestion-api#regional_endpoints).

For additional information on accessing Chronicle, see the [Chronicle documentation](https://cloud.google.com/chronicle/docs/reference/ingestion-api#getting_api_authentication_credentials).

## Log Batch Creation Request Limits

`batch_log_count_limit_grpc`, `batch_request_size_limit_grpc`, `batch_log_count_limit_http`, `batch_request_size_limit_http` are all used for ensuring log batch creation requests don't exceed Chronicle's backend limits - the former two for Chronicle's gRPC endpoint, and the latter two for Chronicle's HTTP endpoint. If a request either exceeds the configured size limit or contains more logs than the configured log count limit, the request will be split into multiple requests that adhere to these limits, with each request containing a subset of the logs contained in the original request. Any single logs that result in the request exceeding the size limit will be dropped.

## Example Configuration

### Basic Configuration

```yaml
chronicle:
  creds_file_path: "/path/to/google/creds.json"
  log_type: "ABSOLUTE"
  customer_id: "customer-123"
```

### Basic Configuration with Regional Endpoint

```yaml
chronicle:
  endpoint: https://malachiteingestion-pa.googleapis.com
  creds_file_path: "/path/to/google/creds.json"
  log_type: "ONEPASSWORD"
  customer_id: "customer-123"
```

### Configuration with Ingestion Labels

```yaml
chronicle:
  creds_file_path: "/path/to/google/creds.json"
  log_type: ""
  customer_id: "customer-123"
  ingestion_labels:
    env: dev
    zone: USA
```
