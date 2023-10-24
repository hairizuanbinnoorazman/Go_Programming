# Basic Golang Application + Deployments

Contains basic example in order to deal with control a Golang application via Environment variables. The following repository serves to provide the basic application to deploy the same app into various different environments and hopefully have 0 issues between them.

Here is the full list of environment supported for this:

- [Basic Golang Application + Deployments](#basic-golang-application--deployments)
  - [Golang App on local workstation](#golang-app-on-local-workstation)
  - [Docker](#docker)
  - [Google Compute Engine](#google-compute-engine)
  - [Google Kubernetes Engine](#google-kubernetes-engine)


## Golang App on local workstation

Build the application and then run it

```bash
CGO_ENABLED=0 go build -o app main.go
```

Then, we can simply just run it after that

```bash
./app
```

## Docker

Build the container image (there are multiple versions of it) and then run it

```bash
docker build -t lol -f ./deployment/docker/alpine.Dockerfile .
docker run -p 8080:8080 lol
```

## Google Compute Engine

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

## Google Kubernetes Engine

First build the container with the tag that allows us to push it into container registry/artifact registry.

For container registry:

```bash
docker build -t gcr.io/<project id>/basic:v1 -f ./deployments/docker/slim.Dockerfile .
docker push gcr.io/<project id>/basic:v1
```

Then, we would need to run the following command to create a Auto-GKE cluster

```bash
gcloud container clusters create-auto "autopilot-cluster-1" --region "us-central1"
```

Next, we would run the following command to create the deployment service

```bash
kubectl create deployment basic-app --image=gcr.io/<project id>/basic:v1
```