# macOS Installation

## Installing

The installation for macOS is supported via [Homebrew](https://brew.sh/).

```sh
brew tap observiq/homebrew-observiq-otel-collector
brew update
brew install observiq/observiq-otel-collector/observiq-otel-collector
```

You can then run `brew services start observiq/observiq-otel-collector/observiq-otel-collector` to start the collector with the supplied configuration.

To verify the collector is working you can check the output at `/tmp/observiq-otel-collector.log`.

### Install/Update script
The collector may be installed through a shell script which will wrap brew commands.

This script may also be used to update an existing installation.

To install using the installation script, you may run:
```sh
sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh
```

### Additional Install Steps

If you plan on collecting metrics via the [JMX Receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v0.50.0/receiver/jmxreceiver/README.md) you can copy the `opentelemetry-java-contrib-jmx-metrics.jar` to the default location to make configuration easier.

```sh
sudo cp $(brew --prefix observiq/observiq-otel-collector/observiq-otel-collector)/lib/opentelemetry-java-contrib-jmx-metrics.jar /opt
```
## Configuring the Collector

After installing, the `observiq-otel-collector` you can change the configuration file printed out at the end of the brew installation.

The default configuration file can be found at `$(brew --prefix observiq/observiq-otel-collector/observiq-otel-collector)/config.yaml`.

After changing the configuration file run `brew services restart observiq/observiq-otel-collector/observiq-otel-collector` for the changes to take effect.

For more information on configuring the collector, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

## Collector Services Commands

The collector uses `brew services` the following commands control the collector lifecycle.

```sh
# Start the collector
launchctl start com.observiq.collector

# Stop the collector
launchctl stop com.observiq.collector
```

## Uninstalling

### Uninstall brew

To uninstall the collector run the following commands:

```sh
brew uninstall observiq/observiq-otel-collector/observiq-otel-collector

# To remove the plist file
launchctl remove com.observiq.collector

# If you moved the opentelemetry-java-contrib-jmx-metrics.jar
sudo rm /opt/opentelemetry-java-contrib-jmx-metrics.jar
```

### Uninstall script

To uninstall an installation made with the install script, run:
```sh
sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh -r
```
