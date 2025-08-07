locals {
  name   = "bite-tracker"
  zone_a = "eu-central-1a"
  zone_b = "eu-central-1b"
  default_tags = {
    App = local.name
    Env = "dev"
  }
}


