VERSION?=0.0.4
GCP_PROJECT?=$$(gcloud config list project --format yaml  | yq e .core.project -)

build:
	docker build -t gcr.io/${GCP_PROJECT}/full-observability:${VERSION} .
push:
	docker push gcr.io/${GCP_PROJECT}/full-observability:${VERSION}
deploy:
	cd deployment && kustomize edit set image full-observability=gcr.io/${GCP_PROJECT}/full-observability:${VERSION}
	kustomize build ./deployment | kubectl apply -f -
full: build push deploy
req:
	kubectl exec $$(kubectl get pods | grep app- | awk '{ printf $1 }') -- curl localhost:8080
req2:
	kubectl exec $$(kubectl get pods | grep app2- | awk '{ printf $1 }') -- curl localhost:8080
req3:
	kubectl exec $$(kubectl get pods | grep app3- | awk '{ printf $1 }') -- curl localhost:8080