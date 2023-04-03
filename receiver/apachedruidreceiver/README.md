# Apache Druid Receiver

| Status                   |           |
|--------------------------|-----------|
| Stability                | [development] |
| Supported pipeline types | metrics       |


This Apache Druid receiver allows Apache Druid's [Metric Emitters](https://druid.apache.org/docs/latest/configuration/index.html#enabling-metrics) to send metrics over HTTP or HTTPS from the Apache Druid database to an OpenTelemetry collector.

## Getting Started

To successfully operate this receiver, you must follow these steps in order:
1. Have an Apache Druid database instance.
Steps 2 and 3 are optional, if you would like to configure your Druid emitter to send metrics to a TLS-enabled receiver.
2. Receive a properly CA signed SSL certificate for use on the collector host.
3. Configure the receiver using the previously acquired SSL certificate, and then start the collector.
4. Alter your Druid configuration to enable metric emission over HTTP: https://druid.apache.org/docs/latest/configuration/index.html#enabling-metrics
  - Currently supported processes include Broker and Historical processes.
  - Receiver-relevant configuration properties include:
    - Required configurations include:
      - druid.emitter
        - must be druid.emitter=http
      - druid.monitoring.monitors
        - must be druid.monitoring.monitors=["org.apache.druid.server.metrics.QueryCountStatsMonitor"]
      - druid.emitter.http.recipientBaseUrl
        - This is the endpoint on which the receiver will be listening for POST requests from the Druid emitter.
        - Example: http://localhost:8080
    - Optional configurations include:
      - druid.emitter.http.basicAuthentication
        - Format: login:password
        - Conditionally optional, based on whether these fields were provided in the receiver configuration.
      - druid.enableTlsPort and its counterpart, druid.enablePlaintextPort
        - default to false and true, respectively
      - See "HTTP Emitter Module TLS Overrides" in the Druid documentation linked above for the necessary properties to alter if you performed optional steps 2 and 3.
        - It is recommended that you check whether your Druid instance is already configured to use TLS for its internal HTTP client to possibly simplify the emitter TLS configuration process, if relevant.
5. Restart your Druid instance to update its config.

## Configuration

- `tls` (optional)
    - `cert_file` 
    - `key_file`
- `endpoint` 
  - The endpoint on which the receiver will await POST requests from the Apache Druid emitter.
- `basic_auth` (optional)
    - `username`
    - `password`
  - If these values are set, the receiver expects to see it in any valid requests within the header? // TODO IAN


### Example:

```yaml
receivers:
  apachedruid:
    metrics:
      tls:
        key_file: some_key_file
        cert_file: some_cert_file
      endpoint: 0.0.0.0:12345
      basic_auth:
        username: john.doe
        password: 1234abcd
```