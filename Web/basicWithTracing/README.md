# Basic

Contains basic example in order to deal with control a Golang application via Environment variables


```bash
kubectl -n kube-system create serviceaccount tiller
kubectl create clusterrolebinding tiller \                  
  --clusterrole cluster-admin \
  --serviceaccount=kube-system:tiller
helm init --service-account tiller
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm install jaegertracing/jaeger-operator

kubectl apply -f jaeger.yaml
kubectl apply -f deploy.yaml

kubectl port-forward service/simplest-query 8088:16686
```


