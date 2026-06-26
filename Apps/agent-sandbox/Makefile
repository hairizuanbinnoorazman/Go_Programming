PROJECT_ID ?= new-demo-project-462517
REGION ?= asia-southeast1
ZONE ?= asia-southeast1-a
CLUSTER ?= agent-sandbox-standard
NAMESPACE ?= agent-sandbox-snapshots
AR_REPOSITORY ?= demo
IMAGE_TAG ?= $(shell git rev-parse --short HEAD 2>/dev/null || date +%Y%m%d%H%M%S)
IMAGE ?= $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(AR_REPOSITORY)/agent-sandbox-demo:$(IMAGE_TAG)
KUBECTL = kubectl --token="$$(gcloud auth print-access-token)"

.PHONY: build test image deploy apply rollout port-forward undeploy cluster-credentials

build:
	GOWORK=off go build ./...

test:
	GOWORK=off go test ./...

image:
	gcloud builds submit --project="$(PROJECT_ID)" --tag="$(IMAGE)" .

cluster-credentials:
	gcloud container clusters get-credentials "$(CLUSTER)" --project="$(PROJECT_ID)" --zone="$(ZONE)"

apply:
	sed "s/namespace: agent-sandbox-snapshots/namespace: $(NAMESPACE)/g" k8s/deployment.yaml | $(KUBECTL) apply -f -
	$(KUBECTL) -n "$(NAMESPACE)" set image deployment/agent-sandbox-demo app="$(IMAGE)"

rollout:
	$(KUBECTL) -n "$(NAMESPACE)" rollout status deployment/agent-sandbox-demo --timeout=5m

deploy: test image cluster-credentials apply rollout
	@echo "Deployed $(IMAGE)"
	@echo "Run: make port-forward"

port-forward:
	$(KUBECTL) -n "$(NAMESPACE)" port-forward service/agent-sandbox-demo 8080:8080

undeploy:
	sed "s/namespace: agent-sandbox-snapshots/namespace: $(NAMESPACE)/g" k8s/deployment.yaml | $(KUBECTL) delete -f -
