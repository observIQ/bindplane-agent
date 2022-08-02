# `active_directory` plugin

The `active_directory` plugin consumes [Active Directory](https://en.wikipedia.org/wiki/Active_Directory) log entries from the local filesystem and outputs parsed entries.

## Supported Platforms

- Windows

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `enable_dns_server` | `bool` | `true` | Enable to collect DNS server logs |
| `enable_dfs_replication` | `bool` | `true` | Enable to collect DFS replication logs |
| `enable_file_replication` | `bool` | `false` | Enable to collect file replication logs |
| `poll_interval` | `string` | `1s` | Set the rate that logs are being collected  |
| `max_reads` | `int` | `1000` | Maximum number of logs collected |
| `start_at` | `enum` | `end` | Start reading file from 'beginning' or 'end' |

## Example usage

### Configuration

Using defaults:

```yaml
pipeline:
- type: active_directory
- type: stdout

```

With file replication enabled:

```yaml
pipeline:
- type: active_directory
  enable_file_replication: true
- type: stdout

```
