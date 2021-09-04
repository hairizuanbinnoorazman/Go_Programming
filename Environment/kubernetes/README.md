# Setting up Kubernetes environment

First part is to create a cluster. In order to make a pretty much barebones cluster pretty easily in GCP, we can use the following command:

```bash
gcloud beta container clusters create "cluster-2" --zone "asia-southeast1-a" --no-enable-basic-auth --machine-type "e2-standard-4" --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --max-pods-per-node "110" --num-nodes "3" --enable-ip-alias --no-enable-intra-node-visibility --default-max-pods-per-node "110" --no-enable-master-authorized-networks --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --enable-shielded-nodes --node-locations "asia-southeast1-a" --logging=NONE --monitoring=NONE
```

To connect to the cluster

```bash
gcloud container clusters get-credentials cluster-2 --zone asia-southeast1-a
```

## Application Monitoring & Dashboards

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# helm show values prometheus-community/kube-prometheus-stack > customize.yaml

helm upgrade --install -f prom.yaml kube-prometheus-stack prometheus-community/kube-prometheus-stack

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

## S3 Storage

Let's assume a scenario we're building a private cloud outside of the standard Cloud Providers. What are the alternatives that we have?

The easiest example is the Minio Project. https://github.com/minio/operator/tree/master/helm/minio-operator

```bash
helm repo add minio https://operator.min.io/

helm upgrade --install -f minio.yaml minio-operator minio/minio-operator

# Access to "tenant" management
kubectl get secret $(kubectl get serviceaccount console-sa --namespace default -o jsonpath="{.secrets[0].name}") --namespace default -o jsonpath="{.data.token}" | base64 --decode 
kubectl  port-forward svc/console 9090:9090

# Access to a single tenant's s3
kubectl port-forward svc/minio1-console 9090:9090
```

To test this, we can run it in a separate nginx pod

```bash
kubectl create deployment nginx --image=nginx
kubectl exec -it <pod-name> -- /bin/bash
apt update && apt install -y iputils-ping

wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc
mv mc /usr/local/bin/mc

mc alias set yahoo http://minio1-hl.default.svc.cluster.local:9000 minio minio123
mc admin info yahoo
mc ls yahoo
mc mb yahoo/testtest
```

## Application Logging

One potential logging stack to use is Loki

This one install both loki and promtail as well. Promtail would be the component that would be deployed to extract logs from the pods which would be piped into loki

Refer to the values file here:  
https://github.com/grafana/helm-charts/blob/main/charts/loki-distributed/values.yaml

Refer to the configuration here:  
https://grafana.com/docs/loki/latest/configuration/

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

helm upgrade --install -f loki.yaml loki grafana/loki-distributed
```

We need to install the agent to sniff out the logging items

```bash
helm upgrade --install -f promtail.yaml promtail grafana/promtail
```

## Distributed Tracing

One potential tracing tool is tempo

This is the configuration file:  
https://grafana.com/docs/tempo/latest/configuration/#storage

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

helm upgrade --install -f tempo.yaml tempo grafana/tempo-distributed
```

## Cleanup

Delete everything that was installed

```bash
helm delete tempo
helm delete promtail
helm delete loki
helm delete minio-operator
helm delete kube-prometheus-stack
```

Delete volumes

```bash
kubectl delete pvc --all
```

To delete the cluster

```bash
gcloud container clusters delete cluster-2 --zone asia-southeast1-a
```

## Potential Tools

- etcd operator
- kafka operator
- k6 performance tooling?
- argocd
- jenkins
- airflow
- http://allure.qatools.ru/
- https://github.com/devopsprodigy/kubegraf
- https://github.com/armosec/kubescape
- vault-operator
- open-policy-agent
- keda
