# `aerospike` plugin

The `aerospike` plugin consumes [Aerospike](https://aerospike.com/) log entries from the local filesystem and outputs parsed entries.

## Supported Platforms

- Linux
- Windows
- MacOS
- Kubernetes

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `journald_directory` | `string` | `"/var/log/aerospike/aerospike.log"` | The absolute path to the general Aerospike logs |
| `start_at` | `enum` | `end` | Start reading file from 'beginning' or 'end' |

## Prerequisites

To set up logs in Aerospike, check out their documentation [here](https://download.aerospike.com/docs/operations/manage/log/index.html).

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: aerospike
- type: stdout

```

With non-standard log paths:

```yaml
pipeline:
- type: aerospike
  journald_directory: "/path/to/logs"
- type: stdout

```
