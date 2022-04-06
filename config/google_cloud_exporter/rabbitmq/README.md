# RabbitMQ Metrics with Google Cloud

The RabbitMQ Receiver can be used to send RabbitMQ metrics to Google Cloud Monitoring. See the [documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/rabbitmqreceiver/README.m) for more details.

## Prerequisites

See the [prerequisites section](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/rabbitmqreceiver/README.md#prerequisites) for RabbitMQ prerequisites.

See the [prerequisites](../prerequisites.md) doc for Google Cloud prerequisites.

## Authentication Environment Variables

The configuration assumes the following environment variables are set:
- `RABBITMQ_USERNAME`
- `RABBITMQ_PASSWORD`

Set the variables by creating a [systemd override](https://wiki.archlinux.org/title/systemd#Replacement_unit_files).

Run the following command
```bash
sudo systemctl edit observiq-otel-collector
```

If this is the first time an override is being created, the file will be empty. Paste the following contents into the file. If the `Service` section is already present, append the two `Environment` lines to the `Service` section.

Replace `otel` with your RabbitMQ username and password.
```
[Service]
Environment=RABBITMQ_USERNAME=otel
Environment=RABBITMQ_PASSWORD=otel
```

After restarting the collector, the configuration will attempt to use the configured username and password.

```bash
sudo systemctl restart observiq-otel-collector
```
