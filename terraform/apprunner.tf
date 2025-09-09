# TODO: Remove
# resource "aws_apprunner_service" "bite_tracker" {
#   service_name = local.name
#
#   auto_scaling_configuration_arn = aws_apprunner_auto_scaling_configuration_version.bite_tracker.arn
#
#   instance_configuration {
#     cpu               = "0.25 vCPU"
#     memory            = 512
#     instance_role_arn = aws_iam_role.bite_tracker_instance_role.arn
#   }
#
#   source_configuration {
#     auto_deployments_enabled = false
#
#     image_repository {
#       image_configuration {
#         port = "8080"
#         runtime_environment_variables = {
#           BT_LISTEN_ADDR        = ":8080"
#           BT_TOKEN_AGE          = "24h"
#           BT_DB_NAME            = aws_db_instance.bite_tracker.db_name
#           BT_DB_HOST            = aws_db_instance.bite_tracker.address
#           BT_DB_PORT            = aws_db_instance.bite_tracker.port
#           BT_DB_SSL_MODE        = "require"
#           BT_DB_BOOTSTRAP_ROLES = tostring(var.bootstrap_db_roles)
#         }
#         runtime_environment_secrets = {
#           BT_DB_APP_USER_USERNAME = aws_secretsmanager_secret.db_app_user_username.arn
#           BT_DB_APP_USER_PASSWORD = aws_secretsmanager_secret.db_app_user_password.arn
#           # Only use admin on first run for bootstrapping roles
#           BT_DB_MIGRATE_USER_USERNAME = var.bootstrap_db_roles ? aws_secretsmanager_secret.db_admin_user_username.arn : aws_secretsmanager_secret.db_app_user_username.arn
#           BT_DB_MIGRATE_USER_PASSWORD = var.bootstrap_db_roles ? aws_secretsmanager_secret.db_admin_user_password.arn : aws_secretsmanager_secret.db_app_user_password.arn
#           BT_HMAC_TOKEN_SECRET        = aws_secretsmanager_secret.hmac_token_secret.arn
#
#         }
#       }
#       image_identifier      = "${aws_ecr_repository.bite_tracker.repository_url}:latest"
#       image_repository_type = "ECR"
#
#     }
#
#     authentication_configuration {
#       access_role_arn = aws_iam_role.bite_tracker_access_role.arn
#     }
#   }
#
#   network_configuration {
#     egress_configuration {
#       egress_type       = "VPC"
#       vpc_connector_arn = aws_apprunner_vpc_connector.bite_tracker.arn
#     }
#   }
#
#   tags = local.default_tags
# }
#
# resource "aws_apprunner_auto_scaling_configuration_version" "bite_tracker" {
#   auto_scaling_configuration_name = local.name
#
#   # Don't want autoscaling on demo
#   max_concurrency = 50
#   max_size        = 1
#   min_size        = 1
# }
#
# resource "aws_apprunner_vpc_connector" "bite_tracker" {
#   vpc_connector_name = local.name
#   subnets            = [aws_subnet.bite_tracker_app_a.id]
#   security_groups    = [aws_security_group.bite_tracker_app.id]
#   tags               = local.default_tags
# }
