# File Logging with Google Cloud

The Filelog receiver can be used to read log files and send the entries to Google Cloud Logging.

## Limitations

The collector must be installed on the target system. The user running the collector must have permission to read the files targeted by the Filelog receiver's configuration.

## Prerequisites

See the [prerequisites](../README.md) doc for Google Cloud prerequisites.

Edit the configuration and replace `project: REPLACE_ME` with the project id you wish to write logs to.

## Examples

**Raw log file**

Filelog can read raw log files (without parsing) and forward to Google Cloud Logging. This a
basic implementation. Generally, you will want to parse the log entries in order to make searching
easier within Google Cloud Logging.

See [config.yaml](./config.yaml).

**Regex**

Filelog can read raw log files, parse them with regex, and forward to Google Cloud Logging.

See [config_regex.yaml](./config_regex.yaml) for an example using Redis.

**Json**

Filelog can read raw log files, parse them as json, and forward to Google Cloud Logging.

See [config_json.yaml](./config_json.yaml) for an example using arbitrary json.
