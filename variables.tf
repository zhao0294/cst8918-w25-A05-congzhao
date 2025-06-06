variable "labelPrefix" {
  description = "Prefix for all resource names"
  type        = string
}

variable "region" {
  default = "eastus"
}

variable "admin_username" {
  default = "azureadmin"
}

variable "ssh_public_key" {
  description = "Path to your local SSH public key"
  type        = string
}