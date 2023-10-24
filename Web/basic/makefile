version?=v1
gcp_project=$(shell gcloud config get project)

upload:
	ifndef APP_BUCKET
	$(error APP_BUCKET is not set)
	endif
	GOOS=linux GOARCH=amd64 go build -o web-basic-${version} .
	gsutil cp web-basic-${version} gs://${APP_BUCKET}/web-basic/web-basic-${version}
	rm web-basic-${version}

gcp_docker_build:
	docker build -t gcr.io/${gcp_project}/basic-app:${version} -f ./deployment/docker/slim.Dockerfile .

gcp_docker_push:
	docker push gcr.io/${gcp_project}/basic-app:${version}

gcp_instance:
	gcloud compute instances create instance-1 --project=$(gcp_project) --zone=us-central1-a --machine-type=e2-medium

gke_cluster:
	gcloud container clusters create-auto "autopilot-cluster-1" --region "us-central1"

create_deployment:
	kubectl create deployment basic-app --image=gcr.io/$(gcp_project)/basic-app:v1

tf_create_gke:
	TF_VAR_gcp_project_id=$(gcp_project) terraform plan -target=module.gke -var-file=local.tfvars -out gke.plan

tf_destroy_gke:
	TF_VAR_gcp_project_id=$(gcp_project) terraform plan -target=module.gke -destroy -var-file=local.tfvars -out gke.plan