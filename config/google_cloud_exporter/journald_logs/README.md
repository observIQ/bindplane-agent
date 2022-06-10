# Journald Logging with Google Cloud

The Journald receiver can be used to send [Journald](https://wiki.archlinux.org/title/Systemd/Journal) logs to Google Cloud Logging.

## Limitations

The collector must be installed on the Journald system. The system must use Systemd as the init system, as Journald is the central logging solution for Systemd equipped systems. You can validate by running the `journalctl` command.

## Prerequisites

See the [prerequisites](../README.md) doc for Google Cloud prerequisites.
