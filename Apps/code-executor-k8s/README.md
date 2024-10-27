# Code Executor on Kubernetes

This is a sample project on how we acan build a code executor tool and have it run on Kubernetes

## Features

- Accept code on a webserver
- Create Configmap
- Create Job that references the configmap
- Once job completes - reports back to the webserver

## Issues

Cleanup pods

```bash
kubectl delete pods $(kubectl get pods | grep test-test | awk '{print $1}')
```