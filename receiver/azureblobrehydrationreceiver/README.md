# Azure Blob Storage Rehydration Receiver
Rehydrates OTLP from Azure Blob Storage that was stored using the Azure Blob Exporter [../../exporter/azureblobexporter/README.md].

## Important Note
This is not a traditional receiver that continually produces data but rather rehydrates all blobs found within a specified time range. Once all of the blobs have been rehydrated in that time range the receiver will stop producing data.

## Minimum Agent Versions
- Introduced: [v1.36.0](https://github.com/observIQ/bindplane-agent/releases/tag/v1.36.0)

## Supported Pipelines
- Metrics
- Logs
- Traces

## How it works
1. The receiver polls blob storage for all blobs in the specified container and prefix, if provided.
2. The receiver will parse each blob's path to determine if it matches the path created by the [Azure Blob Exporter](../../exporter/azureblobexporter/README.md#blob-path).
3. If the blob path matches the exporter's path, the receiver will parse the timestamp represented by the path.
4. If the timestamp is within the configured range the receiver will download the blob and parse its contents into OTLP data.

## Configuration
| Field              | Type      | Default          | Required | Description                                                                                                                    |
|--------------------|-----------|------------------|----------|--------------------------------------------------------------------------------------------------------------------------------|
| connection_string  |  string   |                  | `true`   | The connection string to the Azure Blob Storage account. Can be found under the `Access keys` section of your storage account.                                                         |
| container          |  string   |                  | `true`   | The name of the container to rehydrate from.                                                                                                                                           |
| root_folder        |  string   |                  | `false`  | The root folder that prefixes the blob path. Should match the `root_folder` value of the Azure Blob Exporter.                                                                          |
| starting_time      |  string   |                  | `true `  | The UTC start time that represents the start of the time range to rehydrate from. Must be in the form `YYYY-MM-DDTHH:MM`.                                                              |
| ending_time        |  string   |                  | `true `  | The UTC end time that represents the end of the time range to rehydrate from. Must be in the form `YYYY-MM-DDTHH:MM`.                                                                  |
| delete_on_read     |  bool     | `false`          | `false ` | If `true` the blob will be deleted after being rehydrated.                                                                                                                             |
| poll_interval      |  string   | `1m`             | `false ` | How often to read a new set of blobs. This value is mostly to control how often the blob API is called to ensure once rehydration is done the receiver isn't making to many API calls. |
| storage            |  string   |                  | `false ` | The component ID of a storage extension. The storage extension prevents duplication of data after a collector restart by remembering which blobs were previously rehydrated.           |
