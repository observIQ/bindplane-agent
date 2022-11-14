# File Plugin

Log parser for generic files

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_path | Specify a single path or multiple paths to read one or many files. You may also use a wildcard (*) to read multiple files within a directory | []string |  | true |  |
| exclude_file_path | Specify a single path or multiple paths to exclude one or many files from being read. You may also use a wildcard (*) to exclude multiple files from being read within a directory | []string | `[]` | false |  |
| multiline_line_start_pattern | A Regex pattern that matches the start of a multiline log entry in the log file | string | `` | false |  |
| encoding | Specify the encoding of the file(s) being read. In most cases, you can leave the default option selected | string | `utf-8` | false | `nop`, `utf-8`, `utf-16le`, `utf-16be`, `ascii`, `big5` |
| parse_format | Format of parsed logs (`none`, `json`, or `regex`) | string | `none` | false | `none`, `json`, `regex` |
| regex_pattern | Pattern for regex parsed log | string | `` | false |  |
| log_type | Adds the specified 'Type' as a label to each log message | string | `file` | false |  |
| include_file_name | Whether to add the file name as the attribute log.file.name. | bool | `true` | false |  |
| include_file_path | Whether to add the file path as the attribute log.file.path. | bool | `false` | false |  |
| start_at | At startup, where to start reading logs from the file (`beginning` or `end`) | string | `end` | false | `beginning`, `end` |
| retain_raw_logs | When enabled will preserve the original log message on the body in a `raw_log` key | bool | `false` | false |  |
| offset_storage_dir | The directory that the offset storage file will be created | string | `$OIQ_OTEL_COLLECTOR_HOME/storage` | false |  |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/file_logs.yaml
    parameters:
      file_path: [$FILE_PATH]
      exclude_file_path: []
      encoding: utf-8
      parse_format: none
      log_type: file
      include_file_name: true
      include_file_path: false
      start_at: end
      retain_raw_logs: false
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
