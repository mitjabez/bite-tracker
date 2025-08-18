resource "aws_secretsmanager_secret" "db_admin_user_username" {
  name                    = "${local.name}-db-admin-user-username"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "db_admin_user_password" {
  name                    = "${local.name}-db-admin-user-password"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "db_app_user_username" {
  name                    = "${local.name}-db-app-user-username"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "db_app_user_password" {
  name                    = "${local.name}-db-app-user-password"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "hmac_token_secret" {
  name                    = "${local.name}-hmac-token-secret"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "db_admin_user_username" {
  secret_id     = aws_secretsmanager_secret.db_admin_user_username.id
  secret_string = var.db_admin_user_username
}

resource "aws_secretsmanager_secret_version" "db_admin_user_password" {
  secret_id     = aws_secretsmanager_secret.db_admin_user_password.id
  secret_string = var.db_admin_user_password
}

resource "aws_secretsmanager_secret_version" "db_app_user_username" {
  secret_id     = aws_secretsmanager_secret.db_app_user_username.id
  secret_string = var.db_app_user_username
}

resource "aws_secretsmanager_secret_version" "db_app_user_password" {
  secret_id     = aws_secretsmanager_secret.db_app_user_password.id
  secret_string = var.db_app_user_password
}

resource "aws_secretsmanager_secret_version" "hmac_token_secret" {
  secret_id     = aws_secretsmanager_secret.hmac_token_secret.id
  secret_string = var.hmac_token_secret
}
