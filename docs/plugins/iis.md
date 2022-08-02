# `iis` plugin

The `iis` plugin consumes [Microsoft IIS](https://www.iis.net/) logs entries from the local filesystem and outputs parsed entries.

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- |--- | --- |
| `file_path` | `string` | `"C:/inetpub/logs/LogFiles/W3SVC*/**/*.log"` | The absolute path to the Microsoft IIS logs |
| `exclude_file_log_path` | `strings` | `[]` | Specify a single path or multiple paths to exclude one or many files from being read. You may also use a wildcard (*) to exclude multiple files from being read within a directory. |
| `location` | `enum` | `"UTC"` | The geographic location (timezone) to use when parsing the timestamp |
| `include_file_name` | `bool` | `true` | Include File Name as a label |
| `include_file_path` | `bool` | `false` | Include File Path as a label |
| `include_file_name_resolved` | `bool` | `false` | Same as include_file_name, however, if file name is a symlink, the underlying file's name will be set as a label |
| `include_file_path_resolved` | `bool` | `false` | Same as include_file_path, however, if file path is a symlink, the underlying file's path will be set as a label |
| `start_at` | `enum` | `end` | Start reading file from 'beginning' or 'end' | 

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: iis
- type: stdout

```

With non-standard log path:

```yaml
pipeline:
- type: iis
  file_path: "/path/to/logs"
- type: stdout

```
