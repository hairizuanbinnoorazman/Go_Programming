# Basic Golang Application + Deployments

Contains basic example in order to deal with control a Golang application via Environment variables. The following repository serves to provide the basic application to deploy the same app into various different environments and hopefully have 0 issues between them.

Here is the full list of environment supported for this:

- [Basic Golang Application + Deployments](#basic-golang-application--deployments)
  - [Deploy via Bash/Makefiles](#deploy-via-bashmakefiles)
    - [Golang App on local workstation](#golang-app-on-local-workstation)
    - [Docker](#docker)
    - [Google Compute Engine](#google-compute-engine)
    - [Google Kubernetes Engine](#google-kubernetes-engine)
    - [Google Cloud Functions](#google-cloud-functions)
    - [Google Cloud Run](#google-cloud-run)
    - [Google App Engine](#google-app-engine)
    - [Amazon EC2 Instance](#amazon-ec2-instance)
    - [Amazon Elastic Container Service (ECS)](#amazon-elastic-container-service-ecs)
  - [Deploy via Terraform](#deploy-via-terraform)
    - [Docker](#docker-1)
    - [Google Compute Engine](#google-compute-engine-1)
    - [Google Kubernetes Engine](#google-kubernetes-engine-1)


## Deploy via Bash/Makefiles

### Golang App on local workstation

Build the application and then run it

```bash
CGO_ENABLED=0 go build -o app main.go
```

Then, we can simply just run it after that

```bash
./app
```

### Docker

Build the container image (there are multiple versions of it) and then run it

```bash
docker build -t lol -f ./deployment/docker/alpine.Dockerfile .
docker run -p 8080:8080 lol
```

### Google Compute Engine

First build the linux binary for the application

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .
```

Then, we would need to create the Google Compute Engine instance

```bash
gcloud compute instances create instance-1 --project=<project id> --zone=us-central1-a --machine-type=e2-medium
```

Then, scp the binary into the server

```bash
scp app hairizuan@<ip hostname>:/usr/local/bin/app
ssh hairizuan@<ip hostname>
chmod +x /usr/local/bin/app
app
```

### Google Kubernetes Engine

First build the container with the tag that allows us to push it into container registry/artifact registry.

For container registry:

```bash
docker build -t gcr.io/<project id>/basic:v1 -f ./deployments/docker/slim.Dockerfile .
docker push gcr.io/<project id>/basic:v1
```

Then, we would need to run the following command to create a Auto-GKE cluster

```bash
gcloud container clusters create-auto "autopilot-cluster-1" --region "us-central1"

# Get access to GKE cluster
gcloud container clusters get-credentials autopilot-cluster-1 --zone us-central1-a
```

Next, we would run the following command to create the deployment service

```bash
kubectl create deployment basic-app --image=gcr.io/<project id>/basic:v1
```

### Google Cloud Functions

TODO

### Google Cloud Run

First build the container with the tag that allows us to push it into container registry/artifact registry.

For container registry:

```bash
docker build -t gcr.io/<project id>/basic-app:v1 -f ./deployments/docker/slim.Dockerfile .
docker push gcr.io/<project id>/basic-app:v1
```

Run gcloud command to create Cloud Run service

```bash
gcloud run deploy basic-app --image=gcr.io/<project id>/basic-app:v1 --concurrency=10 --max-instances=1 --platform=managed --allow-unauthenticated --ingress=all --cpu=1 --memory=500Mi --region=us-east1
```

### Google App Engine

TODO

### Amazon EC2 Instance

We would first need to setup AWS CLI tool. We can do it in multiple ways
- Set up a "super" user that has permissions to do everything. We can then embed the `access_keys` and `secret_keys`
- Set a "empty" user that has `access_keys` and `secret_keys` but at the same time, we would create roles and allow our empty user to "assume" our role. This would allow us to drop access if ever there is a need for it.

If we are to utilize the console to create the instance - to manualy ssh in, we need to utilize the `ec2-user` user in the instance. Depends on what is being OS being used: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/managing-users.html

- Need to define default VPC and Security Group if it doesn't exist

```bash
# Get instance IDs
aws ec2 describe-instances | yq '.Reservations.[].Instances.[].InstanceId'  -
aws ec2 run-instances --image-id=ami-0dbc3d7bc646e8516 --count=1 --instance-type=t2.micro --key-name="Hairizuan Laptop - Key 2"
aws ec2 stop-instances --instance-ids <values>
aws ec2 terminate-instances --instance-ids <values>
```

### Amazon Elastic Container Service (ECS)

We would first need to push the container into Amazon ECR (Elastic Container Registry). To push image into the registry, we would first build the image, then we would need to set up credentials

```bash
# For convenience purposes - to get account id
aws sts get-caller-identity

# Need to be called once in awhile since auth token got expired
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <account-id>.dkr.ecr.<region>.amazonaws.com
```

ECR is quite different from Google's Container Registry - for each type of "application", we would need to create a new "repository" which we can then push the images in. E.g. in the case of the application here - we would first need to create a repository on ECR with the name `basic` - afterwhich, we can then push: `<aws specific registry url>/basic:v1`. We can't do this: `<aws specific registry url>/basic/basic:v1` - this would fail (docker command for pushing it would keep retrying till it eventually would fail)


```bash
aws ecr create-repository --repository-name=basic --region=<region>
```

We would then to build the docker image and then push it in.

```bash
docker build -t <account-id>.dkr.ecr.<region>.amazonaws.com/basic:v1 -f ./deployment/docker/slim.Dockerfile .
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/basic:v1
```

Next steps would be the following:
- Creating task definitions.
- We're going for a "simple" deployment - so we'll opt for a ECS Service -> essentially, this wraps tasks (although some pages say that ECS Service -> ECS Taskset -> ECS Tasks). However, we don't need to interact with ECS Taskset - we simply rely on service to manage number of tasks to be run.
- Create service call is "not idempotent" - need to run update-service instead of create service when trying to update it.
- Don't forget to add container port - we would need this since the basic-golang app is only exposed on port 8080

```bash
aws ecs register-task-definition --family basic --requires-compatibilities FARGATE --container-definitions '[{"name":"web","image":"<account-id>.dkr.ecr.<region>.amazonaws.com/basic:v1","essential":true,"portMappings":[{"containerPort":8080}]}]' --memory 2048 --cpu 512 --network-mode awsvpc --execution-role-arn arn:aws:iam::<account-id>:role/ecsTaskExecutionRole

aws ecs list-task-definitions

aws ecs create-service --cluster default --service-name basic-test --task-definition arn:aws:ecs:<region-id>:<account-id>:task-definition/basic:1 --desired-count 2 --launch-type FARGATE --network-configuration 'awsvpcConfiguration={subnets=[<subnet-id>],securityGroups=[<security-group-id>],assignPublicIp=ENABLED}'

aws ecs update-service --service basic-test --task-definition arn:aws:ecs:<region-id>:<account-id>:task-definition/basic:2 
```

For cleaning it up

```bash
aws ecs update-service --service basic-test --task-definition arn:aws:ecs:<region-id>:<account-id>:task-definition/basic:2 --desired-count=0

# Apparently can't delete if it scaled more than 0
aws ecs delete-service --cluster=default --service=basic-test

# Check if tasks still running
aws ecs list-services
aws ecs list-tasks

aws ecs deregister-task-definition --task-definition basic:1
aws ecs deregister-task-definition --task-definition basic:2
aws ecs delete-task-definitions --task-definitions='["basic:1","basic:2"]'

# Apparently, there is delete protection on repositories (that's somewhat logical)
aws ecr delete-repository --repository-name basic --region <region> --force
```

## Deploy via Terraform


Some drawbacks (comparing it to other tools like Ansible).  
- Unable to simply pass variables from tfvars down to modules easily. Design conflict with initial aim of terraform. Compare the experience of this compared to Ansible  
  https://github.com/hashicorp/terraform/issues/32508

Useful references/notes
- Name of module defined in root `main.tf` is used in the `-target` flag. If not matched, it will not trigger
- If `deletion_proection` is true. First, make the change, create plan and apply it. After which, attempt to destroy after. (Alternatively, can change tfstate file - but that's not recommended)
- terraform registry: https://registry.terraform.io/
- https://spacelift.io/blog/how-to-use-terraform-variables
- https://robertdebock.nl/learn-terraform/

### Docker

NOTE: Seems hard to also include build step in terraform (need to replicate the dockerfile into terraform itself - it seems to have issues with relative paths? Might require a bit of debugging)

From terraform folder, run the following command

```bash
terraform plan -target=module.localdocker -var-file=local.tfvars -out localdocker.plan
terraform apply localdocker.plan
```

To destroy it

```bash
terraform plan -target=module.localdocker -var-file=local.tfvars -destroy -out localdocker.plan
terraform apply localdocker.plan
```

### Google Compute Engine

Require step to build Google Cloud Image (which we would then use to deploy the VM). This is done by relying on packer. Packer provides the tooling that we would need to first create the image that we wish to rely on and then - we can create instances out of it.

```bash
# Just to test that packer script works
packer validate -var 'gcp_project_id=aaa'  gce.pkr.hcl

# Building it
packer build -var 'gcp_project_id=XXX' gce.pkr.hcl
```

### Google Kubernetes Engine

It is assumed that the images have already been pushed into a container registry - all we have to do is to deploy it into a cluster.

```bash
TF_VAR_gcp_project_id=<project id> terraform plan -target=module.gke -var-file=local.tfvars -out gke.plan
terraform apply gke.plan
```

To destroy it 

```bash
TF_VAR_gcp_project_id=<project id> terraform plan -target=module.gke -var-file=local.tfvars -destroy -out gke.plan
terraform apply gke.plan
```

Convienience make commands available to create the plans via:

```bash
# Create plan for creating the infra
make tf_create_gke

# Create plan for destroying the infra
make tf-destroy_gke
```