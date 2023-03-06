# W3C Plugin

Log Parser for W3C

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_log_path | Specify a single path or multiple paths to read one or many files. You may also use a wildcard (*) to read multiple files within a directory. | []string |  | true |  |
| exclude_file_log_path | Specify a single path or multiple paths to exclude one or many files from being read. You may also use a wildcard (*) to exclude multiple files from being read within a directory. | []string | `[]` | false |  |
| encoding | Specify the encoding of the file(s) being read. In most cases, you can leave the default option selected. | string | `utf-8` | false | `utf-8`, `utf-16le`, `utf-16be`, `ascii`, `big5` |
| log_type | Adds the specified 'Type' as a label to each log message. | string | `w3c` | false |  |
| start_at | At startup, where to start reading logs from the file. Must be set to "beginning" if 'headers' is not specified. | string | `beginning` | false | `beginning`, `end` |
| max_concurrent_files | Max number of W3C files that will be open during a polling cycle | int | `512` | false |  |
| include_file_name | Include File Name as a label | bool | `true` | false |  |
| include_file_path | Include File Path as a label | bool | `false` | false |  |
| include_file_name_resolved | Same as include_file_name, however, if file name is a symlink, the underlying file's name will be set as a label | bool | `false` | false |  |
| include_file_path_resolved | Same as include_file_path, however, if file path is a symlink, the underlying file's path will be set as a label | bool | `false` | false |  |
| header | The W3C header which specifies the field names. Field names will be auto detected if unspecified. | string |  | false |  |
| delimiter | Delimiter character used between fields (Defaults to a tab character) | string | `	` | false |  |
| header_delimiter | Delimiter character used between header fields (Defaults to the value of 'delimiter') | string |  | false |  |
| offset_storage_dir | The directory that the offset storage file will be created | string | `$OIQ_OTEL_COLLECTOR_HOME/storage` | false |  |

## Example Config:

Below is an example of a basic config

```yaml
receivers:
  plugin:
    path: ./plugins/w3c_logs.yaml
    parameters:
      file_log_path: [$FILE_LOG_PATH]
      exclude_file_log_path: []
      encoding: utf-8
      log_type: w3c
      start_at: beginning
      max_concurrent_files: 512
      include_file_name: true
      include_file_path: false
      include_file_name_resolved: false
      include_file_path_resolved: false
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
