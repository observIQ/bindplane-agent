# Mysql Metrics with Google Cloud

The Mysql Receiver can be used to send Mysql metrics to Google Cloud Monitoring. See the [documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/mysqlreceiver) for more details.

## Limitations

The collector must be installed on the Mysql system.

## Prerequisites

See the [prerequisites section](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/mysqlreceiver#prerequisites) for Mysql prerequisites.

See the [prerequisites](../prerequisites.md) doc for Google Cloud prerequisites.

## Authentication Environment Variables

The configuration assumes the following environment variables are set:
- `MYSQL_USERNAME`
- `MYSQL_PASSWORD`

Set the variables by creating a [systemd override](https://wiki.archlinux.org/title/systemd#Replacement_unit_files).

Run the following command
```bash
sudo systemctl edit observiq-otel-collector
```

If this is the first time an override is being created, the file will be empty. Paste the following contents into the file. If the `Service` section is already present, append the two `Environment` lines to the `Service` section.
```
[Service]
Environment=MYSQL_USERNAME=otel
Environment=MYSQL_PASSWORD=password
```

After restarting the collector, the configuration will attempt to use the username:password `otel:otel`.

```bash
sudo systemctl restart observiq-otel-collector
```
