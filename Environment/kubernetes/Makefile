environment:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add minio https://operator.min.io/
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update

cluster:
	gcloud beta container clusters create "cluster-2" --zone "asia-southeast1-a" --no-enable-basic-auth --machine-type "e2-standard-4" --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --max-pods-per-node "110" --num-nodes "3" --enable-ip-alias --no-enable-intra-node-visibility --default-max-pods-per-node "110" --no-enable-master-authorized-networks --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --enable-shielded-nodes --node-locations "asia-southeast1-a" --logging=NONE --monitoring=NONE
	gcloud container clusters get-credentials cluster-2 --zone asia-southeast1-a

observability:
	helm upgrade --install -f prom.yaml kube-prometheus-stack prometheus-community/kube-prometheus-stack
	helm upgrade --install -f minio.yaml minio-operator minio/minio-operator
	sleep 30
	make buckets
	helm upgrade --install -f loki.yaml loki grafana/loki-distributed
	helm upgrade --install -f promtail.yaml promtail grafana/promtail
	helm upgrade --install -f tempo.yaml tempo grafana/tempo-distributed

delete:
	helm delete tempo
	helm delete promtail
	helm delete loki
	helm delete minio-operator
	helm delete kube-prometheus-stack
	kubectl delete pvc --all
	gcloud container clusters delete cluster-2 --quiet --zone asia-southeast1-a

buckets:
	kubectl apply -f make-bucket.yaml
	sleep 20
	sh check-bucket.sh make-bucket-testtest
	sh check-bucket.sh make-bucket-haha
	kubectl delete -f make-bucket.yaml

access-grafana:
	kubectl port-forward service/kube-prometheus-stack-grafana 3000:80
