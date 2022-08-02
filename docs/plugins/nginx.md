# `nginx` plugin

The `nginx` plugin consumes [nginx](https://www.nginx.com/) log entries from the local filesystem and outputs parsed entries.

## Supported Platforms

- Linux
- Windows
- MacOS
- Kubernetes

## Configuration Fields

| Field               | Default                      | Description |                                                                                                                                                                                                  
| ---                 | ---                          | ---         |                                                                                                                                                                                                     
| `source`            | `file`                       | Specifies where the logs are located. When choosing the 'file' option, the agent reads logs from the log path(s) specified. When choosing the 'Kubernetes' option, logs are read from /var/log/containers based on the Pod and Container specified |
| `enable_access_log` | `true`                       | Enable access log collection |
| `access_log_path`   | `/var/log/nginx/access.log`  | Path to access log file      |
| `log_format`        | `default`                    | Specifies the format of the access log entries. When choosing 'default', the agent will expect the log entries to match the default nginx log format. When choosing 'observiq', the agent will expect the log entries to match an optimized JSON format that adheres to the observIQ specification. See the Nginx documentation for more information |
| `enable_error_log`  | `true`                       | Enable error log collection  |
| `error_log_path`    | `/var/log/nginx/error.log`   | Path to error log file       |
| `start_at`          | `beginning`                  | Start reading file from 'beginning' or 'end' |
| `encoding`          | `utf-8`                      | Specify the encoding of the file(s) being read. In most cases, you can leave the default option selected. |
| `cluster_name`      | `""`                         | Cluster name to be added as a resource label. Used when source is set to kubernetes |
| `pod_name`          | `nginx`                      | The pod name without the unique identifier on the end. It should match the deployment, daemonset, statefulset or other resource name. Used when source is set to kubernetes |
| `container_name`    | `*`                          | The container name, useful if the pod has more than one container. Used when source is set to kubernetes |

## Log Format

### Default Log Format

The default log format assumes the use of the combined format documented [here](http://nginx.org/en/docs/http/ngx_http_log_module.html).

Combined format configuration:
```
log_format combined '$remote_addr - $remote_user [$time_local] '
                    '"$request" $status $body_bytes_sent '
                    '"$http_referer" "$http_user_agent"';
```

Combined format sample log:
```
10.33.104.40 - - [11/Jan/2021:11:25:01 -0500] "GET / HTTP/1.1" 200 612 "-" "curl/7.58.0"
```

### observIQ Log format

The observIQ log format is an enhanced log format that includes many useful fields that do not exist in the default
logging format, such as upstream information and http_x_forwarded_for headers.

observIQ log format configuration:
```
log_format observiq '{"remote_addr":"$remote_addr","remote_user":"$remote_user","time_local":"$time_local","request":"$request","status":"$status","body_bytes_sent":"$body_bytes_sent","http_referer":"$http_referer","http_user_agent":"$http_user_agent","request_length":"$request_length","request_time":"$request_time","upstream_addr":"$upstream_addr","upstream_response_length":"$upstream_response_length","upstream_response_time":"$upstream_response_time","upstream_status":"$upstream_status","proxy_add_x_forwarded_for":"$proxy_add_x_forwarded_for","bytes_sent":"$bytes_sent","time_iso8601":"$time_iso8601","upstream_connect_time":"$upstream_connect_time","upstream_header_time":"$upstream_header_time","http_x_forwarded_for":"$http_x_forwarded_for"}';
```

observIQ log format sample log:
```
{"remote_addr":"10.33.104.40","remote_user":"-","time_local":"25/Feb/2021:16:20:01 -0500","request":"GET /about-us?app=prod&user=james&app=stage HTTP/1.1","status":"404","body_bytes_sent":"178","http_referer":"-","http_user_agent":"curl/7.58.0","request_length":"114","request_time":"0.000","upstream_addr":"-","upstream_response_length":"-","upstream_response_time":"-","upstream_status":"-","proxy_add_x_forwarded_for":"10.33.104.40","bytes_sent":"342","time_iso8601":"2021-02-25T16:20:01-05:00","upstream_connect_time":"-","upstream_header_time":"-","http_x_forwarded_for":"-"}
```

## Example usage
 
### Default Configuration File Source

Stanza Pipeline

```yaml
pipeline:
- type: nginx
- type: stdout

```

Input Entry (Access Log)

```
10.33.104.40 - - [11/Jan/2021:11:25:01 -0500] "GET / HTTP/1.1" 200 612 "-" "curl/7.58.0"
```

Output Entry (Access Log)

```json
{
  "timestamp": "2021-01-11T11:25:01-05:00",
  "severity": 30,
  "severity_text": "200",
  "labels": {
    "file_name": "access.log",
    "log_type": "nginx.access",
    "plugin_id": "nginx"
  },
  "record": {
    "body_bytes_sent": "612",
    "http_referer": "-",
    "http_user_agent": "curl/7.58.0",
    "method": "GET",
    "path": "/",
    "protocol": "HTTP",
    "protocol_version": "1.1",
    "remote_addr": "10.33.104.40",
    "remote_user": "-",
    "status": "200"
  }
}
```

Input Entry (Error Log)

```
2021/02/25 16:51:01 [emerg] 18747#18747: duplicate "log_format" name "oiq" in /root/nginx.conf.bad:43
```

Output Entry (Error Log)

```json
{
  "timestamp": "2021-02-25T16:51:01-05:00",
  "severity": 90,
  "severity_text": "emerg",
  "labels": {
    "file_name": "error.log",
    "log_type": "nginx.error",
    "plugin_id": "nginx"
  },
  "record": {
    "message": "duplicate \"log_format\" name \"oiq\" in /root/nginx.conf.bad:43",
    "pid": "18747",
    "tid": "18747"
  }
}
```

### observIQ Configuration File Source

Stanza Pipeline

```yaml
pipeline:
- type: nginx
  format: observiq
- type: stdout

```

Input Entry (Access Log)

```
{"remote_addr":"10.33.104.40","remote_user":"-","time_local":"25/Feb/2021:16:20:01 -0500","request":"GET /about-us?app=prod&user=james&app=stage HTTP/1.1","status":"404","body_bytes_sent":"178","http_referer":"-","http_user_agent":"curl/7.58.0","request_length":"114","request_time":"0.000","upstream_addr":"-","upstream_response_length":"-","upstream_response_time":"-","upstream_status":"-","proxy_add_x_forwarded_for":"10.33.104.40","bytes_sent":"342","time_iso8601":"2021-02-25T16:20:01-05:00","upstream_connect_time":"-","upstream_header_time":"-","http_x_forwarded_for":"-"}
```

Output Entry (Access Log)

```json
{
  "timestamp": "2021-02-25T16:20:01-05:00",
  "severity": 50,
  "severity_text": "404",
  "labels": {
    "file_name": "access.log",
    "log_type": "nginx.access",
    "plugin_id": "nginx"
  },
  "record": {
    "body_bytes_sent": "178",
    "bytes_sent": "342",
    "http_referer": "-",
    "http_user_agent": "curl/7.58.0",
    "http_x_forwarded_for": "-",
    "method": "GET",
    "path": "/about-us?app=prod&user=james&app=stage",
    "protocol": "HTTP",
    "protocol_version": "1.1",
    "proxy_add_x_forwarded_for": "10.33.104.40",
    "remote_addr": "10.33.104.40",
    "remote_user": "-",
    "request": "GET /about-us?app=prod&user=james&app=stage HTTP/1.1",
    "request_length": "114",
    "request_time": "0.000",
    "status": "404",
    "time_iso8601": "2021-02-25T16:20:01-05:00",
    "upstream_addr": "-",
    "upstream_connect_time": "-",
    "upstream_header_time": "-",
    "upstream_response_length": "-",
    "upstream_response_time": "-",
    "upstream_status": "-"
  }
}
```

### Default Configuration Kubernetes Source

Stanza can collect Nginx logs while running on Kubernetes. Use the provided script to deploy a sample environment:
1. Start [Minikube](https://minikube.sigs.k8s.io/docs/start/)
2. Deploy and expose Nginx
3. Create plugin configmap with nginx and kubernetes container plugins
4. Deploy Stanza as a daemonset
5. Detect the ip and port used to expose Nginx
6. Curl the Nginx endpoint to generate a log
7. Get Stanza's output

<details>
  <summary>click to expand `deploy.sh`</summary>

```bash
minikube start

kubectl apply -f https://k8s.io/examples/application/deployment.yaml
sleep 1
kubectl rollout status deploy/nginx-deployment

kubectl expose deployment nginx-deployment --port=80 --type=NodePort

kubectl create configmap plugin \
    --from-file plugins/nginx.yaml \
    --from-file plugins/kubernetes_container.yaml

cat <<EOF | kubectl create -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: stanza-config
data:
  config.yaml: |
    pipeline:
    - type: nginx 
      source: kubernetes 
      pod: nginx-deployment 
    - type: stdout 
EOF

kubectl apply -f https://raw.githubusercontent.com/observIQ/stanza-plugins/v0.0.88/dev/k8s/daemonset.yaml
sleep 1
kubectl rollout status ds/stanza

TARGET_IP=$(minikube ip)
TARGET_PORT=$(kubectl get svc nginx-deployment -o json | jq '.spec.ports[0].nodePort')

curl "${TARGET_IP}:${TARGET_PORT}"

kubectl logs ds/stanza
```
</details>

Stanza Pipeline

```yaml
pipeline:
- type: nginx
  source: kubernetes
  pod: nginx-deployment
- type: stdout

```

Input Entry (Access Log)

```
172.17.0.1 - - [08/Dec/2021:17:40:04 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.79.1" "-"
```

Output Entry (Access Log)

```json
{
  "timestamp": "2021-12-08T17:40:20Z",
  "severity": 30,
  "severity_text": "200",
  "labels": {
    "k8s-ns/kubernetes.io/metadata.name": "default",
    "k8s-pod/app": "nginx",
    "k8s-pod/pod-template-hash": "66b6c48dd5",
    "log_type": "nginx.access",
    "plugin_id": "nginx",
    "stream": "stdout"
  },
  "resource": {
    "k8s.cluster.name": "",
    "k8s.container.id": "ed7fc720c68357e0eb60d457e7483d47665eec3ac8962370a865fd810778237e",
    "k8s.container.name": "nginx",
    "k8s.deployment.name": "nginx-deployment",
    "k8s.namespace.name": "default",
    "k8s.namespace.uid": "bb7abccb-368b-4921-b577-c01a1379adc0",
    "k8s.node.name": "",
    "k8s.pod.name": "nginx-deployment-66b6c48dd5-pxrwz",
    "k8s.pod.uid": "11cdc1ea-4ab4-48f0-81e1-83de404cb639",
    "k8s.replicaset.name": "nginx-deployment-66b6c48dd5"
  },
  "record": {
    "body_bytes_sent": "612",
    "http_referer": "-",
    "http_user_agent": "curl/7.79.1",
    "method": "GET",
    "path": "/",
    "protocol": "HTTP",
    "protocol_version": "1.1",
    "remote_addr": "172.17.0.1",
    "remote_user": "-",
    "status": "200"
  }
}
```