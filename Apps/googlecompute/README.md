# App with Google Compute

This application serves as a test application that would to call the google compute api to list compute engine instances in a project

Application calls the google compute engine and list all instances in a project every few seconds

# Testing it out

While on local computer, you can try to experiment to see if it is possible to just use the default application credentials provided by gcloud would be sufficient to run the program.

```bash
# Edit to provide Project ID first
go build -o lol .
go run ./lol
```

For docker builds

```bash
docker build -t gcr.io/XXX-PROJECT-ID-XXX/sample-google-compute:v1
docker push -t gcr.io/XXX-PROJECT-ID-XXX/sample-google-compute:v1
```

On GKE - utilize workload identity - can skip embedding the of secrets into the application

```bash
gcloud iam service-accounts create sample-google-compute

gcloud projects add-iam-policy-binding XXXXXX \
--member="serviceAccount:sample-google-compute@XXXXXX.iam.gserviceaccount.com" \
--role="roles/editor"

gcloud iam service-accounts add-iam-policy-binding \
sample-google-compute@XXXXX.iam.gserviceaccount.com \
--member="serviceAccount:XXXXXX.svc.id.goog[default/default]" \
--role="roles/iam.workloadIdentityUser"

kubectl annotate serviceaccount \
  --namespace default \
  default \
  iam.gke.io/gcp-service-account=sample-google-compute@XXXXXX.iam.gserviceaccount.com
```


gcloud auth application-default login --no-launch-browser