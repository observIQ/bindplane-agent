# Cadvisor Metrics with Google Cloud

The Prometheus Receiver can be used to send [Cadvisor](https://github.com/google/cadvisor) metrics to Google Cloud Monitoring.

## Prerequisites

**Google Cloud**

See the [prerequisites](../prerequisites.md) doc for Google Cloud prerequisites.

## Deployment

Run as a container

```bash
docker run -d \
    --restart=always \
    --volume="$(pwd)/config.yaml:/etc/otel/config.yaml" \
    --volume="$(pwd)/credentials.json:/etc/otel/credentials.json" \
    -e "GOOGLE_APPLICATION_CREDENTIALS=/etc/otel/credentials.json" \
    -e "DOCKER_AGENT_HOSTNAME=$(hostname)" \
    observiq/observiq-otel-collector:v0.4.1
```

Run with Docker Compose

```yaml
version: "3.9"
services:
  collector:
    restart: always
    # Run as root if using a configuration that requires
    # root privileges.
    #user: root
    image: observiq/observiq-otel-collector:v0.4.1
    restart: always
    deploy:
      resources:
        limits:
          cpus: 0.50
          memory: 256M
    environment:
      - GOOGLE_APPLICATION_CREDENTIALS=/etc/otel/credentials.json
      - DOCKER_AGENT_HOSTNAME=${HOSTNAME}
    volumes:
      - ${PWD}/config.yaml:/etc/otel/config.yaml
      - ${PWD}/credentials.json:/etc/otel/credentials.json
```
