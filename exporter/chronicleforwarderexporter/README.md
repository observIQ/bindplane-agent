# Chronicle Forwarder Exporter

The Chronicle Forwarder Exporter is designed for forwarding logs to a Chronicle endpoint using either Syslog or File-based methods. This exporter supports customization of data export types and various configuration options to tailor the connection and data handling to specific needs.

## Minimum Agent Versions

- Introduced: [v1.42.0](https://github.com/observIQ/bindplane-agent/releases/tag/v1.42.0)

## Supported Pipelines

- Logs

## How It Works

1. For Syslog, it establishes a network connection to the specified Chronicle forwarder endpoint.
2. For File, it writes logs to a specified file path.

## Configuration

| Field                | Type         | Default Value | Required | Description                                       |
| -------------------- | ------------ | ------------- | -------- | ------------------------------------------------- |
| export_type          | string       | `syslog`      | `true`   | Type of export, either `syslog` or `file`.        |
| syslog               | SyslogConfig |               | `true`   | Configuration for Syslog connection.              |
| file                 | File         |               | `true`   | Configuration for File connection.                |
| raw_log_field        | string       |               | `false`  | The field name to send raw logs to Chronicle.     |
| syslog.host          | string       | `127.0.0.1`   | `false`  | The Chronicle forwarder endpoint host.            |
| syslog.port          | int          | `10514`       | `false`  | The port to send logs to.                         |
| syslog.network       | string       | `tcp`         | `false`  | The network protocol to use (e.g., `tcp`, `udp`). |
| syslog.tls.key_file  | string       |               | `false`  | Configure the receiver to use TLS.                |
| syslog.tls.cert_file | string       |               | `false`  | Configure the receiver to use TLS.                |
| file.path            | string       |               | `false`  | The path to the file for storing logs.            |

## Example Configurations

### Syslog Configuration Example

```yaml
chronicleforwarderexporter:
  export_type: "syslog"
  syslog:
    host: "syslog.example.com"
    port: 10514
    network: "tcp"
```

### File Configuration Example

```yaml
chronicleforwarderexporter:
  export_type: "file"
  file:
    path: "/path/to/logfile"
```

---
