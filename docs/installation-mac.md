# macOS Installation

### Installing
The collector may be installed through a shell script.

This script may also be used to update an existing installation.

To install using the installation script, you may run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh
```

## Configuring the Collector

After installing the `observiq-otel-collector` you can change the configuration file printed out at the end of the installation.

The default configuration file can be found at `/opt/observiq-otel-collector/config.yaml`.

After changing the configuration file run `sudo launchctl stop com.observiq.collector; sudo launchctl start com.observiq.collector` for the changes to take effect.

For more information on configuring the collector, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

## Collector Services Commands

The collector uses `launchctl` to control the collector lifecycle using the following commands.

```sh
# Start the collector
sudo launchctl start com.observiq.collector

# Stop the collector
sudo launchctl stop com.observiq.collector
```

## Uninstalling

To uninstall an installation made with the install script, run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh -r
```
