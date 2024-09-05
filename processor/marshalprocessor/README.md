# Marshal Processor

This processor is used to marshal parsed logs into JSON, XML, or KV format.

This processor is intended to be wrapped into the Marshal processor in Bindplane.

NOTE: XML support is in progress and not yet available.

## Supported pipelines

- Logs

## How It Works

1. This processor expects its input to contain a parsed log body. It will marshal the body fields only; if additional fields from the log are desired, they must first be moved to the body.

2. The body can be marshaled to string-encoded JSON, XML, or KV.

   - For KV:
     - Fields will be converted to "key1=value1 key2=value2 key3=value3..." if no separators are configured
     - The parsed fields should be flattened first so that every key is at the top level.

3. The output of this processor will be the same as the input, but with a modified log body. Any body incompatible with the marshal type will be unchanged.

## Configuration

| Field           | Type   | Default | Description                                         |
| --------------- | ------ | ------- | --------------------------------------------------- |
| marshalTo       | string | ""      | The format to marshal into. Can be JSON, XML, or KV |
| kvSeparator     | rune   | "="     | The separator between key and value                 |
| kvPairSeparator | rune   | " "     | The separator between KV pairs                      |

## Example JSON Config

```yaml
processors:
  marshal:
    marshalTo: "JSON"
service:
  pipelines:
    logs:
      processors: [marshal]
```

## Example KV config

```yaml
processors:
  marshal:
    marshalTo: "KV"
    kvSeparator: ","
    kvPairSeparator: ":"
service:
  pipelines:
    logs:
      processors: [marshal]
```

## Example Output

The parsed body will be replaced by a marshaled body. All other fields are untouched.

### 1. Nested parsed body to JSON

In the example below, "bindplane-otel-attributes" represents attributes that have been moved to the body.

#### Parsed body

```
"body": {
    "kvlistValue": {
        "values": [
            { "key": "severity", "value": { "doubleValue": 155 } },
            {
                "key": "nested",
                "value": {
                    "kvlistValue": {
                        "values": [
                            { "key": "n2", "value": { "doubleValue": 2 } },
                            { "key": "n1", "value": { "doubleValue": 1 } }
                        ]
                    }
                }
            },
            { "key": "name", "value": { "stringValue": "test" } },
            {
                "key": "bindplane-otel-attributes",
                "value": {
                    "kvlistValue": {
                        "values": [
                            {
                                "key": "baba",
                                "value": { "stringValue": "you" }
                            },
                            {
                                "key": "host",
                                "value": { "stringValue": "myhost" }
                            }
                        ]
                    }
                }
            }
        ]
    }
},

```

#### JSON output

```
"body": {
    "stringValue": {
        {
            "bindplane-otel-attributes":
                {
                    "baba":"you",
                    "host":"myhost"
                },
            "name":"test",
            "nested":
                {
                    "n1":1,
                    "n2":2
                },
            "severity":155
        }
    }
}
```

### 2: Flattened parsed body to KV with default separators

In the example below, flattening has already been done on the "nested" field and the "bindplane-otel-attributes" field.

#### Parsed flattened body

```
"body": {
    "kvlistValue": {
        "values": [
            { "key": "severity", "value": { "doubleValue": 155 } },
            {
                "key": "nested-n1",
                "value": { "doubleValue": 1 }
            },
            {
                "key": "nested-n2",
                "value": { "doubleValue": 2 }
            },
            { "key": "name", "value": { "stringValue": "test" } },
            {
                "key": "bindplane-otel-attributes-baba",
                "value": { "stringValue": "you" }
            },
            {
                "key": "bindplane-otel-attributes-host",
                "value": { "stringValue": "myhost" }
            }
        ]
    }
}
```

#### KV output

```
"body": {
    "stringValue": {
        bindplane-otel-attributes-baba=you
        bindplane-otel-attributes-host=myhost
        name=test
        nested-n1=1
        nested-n2=2
        severity=155
    }
}
```
