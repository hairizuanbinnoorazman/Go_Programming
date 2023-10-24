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
	docker build -t gcr.io/${gcp_project}/basic-app:${version} -f ./deployment/docker/slim.Dockerfile .

gcp_instance:
	gcloud compute instances create instance-1 --project=$(gcp_project) --zone=us-central1-a --machine-type=e2-medium

gke_cluster:
	gcloud container clusters create-auto "autopilot-cluster-1" --region "us-central1"