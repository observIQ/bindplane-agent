# macOS Plugin

Log parser for macOS

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| enable_system_log | Enable to collect MacOS system logs | bool | `true` | false |  |
| system_log_path | The absolute path to the System log | []string | `[/var/log/system.log]` | false |  |
| enable_install_log | Enable to collect MacOS install logs | bool | `true` | false |  |
| install_log_path | The absolute path to the Install log | []string | `[/var/log/install.log]` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| offset_storage_dir | The directory that the offset storage file will be created | string | `$OIQ_OTEL_COLLECTOR_HOME/storage` | false |  |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/macos_logs.yaml
    parameters:
      enable_system_log: true
      system_log_path: [/var/log/system.log]
      enable_install_log: true
      install_log_path: [/var/log/install.log]
      start_at: end
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
