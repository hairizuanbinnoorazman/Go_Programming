# Autoscaling - Client Side

This is a sample application that is meant to test out horizontal scaling using custom metrics. This has been tested on kubernetes 1.15+. However, experience with that is not that great.

Further tests needs to be done - maybe kubernetes 1.18+ which also incluses autoscaling behaviours - which allow proper tight control over metric scaling.

Before proceeding - do adjust the kustomization.yml file. And push the required image to the docker registry

## Usage

Before deploying, do ensure that you build the required image and push to an image registry that is accessible by Kubernetes cluster

```bash
# To deploy the components
make deploy

# To drop the components
make drop
```

## Metrics to be used in grafana for comparison

Count of number of pods  
`sum(kube_pod_container_info{container="custom-hpa-client"})`

Number of items in queue
`testservice_generated_queue_item`