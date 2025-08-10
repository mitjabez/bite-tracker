variable "ssh_public_key_path" {
  type        = string
  description = "Path to the SSH public key which provides access to the jumpbox"
}

variable "db_admin_username" {
  type        = string
  default     = "postgres"
  description = "Username of the Postgres DB Admin"
}

variable "db_admin_password" {
  type        = string
  sensitive   = true
  description = "Password of the admin Postgres user"
}

variable "db_app_user_password" {
  type        = string
  sensitive   = true
  description = "Password of user to execute Postgres queries from the app"
}

variable "db_migrate_user_password" {
  type        = string
  sensitive   = true
  description = "Password of user to execute Postgres migrations"
}

variable "hmac_token_secret" {
  type        = string
  sensitive   = true
  description = "Secret for signit JWT tokens"
}
