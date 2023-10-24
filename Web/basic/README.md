# Basic Golang Application + Deployments

Contains basic example in order to deal with control a Golang application via Environment variables. The following repository serves to provide the basic application to deploy the same app into various different environments and hopefully have 0 issues between them.

Here is the full list of environment supported for this:

- [Basic Golang Application + Deployments](#basic-golang-application--deployments)
  - [Deploy via Bash/Makefiles](#deploy-via-bashmakefiles)
    - [Golang App on local workstation](#golang-app-on-local-workstation)
    - [Docker](#docker)
    - [Google Compute Engine](#google-compute-engine)
    - [Google Kubernetes Engine](#google-kubernetes-engine)
    - [Google Cloud Functions](#google-cloud-functions)
    - [Google Cloud Run](#google-cloud-run)
    - [Google App Engine](#google-app-engine)
  - [Deploy via Terraform](#deploy-via-terraform)
    - [Docker](#docker-1)
    - [Google Compute Engine](#google-compute-engine-1)
    - [Google Kubernetes Engine](#google-kubernetes-engine-1)


## Deploy via Bash/Makefiles

### Golang App on local workstation

Build the application and then run it

```bash
CGO_ENABLED=0 go build -o app main.go
```

Then, we can simply just run it after that

```bash
./app
```

### Docker

Build the container image (there are multiple versions of it) and then run it

```bash
docker build -t lol -f ./deployment/docker/alpine.Dockerfile .
docker run -p 8080:8080 lol
```

### Google Compute Engine

First build the linux binary for the application

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .
```

Then, we would need to create the Google Compute Engine instance

```bash
gcloud compute instances create instance-1 --project=<project id> --zone=us-central1-a --machine-type=e2-medium
```

Then, scp the binary into the server

```bash
scp app hairizuan@<ip hostname>:/usr/local/bin/app
ssh hairizuan@<ip hostname>
chmod +x /usr/local/bin/app
app
```

### Google Kubernetes Engine

First build the container with the tag that allows us to push it into container registry/artifact registry.

For container registry:

```bash
docker build -t gcr.io/<project id>/basic:v1 -f ./deployments/docker/slim.Dockerfile .
docker push gcr.io/<project id>/basic:v1
```

Then, we would need to run the following command to create a Auto-GKE cluster

```bash
gcloud container clusters create-auto "autopilot-cluster-1" --region "us-central1"

# Get access to GKE cluster
gcloud container clusters get-credentials autopilot-cluster-1 --zone us-central1-a
```

Next, we would run the following command to create the deployment service

```bash
kubectl create deployment basic-app --image=gcr.io/<project id>/basic:v1
```

### Google Cloud Functions

TODO

### Google Cloud Run

First build the container with the tag that allows us to push it into container registry/artifact registry.

For container registry:

```bash
docker build -t gcr.io/<project id>/basic-app:v1 -f ./deployments/docker/slim.Dockerfile .
docker push gcr.io/<project id>/basic-app:v1
```

Run gcloud command to create Cloud Run service

```bash
gcloud run deploy basic-app --image=gcr.io/<project id>/basic-app:v1 --concurrency=10 --max-instances=1 --platform=managed --allow-unauthenticated --ingress=all --cpu=1 --memory=500Mi --region=us-east1
```

### Google App Engine

TODO

## Deploy via Terraform

Some drawbacks (comparing it to other tools like Ansible).  
- Unable to simply pass variables from tfvars down to modules easily. Design conflict with initial aim of terraform. Compare the experience of this compared to Ansible  
  https://github.com/hashicorp/terraform/issues/32508

### Docker

NOTE: Seems hard to also include build step in terraform (need to replicate the dockerfile into terraform itself - it seems to have issues with relative paths? Might require a bit of debugging)

From terraform folder, run the following command

```bash
terraform plan -target=module.localdocker -out localdocker.plan
terraform apply localdocker.plan
```

To destroy it

```bash
terraform plan -target=module.localdocker -destroy -out localdocker.plan
terraform apply localdocker.plan
```

### Google Compute Engine

Require step to build Google Cloud Image (which we would then use to deploy the VM). This is done by relying on packer.

### Google Kubernetes Engine

It is assumed that the images have already been pushed into a container registry - all we have to do is to deploy it into a cluster.
