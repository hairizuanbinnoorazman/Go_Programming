# Leader Election for golang apps in Kubernetes

Leader election for golang apps in Kubernetes by using the `leaderelection` portion of the `client-go` Kubernetes client library.

## Quickstart

```bash
# To build the docker image and push it to a container registry
make build

# To deploy it
make deploy
```