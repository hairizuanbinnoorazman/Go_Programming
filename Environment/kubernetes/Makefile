cluster_name?=cluster-2

environment:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add minio https://operator.min.io/
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add nats https://nats-io.github.io/k8s/helm/charts/
	helm repo add jetstack https://charts.jetstack.io
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add codecentric https://codecentric.github.io/helm-charts
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm repo add kedacore https://kedacore.github.io/charts
	helm repo add apache-airflow https://airflow.apache.org
	helm repo add hashicorp https://helm.releases.hashicorp.com
	helm repo add argo https://argoproj.github.io/argo-helm
	helm repo add gatekeeper https://open-policy-agent.github.io/gatekeeper/charts
	helm repo add pyroscope-io https://pyroscope-io.github.io/helm-chart
	helm repo update

cluster:
	gcloud beta container clusters create $(cluster_name) --zone "asia-southeast1-a" --no-enable-basic-auth --machine-type "e2-standard-4" --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --max-pods-per-node "110" --num-nodes "3" --enable-ip-alias --no-enable-intra-node-visibility --default-max-pods-per-node "110" --no-enable-master-authorized-networks --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --enable-shielded-nodes --node-locations "asia-southeast1-a" --logging=NONE --monitoring=NONE
	gcloud container clusters get-credentials $(cluster_name) --zone asia-southeast1-a

observability:
	helm upgrade --install -f prom.yaml kube-prometheus-stack prometheus-community/kube-prometheus-stack
	helm upgrade --install -f minio.yaml minio-operator minio/minio-operator
	sleep 100
	make buckets
	helm upgrade --install -f loki.yaml loki grafana/loki-distributed
	helm upgrade --install -f promtail.yaml promtail grafana/promtail
	helm upgrade --install pyroscope pyroscope-io/pyroscope
	helm upgrade --install -f tempo.yaml tempo grafana/tempo-distributed

mailhog:
	helm upgrade --install mailhog codecentric/mailhog

nats:
	helm upgrade --install -f nats.yaml nats nats/nats

certmanager:
	helm upgrade --install cert-manager jetstack/cert-manager --set installCRDs=true

vault:
	helm upgrade --install vault hashicorp/vault

mysql:
	helm upgrade --install -f mysql.yaml mysql bitnami/mysql

nginx:
	helm upgrade --install -f ingress-nginx.yaml ingress-nginx ingress-nginx/ingress-nginx

keycloak:
	helm upgrade --install -f keycloak.yaml keycloak codecentric/keycloak

gatekeeper:
	helm install -n gatekeeper-system gatekeeper gatekeeper/gatekeeper --create-namespace

#airflow:
#	helm upgrade --install airflow apache-airflow/airflow

#keda:
#	helm upgrade --install keda.yaml kedacore/keda

delete:
	helm delete $$(helm list | grep deployed | awk '{ printf $$1" " }')
	kubectl delete pvc --all
	gcloud container clusters delete $(cluster_name) --quiet --zone asia-southeast1-a

buckets:
	kubectl apply -f make-bucket.yaml
	sleep 60
	sh check-bucket.sh make-bucket-testtest
	sh check-bucket.sh make-bucket-haha
	sh check-bucket.sh make-bucket-profiler
	kubectl delete -f make-bucket.yaml

access-minio:
	kubectl port-forward service/minio1-console 9090:9090

access-grafana:
	kubectl port-forward service/kube-prometheus-stack-grafana 3000:80

access-prometheus:
	kubectl port-forward service/kube-prometheus-stack-prometheus 9090:9090

access-tempo:
	kubectl port-forward service/tempo-tempo-distributed-query-frontend 16686:16686

access-keycloak:
	kubectl port-forward service/keycloak-http 8080:80
