# Provider Configuration
# ---------------------
terraform {
  required_version = ">= 1.0"

  required_providers {
    tf = {
      source  = "patrikkj/tf"
      version = "0.0.1"
    }
  }
}
provider "tf" {}

# File Resource Examples
# --------------------

# Basic file management
resource "tf_local_file" "app_config" {
  path = "config.json"
  content = jsonencode({
    database_url = "postgresql://db.internal:5432/myapp"
    api_key      = var.api_key
    environment  = var.environment
  })
}

# File with sensitive content
resource "tf_local_file" "secure_file" {
  path    = "secure.txt"
  content = "Sensitive content"
}

# File in nested directory
resource "tf_local_file" "nested_file" {
  path    = "nested/dir/config.yml"
  content = <<-EOT
    environment: ${var.environment}
    debug: false
  EOT
}

# File Data Source Example
# ----------------------
data "tf_local_file" "existing_config" {
  path = "existing_config.yml"
}

# Command Execution Resource Example
# -------------------------------
resource "tf_local_exec" "service_deployment" {
  command = <<-EOT
    echo "Deploying service..."
    echo "Environment: ${var.environment}"
  EOT

  on_destroy      = "echo 'Cleaning up...'"
  fail_if_nonzero = true
}

# Command Execution Data Source Example
# ---------------------------------
data "tf_local_exec" "system_info" {
  command         = "uname -a"
  fail_if_nonzero = true
}

# Outputs
# -------
output "config_content" {
  value     = data.tf_local_file.existing_config.content
  sensitive = true
}

output "system_info" {
  value = data.tf_local_exec.system_info.output
}

# Variables
# --------
variable "api_key" {
  type      = string
  sensitive = true
}

variable "environment" {
  type    = string
  default = "production"
}
