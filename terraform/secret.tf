resource "aws_secretsmanager_secret" "db_app_user_connection_string" {
  name = "bt_db_app_user_connection_string"
  # Force delete since this is a demo app
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "db_migrate_user_connection_string" {
  name                    = "bt_db_migrate_user_connection_string"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "hmac_token_secret" {
  name                    = "bt_hmac_token_secret"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "db_app_user_connection_string" {
  secret_id     = aws_secretsmanager_secret.db_app_user_connection_string.id
  secret_string = "postgresql://bt_app_user:${var.db_app_user_password}@${aws_db_instance.bite_tracker.address}:${aws_db_instance.bite_tracker.port}/${aws_db_instance.bite_tracker.db_name}?sslmode=require"
}

resource "aws_secretsmanager_secret_version" "db_migrate_user_connection_string" {
  secret_id     = aws_secretsmanager_secret.db_migrate_user_connection_string.id
  secret_string = "postgresql://bt_app_user:${var.db_migrate_user_password}@${aws_db_instance.bite_tracker.address}:${aws_db_instance.bite_tracker.port}/${aws_db_instance.bite_tracker.db_name}?sslmode=require"
}

resource "aws_secretsmanager_secret_version" "hmac_token_secret" {
  secret_id     = aws_secretsmanager_secret.hmac_token_secret.id
  secret_string = var.hmac_token_secret
}
