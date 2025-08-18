variable "env" {
  type        = string
  description = "Deployment environment (dev, prod)"
}

variable "region" {
  type        = string
  description = "Region of the infrastructure"
}

variable "db_admin_user_username" {
  type        = string
  description = "Username of the Postgres DB Admin"
}

variable "db_admin_user_password" {
  type        = string
  sensitive   = true
  description = "Password of the admin Postgres user"
}

variable "db_app_user_username" {
  type        = string
  sensitive   = true
  description = "User of user to execute Postgres queries from the app"
}

variable "db_app_user_password" {
  type        = string
  sensitive   = true
  description = "Password of user to execute Postgres queries from the app"
}

variable "bootstrap_db_roles" {
  type        = bool
  description = "Whether to bootstrap DB roles at startup"
}

variable "hmac_token_secret" {
  type        = string
  sensitive   = true
  description = "Secret for signit JWT tokens"
}
