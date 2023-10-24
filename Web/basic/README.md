# Basic Golang Application + Deployments

Contains basic example in order to deal with control a Golang application via Environment variables. The following repository serves to provide the basic application to deploy the same app into various different environments and hopefully have 0 issues between them.

Here is the full list of environment supported for this:

- [Basic Golang Application + Deployments](#basic-golang-application--deployments)
  - [Golang App on local workstation](#golang-app-on-local-workstation)
  - [Docker](#docker)
  - [Google Compute Engine](#google-compute-engine)


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

Then, scp the binary into the server

```bash
scp app hairizuan@<ip hostname>:/usr/local/bin/app
ssh hairizuan@<ip hostname>
chmod +x /usr/local/bin/app
app
```