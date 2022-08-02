# `vcenter` plugin

The `vcenter` plugin consumes [VMware vCenter](https://www.vmware.com/products/vcenter-server.html) log entries from the local filesystem and outputs parsed entries. 

## Configuration Fields

| Field | Type | Default | Description |
| --- | --- |--- | --- |
| `listen_port` | `int` | `5140` | A port which the agent will listen for syslog messages | 
| `listen_ip` | `string` | `"0.0.0.0"` | A ip address of the form `<ip>` | 
| `max_buffer_size` | `string` | `"1024kib"` | Maximum size of buffer that may be allocated while reading TCP input | 
| `enable_tls` | `bool` | `false` | Enable TLS for the TCP listener | 
| `certificate_file` | `string` | `"/opt/cert"` | File path for the X509 TLS certificate chain | 
| `private_key_file` | `string` | `"/opt/key"` | File path for the X509 TLS private key path | 

## Prerequisites
See the VMware documentation to enable Syslog on your vCenter instance [here](https://docs.vmware.com/en/VMware-vSphere/6.7/com.vmware.vsphere.vcsa.doc/GUID-9633A961-A5C3-4658-B099-B81E0512DC21.html).

## Example usage

### Configuration

Using default log paths:

```yaml
pipeline:
- type: vcenter
- type: stdout

```

With TLS enabled:

```yaml
pipeline:
- type: vcenter
  enable_tls: true
- type: stdout

```
