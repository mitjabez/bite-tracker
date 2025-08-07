resource "aws_apprunner_service" "bite_tracker" {
  service_name = "bite-tracker"

  auto_scaling_configuration_arn = aws_apprunner_auto_scaling_configuration_version.bite_tracker.arn

  instance_configuration {
    cpu    = 0.25
    memory = 512
  }

  source_configuration {
    auto_deployments_enabled = false

    image_repository {
      image_configuration {
        port = "8080"
      }
      image_identifier      = "${aws_ecr_repository.bite_tracker.repository_url}:latest"
      image_repository_type = "ECR"
    }

    authentication_configuration {
      access_role_arn = aws_iam_role.bite_tracker.arn
    }
  }

  tags = local.default_tags
}

resource "aws_apprunner_auto_scaling_configuration_version" "bite_tracker" {
  auto_scaling_configuration_name = "bite-tracker"

  # Don't want autoscaling on demo
  max_concurrency = 50
  max_size        = 1
  min_size        = 1
}
