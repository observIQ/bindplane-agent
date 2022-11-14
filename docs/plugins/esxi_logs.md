# VMware ESXi Plugin

Log parser for VMware ESXi

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| listen_port | A port which the agent will listen for ESXi messages | int | `5140` | false |  |
| listen_ip | The local IP address to listen for ESXi connections on | string | `0.0.0.0` | false |  |
| enable_tls | Enable TLS for the vCenter listener | bool | `false` | false |  |
| certificate_file | Path to TLS certificate file | string | `/opt/cert` | false |  |
| private_key_file | Path to TLS private key file | string | `/opt/key` | false |  |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/esxi_logs.yaml
    parameters:
      listen_port: 5140
      listen_ip: 0.0.0.0
      enable_tls: false
      certificate_file: /opt/cert
      private_key_file: /opt/key
```
