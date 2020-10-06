# Setting up Kubernetes environment

## Application Monitoring & Dashboards

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# helm show values prometheus-community/kube-prometheus-stack > customize.yaml

helm upgrade --install -f customize.yaml kube-prometheus-stack prometheus-community/kube-prometheus-stack

export POD_NAME=$(kubectl get pods --namespace default -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace default port-forward $POD_NAME 9090
```

## Scaling based on custom metrics

```
helm upgrade --install -f adapter.yaml prometheus-adapter prometheus-community/prometheus-adapter
```

Check to see the metrics are received on the apiserver so that it can propagate it down to pods

```
kubectl proxy

# Then access the following endpoint
http://localhost:8001/apis/external.metrics.k8s.io/v1beta1
```

## Application Logging
