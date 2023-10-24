variable "gcp_project_id" {
    description = "Define the GCP Project ID that we will interacting with"
    type        = string
    sensitive   = true
}

variable "gcp_region" {
    description = "Region to deploy the zone to"
    type        = string
    default     = "us-east1"
}

variable "gcp_zone" {
    description = "Zone to deploy the cluster to"
    type        = string
    default     = "us-east1-a"
}

variable "cluster_name" {
    description = "Name of cluster to deploy"
    type        = string
    default     = "cluster-1"
}

variable "image_name" {
    description = "Name of image to be deployed"
    type        = string
    default     = "basic-app"
}

variable "image_tag" {
    description = "Tag of image to be deployed"
    type        = string
    default     = "v1"
}