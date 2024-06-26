module "localdocker" {
    source = "./modules/localdocker"

    image_name = "${var.image_name}"
    image_tag  = "${var.image_tag}"
}

module "gke" {
    source = "./modules/gke"

    gcp_project_id = "${var.gcp_project_id}"
}