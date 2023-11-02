# Kubernetes Telemetry with Google Cloud

The agent is capable of sending Kubernetes telemetry to Google Cloud Monitoring.

## Deployment

Deployments are cloud agnostic, meaning the agent(s) can run on most Kubernetes clusters, including development systems such as Minikube.

All metrics will show up under the following [Monitored resource types](https://cloud.google.com/monitoring/api/resources)
- [k8s_cluster](https://cloud.google.com/monitoring/api/resources)
- [k8s_node](https://cloud.google.com/monitoring/api/resources#tag_k8s_node)
- [k8s_pod](https://cloud.google.com/monitoring/api/resources#tag_k8s_pod)
- [k8s_container](https://cloud.google.com/monitoring/api/resources#tag_k8s_container)
