# App to demonstrate confusing helm hook config

This is a basic app that attempts to demonstrate the effect of suddenly introducing helm annotations to a resource without considering side effects.

This is run with Helm v3

Run the following commands:

First comment out helm annotations on configmap file and then deploy the helm chart

```bash
helm upgrade --install yahoo ./basic-app
```

Edit the `values.yaml` and add more values to pod annotations. This would force redeploy of pods. Rerun step to install helm chart

```bash
helm upgrade --install yahoo ./basic-app
```

You should facing the issue of where pods are stuck in ContainerCreating phase
