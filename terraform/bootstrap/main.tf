# Bootstrap Terraform State Storage
# Run this ONCE manually to create the storage account for remote state
# After this, the main terraform can use remote backend

terraform {
  required_version = ">= 1.9"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id = var.subscription_id
}

variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
}

variable "location" {
  description = "Azure region for state storage"
  type        = string
  default     = "swedencentral"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dwk"
}

# Resource group for Terraform state
resource "azurerm_resource_group" "tfstate" {
  name     = "${var.environment}-tfstate-rg"
  location = var.location

  tags = {
    purpose     = "terraform-state"
    managed-by  = "terraform-bootstrap"
    environment = var.environment
  }
}

# Storage account for Terraform state
resource "azurerm_storage_account" "tfstate" {
  name                     = "${var.environment}tfstate${random_string.suffix.result}"
  resource_group_name      = azurerm_resource_group.tfstate.name
  location                 = azurerm_resource_group.tfstate.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = "TLS1_2"

  blob_properties {
    versioning_enabled = true
  }

  tags = {
    purpose     = "terraform-state"
    managed-by  = "terraform-bootstrap"
    environment = var.environment
  }
}

# Random suffix for globally unique storage account name
resource "random_string" "suffix" {
  length  = 6
  special = false
  upper   = false
}

# Container for Terraform state files
resource "azurerm_storage_container" "tfstate" {
  name                  = "tfstate"
  storage_account_id    = azurerm_storage_account.tfstate.id
  container_access_type = "private"
}

# Outputs needed for backend configuration
output "resource_group_name" {
  value       = azurerm_resource_group.tfstate.name
  description = "Resource group for Terraform state"
}

output "storage_account_name" {
  value       = azurerm_storage_account.tfstate.name
  description = "Storage account for Terraform state"
}

output "container_name" {
  value       = azurerm_storage_container.tfstate.name
  description = "Container for Terraform state"
}

output "backend_config" {
  value       = <<-EOT
    # Add this to terraform/providers.tf after 'required_providers' block:
    backend "azurerm" {
      resource_group_name  = "${azurerm_resource_group.tfstate.name}"
      storage_account_name = "${azurerm_storage_account.tfstate.name}"
      container_name       = "${azurerm_storage_container.tfstate.name}"
      key                  = "dwk-aks.tfstate"
    }
  EOT
  description = "Backend configuration to add to main Terraform"
}
