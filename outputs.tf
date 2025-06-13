output "resource_group_name" {
  value = azurerm_resource_group.this.name
}

output "public_ip" {
  value = azurerm_public_ip.this.ip_address
}