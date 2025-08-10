resource "aws_apprunner_service" "bite_tracker" {
  service_name = local.name

  auto_scaling_configuration_arn = aws_apprunner_auto_scaling_configuration_version.bite_tracker.arn

  instance_configuration {
    cpu    = "0.25 vCPU"
    memory = 512
  }

  source_configuration {
    auto_deployments_enabled = false

    image_repository {
      image_configuration {
        port = "8080"
        runtime_environment_variables = {
          BT_LISTENADDR = ":8080"
          BT_TOKENAGE   = "24h"
        }
        runtime_environment_secrets = {
          BT_DBAPPURL        = aws_secretsmanager_secret.db_app_user_connection_string.arn
          BT_DBMIGRATEURL    = aws_secretsmanager_secret.db_migrate_user_connection_string.arn
          BT_HMACTOKENSECRET = aws_secretsmanager_secret.hmac_token_secret.arn
        }
      }
      image_identifier      = "${aws_ecr_repository.bite_tracker.repository_url}:latest"
      image_repository_type = "ECR"
    }

    authentication_configuration {
      access_role_arn = aws_iam_role.bite_tracker.arn
    }
  }

  network_configuration {
    egress_configuration {
      egress_type       = "VPC"
      vpc_connector_arn = aws_apprunner_vpc_connector.bite_tracker.arn
    }
  }

  tags = local.default_tags
}

resource "aws_apprunner_auto_scaling_configuration_version" "bite_tracker" {
  auto_scaling_configuration_name = local.name

  # Don't want autoscaling on demo
  max_concurrency = 50
  max_size        = 1
  min_size        = 1
}

resource "aws_apprunner_vpc_connector" "bite_tracker" {
  vpc_connector_name = local.name
  subnets            = [aws_subnet.bite_tracker_app_a.id]
  security_groups    = [aws_security_group.bite_tracker_app.id]
  tags               = local.default_tags
}
