# Basic Web Application server

This folder is to showcase web applications to print secrets from mounted GCP Secrets

# Curling

```
curl localhost:8080/secretenv?env=TEST
curl localhost:8080/secretfile
```

# Local workstation

```
go build -o lol .
SECRET_PATH=./README.md lol
```

# Docker environment

```
docker build -t lol .
docker run -e TEST=somesecret -p 8080:8080 lol
```