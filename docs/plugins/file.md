# `file` plugin

The `file` plugin consumes log files from the local filesystem and outputs entries.

## Supported Platforms

- Linux
- Windows
- MacOS
- Kubernetes

## Configuration Fields

| Field | Default | Description |
| --- | --- | --- |
| `file_log_path` | `''` | Specify a single path or multiple paths to read one or many files. You may also use a wildcard (*) to read multiple files within a directory. |
| `exclude_file_log_path` | `''` | Specify a single path or multiple paths to exclude one or many files from being read. You may also use a wildcard (*) to exclude multiple files from being read within a directory. |
| `enable_multiline` | `false` | Enable to parse Multiline Log Files |
| `multiline_line_start_pattern` | `''` | A Regex pattern that matches the start of a multiline log entry in the log file. |
| `encoding` | `utf-8` | Specify the encoding of the file(s) being read. In most cases, you can leave the default option selected. |
| `log_type` | `file` | Adds the specified 'Type' as a label to each log message. |
| `start_at` | `beginning` | Start reading file from 'beginning' or 'end' |

## Example usage

### Simple file input

Stanza Pipeline
```yaml
pipeline:
- type: file
  file_log_path:
    - "/test.log"
- type: stdout

```

<table>
<tr><td> `./test.log` </td> <td> Output records </td></tr>
<tr>
<td>

```
log1
log2
log3
```

</td>
<td>

```json
{
  "message": "log1"
},
{
  "message": "log2"
},
{
  "message": "log3"
}
```

</td>
</tr>
</table>

### Multiline

Configuration:
```yaml
pipeline:
- type: file
  file_log_path: "/test.log"
  enable_multiline: true
  multiline_line_start_pattern: 'START '
- type: stdout
```

<table>
<tr><td> `./test.log` </td> <td> Output records </td></tr>
<tr>
<td>

```
START log1
log2
START log3
log4
```

</td>
<td>

```json
{
  "message": "START log1\nlog2\n"
},
{
  "message": "START log3\nlog4\n"
}
```

</td>
</tr>
</table>
