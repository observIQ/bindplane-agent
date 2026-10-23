# OpAMP Configuration

The Bindplane Agent can be setup as an agent that is managed by the [Bindplane platform](https://bindplane.com/) via OpAMP.

## Configuration

The agent can be configured to connect to the server a few different ways.

### Config File

The agent can be configured to read its connection config from a `manager.yaml` file. The `--manager` flag can be used to specify the location of this config file, by default it's `./manager.yaml`. The contents of the `manager.yaml` are detailed out in the table below.

| Parameter  | Required | Description                                                                |
| :--------  | :------: | :------------------------------------------------------------------------- |
| endpoint   | X        | The API endpoint to communicate with the server via websocket              |
| secret_key |          | The Secret Key defined for the server to be used for authorization         |
| agent_id   |          | A [ULID](https://github.com/ulid/spec) used to uniquely identify the agent |
| labels     |          | A comma separated list of labels in the form `label=value`                 |
| agent_name |          | Human readable name for the agent                                          |
| tls_config |          | See [tls config](#tls-config) section                                      |

Here's an example of what a common `manager.yaml` looks like:

```yaml
endpoint: ws://bindplane.localnet/v1/opamp
secret_key: 3d83f0cb-2567-42c7-ada6-960842924d11
agent_id: 01H5MG7N9N36J28WEC8A8X5B17
```

#### TLS Config

If TLS is enabled on the server the agent will need to be configured in order to connect. 

**Note**: If using TLS on the server the `endpoint` field will need to have the `wss` protocol for TLS enabled websockets.

| Parameter            | Required | Description                                                                                         |
| :------------------- | :------: | :-------------------------------------------------------------------------------------------------- |
| insecure_skip_verify |          | InsecureSkipVerify controls whether a client verifies the server's certificate chain and host name. |
| key_file             |          | Path to the `.key` file                                                                             |
| cert_file            |          | Path to the Certificate file                                                                        |
| ca_file              |          | Path to the Certificate Authority file                                                              |

### Environment variables

The agent can also use environment variables to set portions of the connection configuration. This is useful for a containerized agent where a mounted volume might not be present. 

If the agent can not find the specified `manager.yaml` file it will search for the environment variables and create a `manager.yaml` at the location of the `--manager` command argument.

**Note**: Only the `OPAMP_ENDPOINT` is required. If this is not set and there is no `manager.yaml` the agent will start in its normal standalone mode.

| Environment Variable  | Required | Description                                                                       |
| :-------------------- | :------: | :-------------------------------------------------------------------------------- |
| OPAMP_ENDPOINT        | X        | The API endpoint to communicate with the server via websocket                     |
| OPAMP_SECRET_KEY      |          | The Secret Key defined for the server to be used for authorization                |
| OPAMP_AGENT_ID        |          | A UUID used to uniquely identify the agent. If not supplied one will be generated |
| OPAMP_LABELS          |          | A comma separated list of labels in the form `label=value`                        |
| OPAMP_AGENT_NAME      |          | Human readable name for the agent                                                 |
| OPAMP_TLS_SKIP_VERIFY |          | Set to `"true"` to skip verification of the OpAMP server's TLS certificate        |
| OPAMP_TLS_CA          |          | File path to a certificate authority file that should be used to validate the server's TLS certificate |
| OPAMP_TLS_CERT        |          | File path to a certificate file that will be used for client TLS authentication |
| OPAMP_TLS_KEY         |          | File path to a private key file that will be used for client TLS authentication |

## Testing the Connection

The `test-connection` command checks the connection to the OpAMP server using the settings in `manager.yaml`. It performs a websocket handshake only and does not send any OpAMP messages. The command fails if `manager.yaml` does not exist and never modifies it.

On Linux:

```shell
sudo /opt/observiq-otel-collector/observiq-otel-collector test-connection
```

On Windows:

```shell
"C:\Program Files\observIQ OpenTelemetry Collector\observiq-otel-collector.exe" test-connection
```

| Flag      | Description                                                         |
| :-------- | :------------------------------------------------------------------ |
| --manager | Path to the manager config. See below for the default               |
| --timeout | Maximum time to wait for the handshake. Defaults to `30s`           |

When `--manager` is not set, the first of the following paths that exists is used. The selected path is printed.

1. `./manager.yaml`
2. `/opt/observiq-otel-collector/manager.yaml` (Linux)
3. `C:\Program Files\observIQ OpenTelemetry Collector\manager.yaml` (Windows)
4. The path in the `MANAGER_YAML_PATH` environment variable

TLS settings are read from `tls_config` in `manager.yaml`.

The command uses the same proxy environment variables as the agent (`HTTPS_PROXY`, `HTTP_PROXY`, and `NO_PROXY`) and reports which proxy, if any, applies to the endpoint:

- **Linux with systemd**: The proxy variables are read from the `observiq-otel-collector` systemd unit, including its overrides and environment files, so the test uses the same proxy as the service. Proxy variables in the current shell are ignored. The command must be run as root or with `sudo`.
- **Linux without systemd, such as containers**: The current environment is used.
- **Windows**: The proxy variables are read from the system environment and the `Environment` value of the `observiq-otel-collector` service (`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\observiq-otel-collector`), which overrides the system environment, so the test uses the same proxy as the service. Proxy variables in the current console are ignored. If the service is not installed, the current environment is used.

If Bindplane rejects the connection with `429 Too Many Requests`, the command retries up to three times. It waits for the duration in the server's `Retry-After` header, or backs off starting at one second when the header is not set, and prints each retry and how long it waits. Bindplane returns `429` when too many agents are connecting or when it detects another connection with the same agent ID, for example when the collector is running.
