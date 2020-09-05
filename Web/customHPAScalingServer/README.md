# Autoscaling - Server Side

This is a sample application that is meant to test out horizontal scaling using custom metrics. This has been tested on kubernetes 1.15+. However, experience with that is not that great.

Further tests needs to be done - maybe kubernetes 1.18+ which also incluses autoscaling behaviours - which allow proper tight control over metric scaling.

Before proceeding - do adjust the kustomization.yml file. And push the required image to the docker registry

## Usage

```bash
# To deploy the components
make deploy

# To drop the components
make drop
```