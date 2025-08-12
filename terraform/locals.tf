locals {
  name   = "bite-tracker"
  region = "eu-central-1"
  zone_a = "${local.region}a"
  zone_b = "${local.region}b"
  default_tags = {
    App = local.name
    Env = "dev"
  }
}

