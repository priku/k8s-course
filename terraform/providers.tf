terraform {
  required_version = ">= 1.9, < 2.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.46.0, < 5.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"
    }
  }

  # Remote backend for state storage (CI/CD)
  backend "azurerm" {
    resource_group_name  = "dwk-tfstate-rg"
    storage_account_name = "dwktfstatep3dqbo"
    container_name       = "tfstate"
    key                  = "dwk-aks.tfstate"
  }
}

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
  subscription_id = "3921f8c4-f431-4e3f-b001-fa10bf905e12"
}

# Get current Azure client configuration
data "azurerm_client_config" "current" {}
