# SAP Netweaver Receiver

This receiver collects metrics from SAP Netweaver instance based on the [SAPControl Web Service Interface](https://www.sap.com/documents/2016/09/0a40e60d-8b7c-0010-82c7-eda71af511fa.html). Stats are collected from the `GetAlertTree` and `EnqGetLockTable` methods.

## Prerequisites

This receiver supports SAP Netweaver 7.10+

## Configuration

The following settings are optional:
- `metrics` (default: see `DefaultMetricsSettings` [here](./internal/metadata/generated_metrics.go): Allows enabling and disabling specific metrics from being collected in this receiver.
- `endpoint` (default = `http://localhost:50013`): The default URL for SAP Netweaver.
- `username` (no default): Specifies the username used to authenticate using basic auth.
- `password` (no default): Specifies the password used to authenticate using basic auth.
- `collection_interval` (default = `10s`): This receiver collects metrics on an interval. This value must be a string readable by Golang's [time.ParseDuration](https://pkg.go.dev/time#ParseDuration).

### Example Configuration

```yaml
receivers:
  sapnetweaver:
    metrics:
    endpoint: http://localhost:50013
    username: root
    password: password
    collection_interval: 10s
```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml).

## Metrics

Details about the metrics produced by this receiver can be found in [metadata.yaml](./metadata.yaml)
