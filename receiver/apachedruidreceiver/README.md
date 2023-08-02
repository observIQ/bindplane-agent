# Apache Druid Receiver
Receives metrics from [Apache Druid](https://druid.apache.org/)
via Apache Druid's [Metric Emitters](https://druid.apache.org/docs/latest/configuration/index.html#enabling-metrics).

## Minimum Agent Versions
- Introduced: [v1.32.0](https://github.com/observIQ/bindplane-agent/releases/tag/v1.32.0)

## Supported Pipelines
- Metrics

## How It Works
1. The user configures their instance of Apache Druid to enable the emission of metrics over HTTP or HTTPS.
2. The user configures this receiver in a pipeline.
3. The user configures a supported component to route telemetry from this receiver.

## Prerequisites
- Running instance of Apache Druid.

## Configuration
| Field               | Type     | Required | Description                                                                                                                                                             |
|---------------------|----------|------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| tls/cert_file       | string   | false    | The full path to a x509 PEM certificate file to be used for the TLS protocol when receiving metrics from the Druid instance.                                            |
| tls/key_file        | string   | false    | The full path to a x509 PEM private key file to be used for the TLS protocol when receiving metrics from the Druid instance.                                            |
| endpoint            | string   | true     | The endpoint on which the receiver will await POST requests from the Apache Druid emitter.                                                                              |
| basic_auth/username | string   | false    | The username expected to be used for basic authentication in the header of POST requests received from the Apache Druid emitter.                                        |
| basic_auth/password | string   | false    | The password expected to be used for basic authentication in the header of POST requests received from the Apache Druid emitter.                                        |

## Example Configurations

### Collect metrics: 
```yaml
receivers:
  apachedruid:
    metrics:
      tls:
        key_file: <full path to x509 PEM private key file>
        cert_file: <full path to x509 PEM certificate>
      endpoint: 0.0.0.0:8080
      basic_auth:
        username: john.doe
        password: 1234abcd
exporters:
  file/no_rotation:
    path: /some/file/path/foo.json
service:
  pipelines:
    metrics:
      receivers: [apachedruid]
      exporters: [file/no_rotation]
```

## How To
### Configuring Apache Druid
The steps below outline how to configure Apache Druid to allow the receiver to receive metrics from it. Step 1 is optional, if you would like to configure your Druid emitter to send metrics to a TLS-enabled receiver.

1. **Enable TLS (optional):** Receive a properly CA signed SSL certificate for use on the collector host.
2. **Configure Druid:** Alter your Druid configuration to enable [metric emission over HTTP](https://druid.apache.org/docs/latest/configuration/index.html#enabling-metrics).

**Configuration Notes:**
1. The receiver currently utilizes/has support for these Druid metrics:
    - query/count
    - query/success/count
    - query/failed/count
    - query/interrupted/count
    - query/timeout/count
    - sqlQuery/time
    - sqlQuery/bytes
2. Receiver-relevant configuration properties include:
    - Required configurations properties:
      - druid.emitter
        - must be druid.emitter=http
      - druid.monitoring.monitors
        - Recommended: druid.monitoring.monitors=["org.apache.druid.server.metrics.QueryCountStatsMonitor"]
      - druid.emitter.http.recipientBaseUrl
        - This is the endpoint on which the receiver will be listening for POST requests from the Druid emitter.
        - Example: http://localhost:8080
    - Optional configurations properties:
      - druid.emitter.http.basicAuthentication
        - Format: login:password
        - Conditionally optional, based on whether these fields were provided in the receiver configuration.
      - druid.enableTlsPort and its counterpart, druid.enablePlaintextPort
        - default to false and true, respectively
      - See "HTTP Emitter Module TLS Overrides" in the Druid documentation linked above for the necessary properties to alter if you performed optional step 1.
        - It is recommended that you check whether your Druid instance is already configured to use TLS for its internal HTTP client to possibly simplify the emitter TLS configuration process, if relevant.

After following the above steps, the instance of Apache Druid is ready for monitoring and the receiver can now be configured.
