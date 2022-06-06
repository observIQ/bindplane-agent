# Big IP F5 with Google Cloud

The bigip Receiver can be used to metrics to Google Cloud Monitoring.

## Prerequisites

See the [prerequisites](../README.md) doc for Google Cloud prerequisites.

## Limitations

To avoid duplicate metrics, ensure each collector is scraping a single f5 system. This ensures
there is a unique `node_id` (hostname of the collector) label for each f5 system being monitored.

## Configuration

An example configuration is located [here](./config.yaml).

1. Copy [config.yaml](./config.yaml) to `/opt/observiq-otel-collector/config.yaml`
2. Follow the [authentication section](./README.md#authentication-environment-variables) for configuring username and password.
2. Restart the collector: `sudo systemctl restart observiq-otel-collector`

## Authentication Environment Variables

The configuration assumes the following environment variables are set:
- `F5_USERNAME`
- `F5_PASSWORD`

Set the variables by creating a [systemd override](https://wiki.archlinux.org/title/systemd#Replacement_unit_files).

Run the following command
```bash
sudo systemctl edit observiq-otel-collector
```

If this is the first time an override is being created, the file will be empty. Paste the following contents into the file. If the `Service` section is already present, append the two `Environment` lines to the `Service` section.

Replace `otel` with your Big IP F5 username and password.
```
[Service]
Environment=F5_USERNAME=otel
Environment=F5_PASSWORD=otel
```

After restarting the collector, the configuration will attempt to use the username:password `otel:otel`.

```bash
sudo systemctl restart observiq-otel-collector
```
