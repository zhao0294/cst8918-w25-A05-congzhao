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

# -------------------
# Network Security Group (NSG)
# -------------------
resource "azurerm_network_security_group" "this" {
  name                = "${var.labelPrefix}-nsg"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name

  security_rule {
    name                       = "Allow-SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "Allow-HTTP"
    priority                   = 1002
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "80"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

# -------------------
# Network Interface (NIC)
# -------------------
resource "azurerm_network_interface" "this" {
  name                = "${var.labelPrefix}-nic"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.this.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.this.id
  }

  network_security_group_id = azurerm_network_security_group.this.id
}

# -------------------
# cloud-init script for Apache installation
# -------------------
data "template_file" "init" {
  template = file("${path.module}/init.sh")
}

data "cloudinit_config" "config" {
  gzip          = false
  base64_encode = false

  part {
    filename     = "init.sh"
    content_type = "text/x-shellscript"
    content      = data.template_file.init.rendered
  }
}

# -------------------
# Linux Virtual Machine
# -------------------
resource "azurerm_linux_virtual_machine" "this" {
  name                  = "${var.labelPrefix}-vm"
  resource_group_name   = azurerm_resource_group.this.name
  location              = azurerm_resource_group.this.location
  size                  = "Standard_B1s"
  admin_username        = var.admin_username
  network_interface_ids = [azurerm_network_interface.this.id]

  disable_password_authentication = true

  admin_ssh_key {
    username   = var.admin_username
    public_key = file(var.ssh_public_key)
  }

  os_disk {
    name                 = "${var.labelPrefix}-osdisk"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "20_04-lts"
    version   = "latest"
  }

  custom_data = data.cloudinit_config.config.rendered
}