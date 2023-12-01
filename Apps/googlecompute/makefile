version?=v0
gcpproject?=$(shell gcloud config get project)
clustername?=compute-demo

docker-build:
	docker build -t gcr.io/$(gcpproject)/googlecompute:$(version) .

docker-push:
	docker push gcr.io/$(gcpproject)/googlecompute:$(version)

create-k8s:
	gcloud container clusters create-auto $(clustername)

connect-k8s:
	gcloud container clusters get-credentials $(clustername)

destroy-k8s:
	gcloud container clusters --quiet delete $(clustername)

deploy:
	echo $(shell sed "s/XXX/$(gcpproject)/g" deploy.yaml > correcteddeploy.yaml  )
	kubectl apply -f correcteddeploy.yaml
	rm correcteddeploy.yaml

create-secret:
	kubectl create secret generic cred-file  --from-file=credfile=key.json

deploy-withcred:
	echo $(shell sed "s/XXX/$(gcpproject)/g" deploy_withcredfile.yaml > correcteddeploy.yaml  )
	kubectl apply -f correcteddeploy.yaml
	rm correcteddeploy.yaml

run-app:
	SERVICE_ACCOUNT_FILE=key.json GCP_PROJECT_ID=$(gcpproject) REGION=asia-southeast1-a ./app

create-service-acct:
	gcloud iam service-accounts create sample-google-compute

delete-service-acct:
	gcloud iam service-accounts delete --quiet sample-google-compute@$(gcpproject).iam.gserviceaccount.com

create-service-acct-key:
	gcloud iam service-accounts keys create key.json \
		--iam-account=sample-google-compute@$(gcpproject).iam.gserviceaccount.com \
		--key-file-type=json

set-iam:	
	gcloud projects add-iam-policy-binding $(gcpproject) \
		--member="serviceAccount:sample-google-compute@$(gcpproject).iam.gserviceaccount.com" \
		--role="roles/editor"
	gcloud iam service-accounts add-iam-policy-binding \
		sample-google-compute@$(gcpproject).iam.gserviceaccount.com \
		--member="serviceAccount:$(gcpproject).svc.id.goog[default/default]" \
		--role="roles/iam.workloadIdentityUser"

edit-namespace:
	kubectl annotate serviceaccount default \
		--namespace=default \
		iam.gke.io/gcp-service-account=sample-google-compute@$(gcpproject).iam.gserviceaccount.com