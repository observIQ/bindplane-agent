# GitHub Receiver

Receives logs from [GitHub](https://github.com/)
via the [GitHub API](https://docs.github.com/en/rest?apiVersion=2022-11-28).

## Minimum Agent Versions

- Introduced: [v1.60.0](https://github.com/observIQ/bindplane-agent/releases/tag/v1.60.0)

## Supported Pipelines

- Metrics
- Logs

## How It Works

1. The user configures their instance of GitHub Enterprise to enable monitoring of audit logs.
2. The user configures this receiver in a pipeline.
3. The user configures a supported component to route telemetry from this receiver.

## Prerequisites

- Created instance of GitHub with the following subscriptions: GitHub Enterprise Cloud
- Access to an admin account for any enterprise, organization, or repo required for audit logs.

## Configuration

| Field          | Type                | Default | Required | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| -------------- | ------------------- | ------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| access_token   | configopaque.String |         | `true`   | Access token is required for audit log generation. Grants access to enterprise and organization if access token shows admin role. https://docs.github.com/en/enterprise-cloud@latest/apps/creating-github-apps/(authenticating-with-a-github-app/generating-a-user-access-token-for-a-github-app)                                                                                                                                                                                                                                                                                                                                                                          |
| log_type       | string              |         | `true`   | Specifies user, organization, or enterprise logs.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| name           | string              |         | `true`   | The name of the user, organization or enterprise.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| poll_interval  | time.Duration       |         | `false`  | The rate at which the receiver will poll for logs. An alternative to webhooks. The minimum poll interval calculated from information at [Rate limits for the REST API](https://docs.github.com/en/enterprise-cloud@latest/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28#primary-rate-limit-for-authenticated-users). The maximum rate limit for authenticated users is 5,000 requests per hour. This makes the minimum polling interval about 0.73 seconds. The default polling interval is 60 seconds, in order to have polling that is frequent enough to capture new logs without overwhelming the API or hitting rate limits too quickly. |
| webhook_config | WebhookConfig       |         | `false`  | Webhooks (not configured yet) that are used when an event triggers on an enterprise, organization, or user. An alternative to polling.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |

## Example Configuration

```yaml
receivers:
  github:
    access_token: access_token
    log_type: log_type
    name: name
    poll_interval: poll_interval
    webhook_config: webhook_config
exporters:
  googlecloud:
    project: my-gcp-project
service:
  pipelines:
    logs:
      receivers: [github]
      exporters: [googlecloud]
```
