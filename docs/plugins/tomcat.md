# `tomcat` plugin

The `tomcat` plugin consumes [Apache Tomcat](https://tomcat.apache.org/) log entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `source` | `file` | Use this field to specify where your logs are coming from. When choosing the 'file' option, the agent reads in logs from the log paths specified below. When choosing the 'Kubernetes' options, the agent reads logs from /var/log/containers based on the Pod and Container specified below. |
| `log_format` | `default`  | When choosing the 'default' option, the agent will expect and parse logs in a format that matches the default logging configuration. When choosing the 'observIQ' option, the agent will expect and parse logs in an optimized JSON format that adheres to the observIQ specification, requiring an update to the server.xml file. |
| `enable_access_log` | `true` | Enable to collect Apache Tomcat access logs |
| `access_log_path` | `"/usr/local/tomcat/logs/localhost_access_log.*.txt"` | Path to access log file |
| `enable_catalina_log` | `true` | Enable to collect Apache Tomcat catalina logs |
| `catalina_log_path` | `"/usr/local/tomcat/logs/catalina.out"` | Path to catalina log file |
| `cluster_name` | `""` | Friendly name to be used as a resource label. Only relevant if the source is "kubernetes". |
| `pod_name` | `'tomcat-*'` | The pod name (without the unique identifier on the end). Only relevant if the source is "kubernetes". |
| `container_name` | `"*"` | The container name of the Nginx container. Only relevant if the source is "kubernetes". |
| `start_at` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using default file log paths:

```yaml
pipeline:
- type: tomcat
- type: stdout

```

Using non-default file parameters:

```yaml
pipeline:
- type: tomcat
  source: file
  log_format: default
  enable_access_log: true
  access_log_path: "path/to/logs"
  enable_catalina_log: true
  catalina_log_path: "path/to/logs"
- type: stdout

```

Using Kubernetes:

```yaml
pipeline:
- type: tomcat
  source: kubernetes
  cluster_name: "stage"
  pod_name: 'tomcat-*'
  container_name: "*"
- type: stdout

```
