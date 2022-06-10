# Syslog with Google Cloud

The Syslog receiver can be used to receive Syslog and forward to Google Cloud Logging.

## Prerequisites

See the [prerequisites](../README.md) doc for Google Cloud prerequisites.

## Examples

**Rsyslog RFC 3164**

Configure [Rsyslog forwarding](https://www.rsyslog.com/sending-messages-to-a-remote-syslog-server/) to the agent system. In this example,
the collector is installed on the Rsyslog system.

```
*.* action(type="omfwd" target="localhost" port="5140" protocol="udp")
```

This example will listen on localhost port `5140/udp` for RFC 3164 syslog, receiving syslog from
the local rsyslog service.

```yaml
syslog:
  udp:
    listen_address: "0.0.0.0:5140"
  protocol: rfc3164
```

**TCP RFC 5424**

This example will listen on all interfaces with port `54526/tcp` for RFC 5424 syslog.

```yaml
syslog:
  tcp:
    listen_address: "0.0.0.0:54526"
  protocol: rfc5424
```
