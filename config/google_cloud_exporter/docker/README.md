# Docker Metrics with Google Cloud

The Docker Receiver can be used to send Docker metrics to Google Cloud Monitoring.

## Prerequisites

**Google Cloud**

See the [prerequisites](../prerequisites.md) doc for Google Cloud prerequisites.

**Docker Socket**

The provided configuration assumes the collector is running on the Docker system. By default, the `endpoint` collector from is `unix:///var/run/docker.sock`.

The user running the collector must have permission to read the [docker socket](https://docs.docker.com/engine/install/linux-postinstall/).

```bash
sudo usermod -aG docker observiq-otel-collector
```
