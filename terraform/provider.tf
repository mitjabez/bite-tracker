provider "aws" {
  region = var.region
  default_tags {
    tags = {
      Env       = var.env
      Terraform = "true"
    }
  }
}

