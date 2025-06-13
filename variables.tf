variable "labelPrefix" {
  description = "Prefix for all resource names"
  type        = string
}

variable "region" {
  type    = string
  default = "eastus"
}

variable "admin_username" {
  type    = string
  default = "azureadmin"
}

variable "ssh_public_key" {
  description = "Path to your local SSH public key"
  type        = string
}