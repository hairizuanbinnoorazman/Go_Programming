module "localdocker" {
    source = "./modules/localdocker"

    image_name = "${var.image_name}"
    image_tag  = "${var.image_tag}"
}