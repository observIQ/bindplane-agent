# macOS Installation

### Installing
The agent may be installed through a shell script.w

This script may also be used to update an existing installation.

To install using the installation script, you may run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-agent/releases/latest/download/install_macos.sh)" install_macos.sh
```

Installation artifacts are signed information on verifying the signature can be found at [Verifying Artifact Signatures](../verify-signature.md).

#### Managed Mode

To install the agent with an OpAMP connection configuration set the following flags. 

```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-agent/releases/latest/download/install_macos.sh)" install_macos.sh -e <your_endpoint> -s <secret-key>
```

To read more about the generated connection configuration file see [OpAMP docs](./opamp.md).

## Configuring the Agent

After installing the `observiq-otel-collector` you can change the configuration file printed out at the end of the installation.

The default configuration file can be found at `/opt/observiq-otel-collector/config.yaml`.

After changing the configuration file run `sudo launchctl unload /Library/LaunchDaemons/com.observiq.collector.plist; sudo launchctl load /Library/LaunchDaemons/com.observiq.collector.plist` for the changes to take effect.

For more information on configuring the agent, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

**Logging**

Logs from the agent will appear in `/opt/observiq-otel-collector/log`. You may run `sudo tail -F /opt/observiq-otel-collector/log/collector.log` to view them.

Stderr for the agent process can be found at `/var/log/observiq_collector.err`.

## Agent Services Commands

The agent uses `launchctl` to control the agent lifecycle using the following commands.

```sh
# Start the agent
sudo launchctl load /Library/LaunchDaemons/com.observiq.collector.plist

# Stop the agent
sudo launchctl unload /Library/LaunchDaemons/com.observiq.collector.plist
```

## Uninstalling

To uninstall an installation made with the install script, run:
```sh
sudo sh -c "$(curl -fsSlL https://github.com/observiq/bindplane-agent/releases/latest/download/install_macos.sh)" install_macos.sh -r
```
