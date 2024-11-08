# Topology Processor

The topology processor is used in custom distributions of the collector to provide snapshot functionality in BindPlane. It is not currently included in the official `bindplane-agent`.

# TODO: write this whole file for topo processor

## Supported pipelines

- Logs
- Metrics
- Traces

## How it works

## Configuration

| Field   | Type   | Default | Required | Description                                                            |
|---------|--------|---------|----------|------------------------------------------------------------------------|
| opamp   | string | `opamp` | `true`   | Specifies the name of the opamp extension for sending custom messages. |


## Examples

### Usage in pipelines

```yaml
receivers:
  filelog:
    include: [/var/log/logfile.txt]

processors:
  topology:
    opamp: opamp

exporters:
  nop:

extensions:
  bindplane:
    labels: "labelA=valueA,labelB=valueB"
  opamp:
    endpoint: "https://localhost:3001/v1/opamp"

service:
  extensions: [bindplane, opamp]
  pipelines:
    logs:
      receivers: [filelog]
      processors: [topologyprocessor]
      exporters: [nop]
```
