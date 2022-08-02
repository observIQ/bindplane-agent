# `cisco_meraki` plugin

The `cisco_meraki` plugin receives logs from [Cisco Meraki](https://meraki.cisco.com/) network devices and outputs a parsed entry.

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- |--- | --- |
| `listen_port` | `int` | `514` | A port which the agent will listen for udp messages |
| `listen_ip` | `string` | `"0.0.0.0"`  | A UDP ip address of the form `<ip>` |

## Prerequisites

It may be necessary to add an inbound firewall rule.

### Windows

- Navigate to Windows Firewall Advanced Settings, and then Inbound Rules
- Create a new rule and set the Rule Type to "Port"
- For Protocol and Ports, select "UDP" and a specific local port of 514
- For Action, select "Allow the connection"
- For Profile, apply to "Domain", "Private", and "Public"
- Set a name to easily identify rule, such as "Allow Syslog Inbound Connections to 514 UDP"

 ### Linux

- Using Firewalld:
```shell
firewall-cmd --permanent --add-port=514/udp
firewall-cmd --reload
```
- Using UFW:
```shell
ufw allow 514
```

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: cisco_meraki
- type: stdout

```

With non-standard port and IP:

```yaml
pipeline:
- type: cisco_meraki
  listen_port: 6514
  listen_ip: "10.0.0.1"
- type: stdout

```
