terraform {
  required_version = ">= 1.17, < 2.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.11"
    }
  }

  # State in Azure Storage in the course resource group, read with Entra ID (no storage keys).
  # The storage account is created once with the commands in README.md.
  backend "azurerm" {
    resource_group_name  = "rg-dwk-aks-sc"
    storage_account_name = "stdwkakstfstatesc"
    container_name       = "tfstate"
    key                  = "dwk-aks.tfstate"
    use_azuread_auth     = true
  }
}

# The subscription comes from ARM_SUBSCRIPTION_ID (see README.md), so its ID is not in the repository
provider "azurerm" {
  features {}

  # azurerm 5 registers no resource providers by default; these are already registered in the subscription
  resource_provider_registrations = "none"
}

data "azurerm_client_config" "current" {}
