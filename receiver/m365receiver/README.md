# Microsoft Office 365 Receiver

| Status |  |
| -------------------------- | ----------- |
| Stability | [dev] |
| Supported pipeline types | metrics, logs |
| Distributions | [contrib] |

Receives metrics from [Microsoft Office 365](https://www.microsoft365.com/)
via the [Microsoft Graph API](https://learn.microsoft.com/en-us/graph/api/overview?view=graph-rest-1.0&preserve-view=true),
and logs via the [Microsoft Management API](https://learn.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-schema).

## Getting Started
To monitor metrics and logs from Microsoft Office 365, some configuration is required for the receiver 
as well as the instance of Office 365 to be monitored. The guide below will outline how to configure both. 
It's recommended to begin with the Office 365 instance since the receiver needs parameters obtained from the instance.

### Configuring Office 365
The steps below outline how to configure Office 365 to allow the receiver to collect metrics from it. 
To use this receiver, the instance of Office 365 needs the following subscriptions: **Microsoft 365 Business Basic**, **Microsoft 365 E5 Compliance** or **Microsoft 365 E3**.

1. **Login to Azure:** Log in to Microsoft Azure under an Admin account for the instance of 365 to be monitored.
2. **Register the receiver in Azure AD:** Navigate to Azure Active Directory. Then go to "App Registrations" and select "New Registration". 
Give the app a descriptive name like "365 Receiver". For "Supported account types", select the Single Tenant option and leave the Redirect URL empty.
3. **Add API Permissions:** Select "View API Permissions" beneath the general application info and click "Add Permissions". The permissions needed for metrics and logs differ, so for whichever monitoring is needed the respective permissions are outlined below.
    - **Metrics:** Select "Microsoft Graph", then "Application Permissions". Find the "Reports" tab and select "Reports.Read.All". Click "Add Permissions" at the bottom of the panel.
    - **Logs:** Select "Office 365 Management APIs", then "Application Permissions". Now select the "ActivityFeed.Read", "ActivityFeed.ReadDlp", and "ServiceHealth.Read" permissions. Click "Add Permissions" at the bottom of the panel.
4. **Grant Admin Consent:** Select the "Grant admin consent for {organization}" button and confirm the pop-up. This will allow the application to access the data returned by the Microsoft Graph and Office 365 Management APIs.
5. **Generate Client Secret:** Select the "Certificates & secrets" tab in the left panel. Under the "Client Secrets" tab, select "New Client Secret." Give it a meaningful description and select the recommended period of 180 days. Save the text in the "Value" column since this is the only time that value will be accessible.
6. **Save Client ID and Tenant ID values:** You will also need the "client_id" value found on the information page for the application that was created. The value will be listed as "Application (client) id." You will also need the tenant value which will be listed as "Directory (tenant) id." Save these values for later.

After following the above steps, the instance of Microsoft Office 365 is ready for monitoring and the receiver can now be configured.

### Configuring the receiver
The Microsoft Office 365 receiver takes the following parameters. `tenant_id`, `client_id`, and `client_secret` are the only required values to receive metrics and logs. These values will have been retrieved while following the [Configuring Office 365](#configuring-office-365) guide. For metrics, the only parameters are the required ones already mentioned. All other parameters are optional ones for collecting logs. The default poll interval is 5 minutes and all logs are collected/true by default. 

- `tenant_id` : (required) identifies the instance of 365 to be monitored
- `client_id` : (required) the identifier this receiver will use to monitor the given tenant/instance
- `client_secret` : (required) the private key this receiver will use, must belong to the given client_id
- `logs`
    - `poll_interval` : time in minutes
    - `general` : true or false
    - `exchange` : true or false
    - `sharepoint` : true or false
    - `azureAD` : true of false
    - `dlp` : true or false
- `storage` : The component ID of a storage extension which can be used when polling for `logs` . The storage extension prevents duplication of data after a collector restart by remembering which data were previously collected.

**Example Configs**

Collect metrics:

```yaml
receivers:
  m365:
    tenant_id: tenant_id
    client_id: client_id
    client_secret: client_secret
exporters:

service:
  pipelines:
    metrics:
      receivers: [m365]
      exporters: 
```

Collect logs (default values):

```yaml
receivers:
  m365:
    tenant_id: tenant_id
    client_id: client_id
    client_secret: client_secret
exporters:

service:
  pipelines:
    logs:
      receivers: [m365]
      exporters: 
```

Collect logs (custom poll interval, storage component, only sharepoint & azureAD logs):
```yaml
receivers:
  m365:
    tenant_id: tenant_id
    client_id: client_id
    client_secret: client_secret
    logs:
      poll_interval: 10m
      general: false
      exchange: false
      dlp: false
    storage: file_storage
exporters:

service:
  pipelines:
    logs:
      receivers: [m365]
      exporters: 
```