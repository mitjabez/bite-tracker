resource "aws_secretsmanager_secret" "db_admin_user_username" {
  name                    = "bt_db_admin_user_username"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "db_admin_user_password" {
  name                    = "bt_db_admin_user_password"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "db_app_user_username" {
  name                    = "bt_db_app_user_username"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "db_app_user_password" {
  name                    = "bt_db_app_user_password"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret" "hmac_token_secret" {
  name                    = "bt_hmac_token_secret"
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
