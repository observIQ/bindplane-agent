# Google Cloud Exporter Prerequisites

The Google Cloud Exporter supports Metrics, Logs, and Traces.

## Google Cloud APIs

Enable the following APIs.
- Cloud Metrics
- Cloud Logging
- Cloud Tracing

To learn more about enabling APIs, see the [documentation](https://cloud.google.com/endpoints/docs/openapi/enable-api).

## Service Account

### GCE

If running within Google Cloud, ensure your system has the Cloud Monitoring, Cloud Logging, and Cloud Trace [scopes](https://developers.google.com/identity/protocols/oauth2/scopes#monitoring). This is all that is required to send metrics from the instance to Cloud Monitoring, you can skip the rest of this step.

### On Premise / Non Google Cloud

If running outside of Google Cloud (On prem, AWS, etc) or without the Cloud Monitoring scope, the Google Exporter requires a service account.

[Create a service account](https://cloud.google.com/iam/docs/creating-managing-service-accounts) with the following roles:
- Metrics: `roles/monitoring.metricWriter`
- Logs: `roles/logging.logWriter`
- Traces: `roles/cloudtrace.agent`

[Create a service account json key](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) and place it on the system that is running the collector. 

**Linux**

In this example, the key is placed at `/opt/observiq-otel-collector/sa.json` and it's permissions are restricted to the user running the collector process.
```bash
sudo cp sa.json /opt/observiq-otel-collector/sa.json
sudo chown observiq-otel-collector: /opt/observiq-otel-collector/sa.json
sudo chmod 0400 /opt/observiq-otel-collector/sa.json
```

Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable by creating a systemd override. A systemd override allows users to modify the systemd service configuration without modifying the service directly. This allows package upgrades to happen seemlessly. You can learn more about systemd units and overrides [here](https://wiki.archlinux.org/title/systemd#Replacement_unit_files).

Run the following command
```bash
sudo systemctl edit observiq-otel-collector
```

If this is the first time an override is being created, paste the following contents into the file:
```
[Service]
Environment=GOOGLE_APPLICATION_CREDENTIALS=/opt/observiq-otel-collector/sa.json
```

If an override is already in place, simply insert the `Environment` parameter into the existing `Service` section.

Restart the collector
```
sudo systemctl restart observiq-otel-collector
```

**Windows**

In this example, the key is placed at `C:/observiq/collector/sa.json`.

Set the `GOOGLE_APPLICATION_CREDENTIALS` with the command prompt `setx` command.

Run the following command
```
setx GOOGLE_APPLICATION_CREDENTIALS "C:/observiq/collector/sa.json"
```

Restart the service using the `services` application.

