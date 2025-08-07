resource "aws_ecr_repository" "bite_tracker" {
  name                 = "bite-tracker"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = local.default_tags
}
