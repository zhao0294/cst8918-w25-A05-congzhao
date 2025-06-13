# CST8918 Lab5 & Lab6 Key Operations and Common Issues

## Table of Contents
- [Lab5 Overview](#lab5-overview)
- [Lab5 Key Steps](#lab5-key-steps)
- [Lab5 Common Issues & Solutions](#lab5-common-issues--solutions)
- [Lab6 Overview](#lab6-overview)
- [Lab6 Key Steps](#lab6-key-steps)
- [Lab6 Common Issues & Solutions](#lab6-common-issues--solutions)
- [Environment Setup](#environment-setup)
- [References](#references)

---

## Lab5 Overview
- Deploy an Azure Linux VM configured with cloud-init
- Validate basic Terraform configuration and resource provisioning
- Use Terratest for automated infrastructure testing

## Lab5 Key Steps
1. **Terraform Configuration**
   - Prepare `main.tf`, `variables.tf`, and cloud-init scripts
2. **Terraform Execution**
   - Run `terraform init`
   - Run `terraform apply -auto-approve`
3. **Terratest Implementation**
   - Use `terraform.Options` to specify directory
   - Test VM creation, state, and network connectivity
4. **SSH Verification**
   - Use Terratest SSH module to verify Apache service is running

## Lab5 Common Issues & Solutions
- **Place go.mod inside the test directory, not the root, to avoid module conflicts**
- **Ensure compatible versions of Azure SDK and Terratest to avoid errors like `undefined` or `missing go.sum`**
- **Terraform output variables (`public_ip`, `resource_group_name`) must be correct for tests to succeed**
- **cloud-init syntax errors may cause VM startup failure or missing services**

---

## Lab6 Overview
- Use Terratest for advanced integration testing
- Automate verification of Terraform-created resources (NICs, VM image versions)
- Implement helper functions to wrap Azure SDK calls

## Lab6 Key Steps
1. **Separate Azure SDK calls into a dedicated `azure` package**
2. **Implement helper functions:**
   - List VMs in a resource group
   - Get NICs attached to a VM
   - Retrieve VM image details
3. **Write test functions to:**
   - Validate resources and service states in one place
   - Print unified test summary after execution
4. **Manage dependencies and version compatibility:**
   - Use independent `go.mod` for the azure package
   - Use `replace` directive in test module `go.mod` to link local azure package

## Lab6 Common Issues & Solutions
- **Azure SDK version conflicts, especially `armcompute` and `azidentity`; pin versions and upgrade to stable releases**
- **Compilation errors like `undefined: azure.ListVirtualMachinesForResourceGroup` due to wrong import paths**
- **Ensure `replace` directive correctly points to local azure package path**
- **Set `ARM_SUBSCRIPTION_ID` environment variable before running tests**

---

## Environment Setup
- Go version: 1.24.x
- Terraform version: >= 1.4.x
- Azure SDK for Go:
  - azidentity v1.10.1
  - armcompute v1.0.0
- Terratest v0.49.0
- SSH Key: stored at `$HOME/.ssh/id_rsa` for SSH connectivity tests

---

## References
- [Terraform Azure Provider Documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs)
- [Azure SDK for Go GitHub](https://github.com/Azure/azure-sdk-for-go)
- [Terratest Official Documentation](https://terratest.gruntwork.io/docs/)
- [Go Modules Official Blog](https://go.dev/blog/using-go-modules)

---

*Author: Cong Zhao  
Date: 2025-06-13*
