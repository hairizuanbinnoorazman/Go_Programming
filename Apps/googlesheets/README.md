# Sample google spreadsheet application

# Deployment

Utilize workload identity to deploy the application without required creds on GKE clusters

```bash
docker build -t gcr.io/XXXXXX/sample-sheets:v2 .
docker push gcr.io/XXXXXX/sample-sheets:v2

gcloud iam service-accounts create sample-sheets

gcloud projects add-iam-policy-binding XXXXXX \
--member="serviceAccount:sample-sheets@XXXXXX.iam.gserviceaccount.com" \
--role="roles/editor"

gcloud iam service-accounts add-iam-policy-binding \
sample-sheets@XXXXXX.iam.gserviceaccount.com \
--member="serviceAccount:XXXXXX.svc.id.goog[default/default]" \
--role="roles/iam.workloadIdentityUser"

gcloud iam service-accounts add-iam-policy-binding \
sample-sheets@XXXXXX.iam.gserviceaccount.com \
--member="serviceAccount:XXXXXX.svc.id.goog[default/default]" \
--role="roles/iam.serviceAccountTokenCreator"

kubectl annotate serviceaccount \
  --namespace default \
   default \
   iam.gke.io/gcp-service-account=sample-sheets@XXXXXX.iam.gserviceaccount.com

```
