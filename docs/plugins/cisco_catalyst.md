# `cisco_catalyst` plugin

The `cisco_catalyst` plugin receives logs from [Cisco Catalyst](https://en.wikipedia.org/wiki/Cisco_Catalyst) network devices and outputs a parsed entry.

## Supported Platforms

- Linux
- Windows
- MacOS
- Kubernetes

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `listen_port` | `int` | `514` | A port which the agent will listen for udp messages |
| `listen_ip` | `string` | `"0.0.0.0"`  | A UDP ip address of the form `<ip>` |
| `add_labels` | `bool` | `true` | Adds net.transport, net.peer.ip, net.peer.port, net.host.ip and net.host.port labels. |

## Prerequisites

No prerequisite actions required.

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: cisco_catalyst
- type: stdout

```

With non-standard port and IP:

```yaml
pipeline:
- type: cisco_catalyst
  listen_port: 6514
  listen_ip: "10.0.0.1"
- type: stdout

```
