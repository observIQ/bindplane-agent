# macOS Installation

## Installing

The installation for macOS is supported via [Homebrew](https://brew.sh/).

```sh
brew tap observiq/homebrew-observiq-otel-collector
brew update
brew install observiq-otel-collector
```

You can then run `brew services start observiq-otel-collector` to start the collector with the supplied configuration.

To verify the collector is working you can check the output at `/tmp/observiq-otel-collector.err`.

### Additional Install Steps

If you plan on JMX collecting metrics via the [JMX Receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v0.45.1/receiver/jmxreceiver/README.md) you can copy the `opentelemetry-java-contrib-jmx-metrics.jar` to the default location to make configuration easier.

```sh
sudo cp $(brew --prefix observiq-otel-collector)/lib/opentelemetry-java-contrib-jmx-metrics.jar /opt
```
## Configuring the Collector

After installing, the `observiq-otel-collector` you can change the configuration file printed out at the end of the brew installation.

The default configuration file can be found at `$(brew --prefix observiq-otel-collector)/config.yaml`.

After changing the configuration file run `brew services restart observiq-otel-collector` for the changes to take effect.

For more information on configuring the collector, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

## Collector Services Commands

The collector uses `brew services` the following commands control the collector lifecycle.

```sh
# Start the collector
brew services start observiq-otel-collector

# Stop the collector
brew services stop observiq-otel-collector

# Restart the collector
brew services restart observiq-otel-collector
```

## Uninstalling

To uninstall the collector run the following commands:

```sh
brew services stop observiq-otel-collector
brew uninstall observiq-otel-collector

# If you moved the opentelemetry-java-contrib-jmx-metrics.jar
sudo rm /opt/opentelemetry-java-contrib-jmx-metrics.jar
```
