# Cisco ASA Plugin

Log parser for Cisco ASA

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| listen_port | A port which the agent will listen for syslog messages | int | `5140` | false |  |
| listen_ip | A syslog ip address | string | `0.0.0.0` | false |  |

## Example Config:

Below is an example of a basic conifg

```yaml
receivers:
  plugin:
    path: ./plugins/cisco_asa_logs.yaml
    parameters:
      listen_port: 5140
      listen_ip: 0.0.0.0
```
