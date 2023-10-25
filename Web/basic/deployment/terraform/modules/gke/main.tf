terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.3.0"
    }

    kubernetes = {
        source = "hashicorp/kubernetes"
        version = "~> 2.23.0"
    }
  }
}

provider "google" {
    project = "${var.gcp_project_id}"
    region  = "${var.gcp_region}"
    zone    = "${var.gcp_zone}"
}

resource "google_container_cluster" "test_cluster" {
    name             = "${var.cluster_name}"
    enable_autopilot = true
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_gke_with_terraform
data "google_client_config" "provider" {}

data "google_container_cluster" "test_cluster" {
    name     = "${var.cluster_name}"
    depends_on = [ 
        google_container_cluster.test_cluster
    ]

}

provider "kubernetes" {
    host  = "https://${data.google_container_cluster.test_cluster.endpoint}"
    token = data.google_client_config.provider.access_token
    cluster_ca_certificate = base64decode(
        data.google_container_cluster.test_cluster.master_auth[0].cluster_ca_certificate,
    )
}

resource "kubernetes_deployment_v1" "basic_app_deployment" {
    metadata {
        name = "basic-app" 
        labels = {
            app = "basic-app"
        }
    }
    spec {
        replicas = 2
        selector {
            match_labels = {
                app = "basic-app"
            }   
        }
        template {
            metadata {
                labels = {
                    app = "basic-app"
                }
            }
            spec {
                container {
                    image = "gcr.io/${var.gcp_project_id}/${var.image_name}:${var.image_tag}"
                    name  = "basic-app"
                    resources {
                        limits = {
                            cpu    = "1"
                            memory = "512Mi"
                        }
                        requests = {
                            cpu    = "0.5"
                            memory = "256Mi"
                        }
                    }
                }
            }
        }
    }
}