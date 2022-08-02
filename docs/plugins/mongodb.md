# `mongodb` plugin

The `mongodb` plugin consumes [MongoDB](https://www.mongodb.com/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `source` | `file` | Use this field to specify where your logs are coming from. When choosing the 'file' option, the agent reads in logs from the log paths specified below.  When choosing the 'Kubernetes' options, the agent reads logs from /var/log/containers based on the Pod and Container specified below. |
| `log_path` | `"/var/log/mongodb/mongodb.log*"` | The path of the log file |
| `cluster_name` | `""`  | Friendly name to be used as a resource label. Only relevant if the source is "kubernetes". |
| `pod_name` | `mongodb` | The pod name (without the unique identifier on the end). Only relevant if the source is "kubernetes". |
| `container_name` | `"*"` | The container name of the Mongodb container. Only relevant if the source is "kubernetes". |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using default file log path:

```yaml
pipeline:
- type: mongodb
- type: stdout

```

Using Kubernetes:

```yaml
pipeline:
- type: mongodb
  source: kubernetes
  cluster_name: "stage"
  pod_name: mongodb
  container_name: "*"
- type: stdout

```
