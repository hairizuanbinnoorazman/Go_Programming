packer {
    required_plugins {
        googlecompute = {
            source  = "github.com/hashicorp/googlecompute"
            version = "~> 1"
        }
    }   
}

variable "gcp_project_id" {
    type      = string
    sensitive = true
}

variable "ssh_username" {
    type    = string
    default = "hairizuan"
}

variable "region" {
    type    = string
    default = "us-central1"
}

variable "zone" {
    type    = string
    default = "us-central1-a"
}

variable "gce_source_image" {
    type    = string
    default = "debian-11-bullseye-v20231010"
}

source "googlecompute" "basic-example" {
    project_id   = var.gcp_project_id
    source_image = var.gce_source_image
    ssh_username = var.ssh_username
    zone         = var.zone
}

build {
    sources = [
        "sources.googlecompute.basic-example"
    ]

    provisioner "file" {
        source      = "../bin/app"
        destination = "/home/hairizuan/app"
    }
    
    provisioner "shell" {
        inline = [
            "sudo mv /home/hairizuan/app /usr/local/bin/app"
        ]
    }
}