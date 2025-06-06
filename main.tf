terraform {
  required_version = ">= 1.1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
    cloudinit = {
      source  = "hashicorp/cloudinit"
      version = "2.3.3"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "cloudinit" {}

# -------------------
# Resource Group
# -------------------
resource "azurerm_resource_group" "this" {
  name     = "${var.labelPrefix}-A05-RG"
  location = var.region
}

# -------------------
# Public IP
# -------------------
resource "azurerm_public_ip" "this" {
  name                = "${var.labelPrefix}-pip"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  allocation_method   = "Dynamic"
  sku                 = "Basic"
}

# -------------------
# Virtual Network
# -------------------
resource "azurerm_virtual_network" "this" {
  name                = "${var.labelPrefix}-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
}

# -------------------
# Subnet
# -------------------
resource "azurerm_subnet" "this" {
  name                 = "${var.labelPrefix}-subnet"
  resource_group_name  = azurerm_resource_group.this.name
  virtual_network_name = azurerm_virtual_network.this.name
  address_prefixes     = ["10.0.1.0/24"]
}