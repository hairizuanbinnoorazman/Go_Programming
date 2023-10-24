variable "image_name" {
    description = "Name of docker image to run"
    type        = string
    default     = "basic-app"
}

variable image_tag {
    description = "Tag of docker image to run"
    type        = string
    default     = "v1"
}