# W3C Plugin

Log Parser for W3C

## Configuration Parameters

| Name | Description | Type | Default | Required | Values |
|:-- |:-- |:-- |:-- |:-- |:-- |
| file_log_path | Specify a single path or multiple paths to read one or many files. You may also use a wildcard (*) to read multiple files within a directory. | []string |  | true |  |
| exclude_file_log_path | Specify a single path or multiple paths to exclude one or many files from being read. You may also use a wildcard (*) to exclude multiple files from being read within a directory. | []string | `[]` | false |  |
| encoding | Specify the encoding of the file(s) being read. In most cases, you can leave the default option selected. | string | `utf-8` | false | `utf-8`, `utf-16le`, `utf-16be`, `ascii`, `big5` |
| log_type | Adds the specified 'Type' as a label to each log message. | string | `w3c` | false |  |
| start_at | At startup, where to start reading logs from the file. Must be set to "beginning" if 'header' is not specified or if 'delete_after_read' is being used. | string | `beginning` | false | `beginning`, `end` |
| max_concurrent_files | Max number of W3C files that will be open during a polling cycle | int | `512` | false |  |
| timestamp_layout | Optional timestamp layout which will parse a timestamp field | string |  | false |  |
| timestamp_parse_from | Field to parse the timestamp from, required if 'timestamp_layout' is set | string |  | false |  |
| timestamp_layout_type | Optional timestamp layout type for parsing the timestamp, suggested if 'timestamp_layout' is set | string | `strptime` | false | `strptime`, `gotime`, `epoch` |
| timestamp_location | The geographic location (timezone) to use when parsing a timestamp that does not include a timezone. The available locations depend on the local IANA Time Zone database. | string | `UTC` | false |  |
| parse_from | Where to parse the data from | string | `body` | false |  |
| parse_to | Where the data will parse to | string | `body` | false | `attributes`, `body` |
| lazy_quotes | If true, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field. Cannot be true if 'ignore_quotes' is true. | bool | `false` | false |  |
| ignore_quotes | If true, all quotes are ignored, and fields are simply split on the delimiter. Cannot be true if 'lazy_quotes' is true. | bool | `false` | false |  |
| delete_after_read | Will delete static log files once they are completely read. When set, 'start_at' must be set to beginning. | bool | `false` | false |  |
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
      timestamp_layout_type: strptime
      timestamp_location: UTC
      parse_from: body
      parse_to: body
      lazy_quotes: false
      ignore_quotes: false
      delete_after_read: false
      include_file_name: true
      include_file_path: false
      include_file_name_resolved: false
      include_file_path_resolved: false
      offset_storage_dir: $OIQ_OTEL_COLLECTOR_HOME/storage
```
