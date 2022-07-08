# macOS Installation

### Installing
The collector may be installed through a shell script.

This script may also be used to update an existing installation.

To install using the installation script, you may run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh
```

#### Managed Mode

To install the collector with an OpAMP connection configuration set the following flags. 

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh -e <your_endpoint> -s <secret-key>
```

To read more about the generated connection configuration file see [OpAMP docs](./opamp.md).

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
