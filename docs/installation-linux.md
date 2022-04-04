# Linux Installation

## Installing

Installation is done through deb and rpm packages. Installing the collector will also install the `observiq-otel-collector` service on systemd systems.

### Install/Update script
The collector may be installed through a shell script which will automatically determine which package to install.

This script may also be used to update an existing installation.

To install using the installation script, you may run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_unix.sh)" install_unix.sh
```

### RPM Installation
First download the RPM package for your architecture from the [releases page](https://github.com/observIQ/observiq-otel-collector/releases/latest).
Then you may install the package using `rpm`, see this example for installing the amd64 package:
```sh
sudo rpm -U ./observiq-otel-collector_linux_amd64.rpm
sudo systemctl enable --now observiq-otel-collector
```

### DEB Installation
First download the DEB package for your architecture from the [releases page](https://github.com/observIQ/observiq-otel-collector/releases/latest).
Then you may install the package using `dpkg`, see this example for installing the amd64 package:
```sh
sudo dpkg -i ./observiq-otel-collector_linux_amd64.deb
sudo systemctl enable --now observiq-otel-collector
```

## Configuring the Collector
After installing, systems with systemd installed will have the `observiq-otel-collector` service up and running!

Logs from the collector will appear in journald. You may run `journalctl -u observiq-otel-collector.service` to view them.

The config file for the collector can be found at `/opt/observiq-otel-collector/config.yaml`. When changing the configuration,the collector service must be restarted in order for config changes to take effect.

For more information on configuring the collector, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

By default, the `observiq-otel-collector` service runs as the "observiq-otel-collector" user. Some OpenTelemetry components may require root permissions.
To run the collector as root, you may create a systemd override.

Run `sudo systemctl edit observiq-otel-collector` and paste the following config:
```
[Service]
User=root
Group=root
```

Restart the collector for these changes to take effect.

## Restarting the Collector
On systemd systems, the collector may be restarted with the following command:
```sh
systemctl restart observiq-otel-collector
```

## Stopping the Collector
On systemd systems, the collector may be stopped with the following command:
```sh
systemctl stop observiq-otel-collector
```

## Starting the Collector
On systemd systems, the collector may be started with the following command:
```sh
systemctl start observiq-otel-collector
```

## Uninstalling

### RPM Uninstall

To uninstall the rpm package, run:
```sh
sudo rpm -e observiq-otel-collector
```

### DEB Uninstall

To uninstall the deb package, run:
```sh
sudo dpkg -r observiq-otel-collector
```

### Uninstall script

To uninstall an installation made with the install script, run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_unix.sh)" install_unix.sh -r
```
