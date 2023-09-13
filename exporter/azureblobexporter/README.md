# Azure Blob Exporter

This exporter allows you to export metrics, traces, and logs to Azure Blob Storage. Telemetry is exported in [OpenTelemetry Protocol JSON format](https://github.com/open-telemetry/opentelemetry-proto).


## Configuration
| Field              | Default          | Required | Description                                                                                                                    |
|--------------------|------------------|----------|--------------------------------------------------------------------------------------------------------------------------------|
| connection_string  |                  | `true`   | The connection string to the Azure Blob Storage account. Can be found under the `Access keys` section of your storage account. |
| container          |                  | `true`   | The name of the container in the storage account to save to.                                                                   |
| blob_prefix        |                  | `false`  | An optional prefix to prepend to blob file name.                                                                               |
| root_folder        |                  | `false`  | An optional root folder that will prefix the blob path.                                                                        |
| partition          | `minute`         | `true`   | Time granularity of blob name. Valid values are `hour` or `minute`.                                                            |
| compression        | `none`           | `false`  | The type of compression applied to the data before sending it to storage. Valid values are `none` and `gzip`.                  |

Blog paths will be in the form:

```
{root_folder}/year=XXXX/month=XX/day=XX/hour=XX/minute=XX
```

## Example Configurations

### Minimal Configuration

```yaml
azureblob:
    connection_string: "DefaultEndpointsProtocol=https;AccountName=storage_account_name;AccountKey=storage_account_key;EndpointSuffix=core.windows.net"
    container: "my-container"
```

Example Blob Names:

```
year=2021/month=01/day=01/hour=01/minute=00/metrics_{random_id}.json
year=2021/month=01/day=01/hour=01/minute=00/logs_{random_id}.json
year=2021/month=01/day=01/hour=01/minute=00/traces_{random_id}.json
```


### Hour Partition Configuration

```yaml
azureblob:
    connection_string: "DefaultEndpointsProtocol=https;AccountName=storage_account_name;AccountKey=storage_account_key;EndpointSuffix=core.windows.net"
    container: "my-container"
    partition: "hour"
```

Example Blob Names:

```
year=2021/month=01/day=01/hour=01/metrics_{random_id}.json
year=2021/month=01/day=01/hour=01/logs_{random_id}.json
year=2021/month=01/day=01/hour=01/traces_{random_id}.json
```

### Full Configuration with compression

```yaml
azureblob:
    connection_string: "DefaultEndpointsProtocol=https;AccountName=storage_account_name;AccountKey=storage_account_key;EndpointSuffix=core.windows.net"
    container: "my-container"
    root_folder: "otel"
    blob_prefix: "linux"
    partition: "minute"
    compression: "gzip"
```

Example Blob Names:

```
otel/year=2021/month=01/day=01/hour=01/minute=00/linuxmetrics_{random_id}.json.gz
otel/year=2021/month=01/day=01/hour=01/minute=00/linuxlogs_{random_id}.json.gz
otel/year=2021/month=01/day=01/hour=01/minute=00/linuxtraces_{random_id}.json.gz
```
