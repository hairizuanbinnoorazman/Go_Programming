terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.1"
    }
  }
}

resource "docker_image" "basic_app" {
  name         = "${var.image_name}:${var.image_tag}"
  keep_locally = false
}

resource "docker_container" "basic_app" {
  image = docker_image.basic_app.image_id
  name  = "basic-app"

  ports {
    internal = 8080
    external = 8080
  }
}