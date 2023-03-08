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

#### Managed Mode

To install the collector with an OpAMP connection configuration set the following flags. 

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_unix.sh)" install_unix.sh -e <your_endpoint> -s <secret-key>
```

To read more about the generated connection configuration file see [OpAMP docs](./opamp.md).

### Installation from local package

To install the collector from a local package use the `-f` with the path to the package.

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_unix.sh)" install_unix.sh -f <path_to_package>
```

### RPM Installation
First download the RPM package for your architecture from the [releases page](https://github.com/observIQ/observiq-otel-collector/releases/latest).
Then you may install the package using `rpm`, see this example for installing the amd64 package:

**Note**: Replace `${VERSION}` with the version of the package you downloaded.

```sh
sudo rpm -U ./observiq-otel-collector_v${VERSION}_linux_amd64.rpm
sudo systemctl enable --now observiq-otel-collector
```

### DEB Installation
First download the DEB package for your architecture from the [releases page](https://github.com/observIQ/observiq-otel-collector/releases/latest).
Then you may install the package using `dpkg`, see this example for installing the amd64 package:

**Note**: Replace `${VERSION}` with the version of the package you downloaded.

```sh
sudo dpkg -i ./observiq-otel-collector_v${VERSION}_linux_amd64.deb
sudo systemctl enable --now observiq-otel-collector
```

## Configuring the Collector
After installing, systems with systemd installed will have the `observiq-otel-collector` service up and running!

**Logging**

Logs from the collector will appear in `/opt/observiq-otel-collector/log`. You may run `sudo tail -F /opt/observiq-otel-collector/log/collector.log` to view them.

Stdout and stderr for the collector process are recorded via journald. You man run `sudo journalctl -u observiq-otel-collector.service` to view them.

**Configuration**

The config file for the collector can be found at `/opt/observiq-otel-collector/config.yaml`. When changing the configuration,the collector service must be restarted in order for config changes to take effect.

For more information on configuring the collector, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

**Permissions**

By default, the `observiq-otel-collector` service runs as the "root" user. Some OpenTelemetry components require root permissions in order to read log files owned by other users.

It may be desirable to run the collector as an unprivileged user. For example, a metrics only collector does not require root access.

To run the collector as the `observiq-otel-collector` user, you may create a systemd override.

Run `sudo systemctl edit observiq-otel-collector` and paste the following config:
```
[Service]
User=observiq-otel-collector
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
