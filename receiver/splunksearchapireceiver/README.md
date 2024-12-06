# Splunk Search API Receiver
This receiver collects Splunk events using the [Splunk Search API](https://docs.splunk.com/Documentation/Splunk/9.3.1/RESTREF/RESTsearch).

## Supported Pipelines
- Logs

## Prerequisites
- Splunk admin credentials
- Configured storage extension

## Configuration
| Field               | Type     | Default                                                                                         | Description                                                                                                                                                             |
|---------------------|----------|-------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| endpoint                   | string   | `required` `(no default)`                                                                 | The endpoint of the splunk instance to collect from.                                                                                                                                         |
| splunk_username            | string   | `(no default)`                                                                           | Specifies the username used to authenticate to Splunk using basic auth.                                                                                                           |
| splunk_password            | string   | `(no default)`                                                                           | Specifies the password used to authenticate to Splunk using basic auth.                                                                                                           |
| auth_token                    | string   | `(no default)`                                                                           | Specifies the token used to authenticate to Splunk using token auth. |
| token_type                    | string   | `(no default)`                                                                           | Specifies the type of token used to authenticate to Splunk using token auth. Accepted values are "Bearer" or "Splunk". |
| job_poll_interval | duration | `5s` | The receiver uses an API call to determine if a search has completed. Specifies how long to wait between polling for search job completion. |
| searches.query | string | `required (no default)` | The Splunk search to run to retrieve the desired events. Queries must start with `search` and should not contain additional commands, nor any time fields (e.g. `earliesttime`) |
| searches.earliest_time | string | `required (no default)` | The earliest timestamp to collect logs. Only logs that occurred at or after this timestamp will be collected. Must be in ISO 8601 or RFC3339 format. |
| searches.latest_time | string | `required (no default)` | The latest timestamp to collect logs. Only logs that occurred at or before this timestamp will be collected. Must be in ISO 8601 or RFC3339 format. |
| searches.event_batch_size | int | `100` | The amount of events to query from Splunk for a single request. |
| storage | component | `required (no default)` | The component ID of a storage extension which can be used when polling for `logs`. The storage extension prevents duplication of data after an exporter error by remembering which events were previously exported. |

### Example Configuration
```yaml
receivers:
  splunksearchapi:
    endpoint: "https://splunk-c4-0.example.localnet:8089"
    tls:
      insecure_skip_verify: true
    splunk_username: "user"
    splunk_password: "pass"
    job_poll_interval: 5s
    searches:
      - query: 'search index=my_index'
        earliest_time: "2024-11-01T01:00:00.000-05:00"
        latest_time: "2024-11-30T23:59:59.999-05:00"
        event_batch_size: 500
    storage: file_storage
exporters:
  googlecloud:
    project: "my-gcp-project"
    log: 
      default_log_name: "splunk-events"
    sending_queue:
      enabled: false

extensions:
  file_storage:
    directory: "./local/storage"

service:
  extensions: [file_storage]
  pipelines:
    logs:
      receivers: [splunksearchapi]
      exporters: [googlecloud]
```