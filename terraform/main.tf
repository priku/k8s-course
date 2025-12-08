# Resource Group
resource "azurerm_resource_group" "aks" {
  name     = var.resource_group_name
  location = var.location
  tags     = var.tags
}

# Random suffix for ACR (must be globally unique)
resource "random_string" "acr_suffix" {
  length  = 8
  special = false
  upper   = false
  numeric = true
}

# ==============================================================================
# Azure Container Registry - Using Azure Verified Module (AVM)
# https://github.com/Azure/terraform-azurerm-avm-res-containerregistry-registry
# ==============================================================================
module "acr" {
  source  = "Azure/avm-res-containerregistry-registry/azurerm"
  version = "0.4.0"

  name                = "dwkacr${random_string.acr_suffix.result}"
  resource_group_name = azurerm_resource_group.aks.name
  location            = azurerm_resource_group.aks.location

  # Basic SKU is sufficient for learning/demo, with zone redundancy disabled
  sku                      = "Basic"
  zone_redundancy_enabled  = false
  retention_policy_in_days = null  # Retention policy only available with Premium SKU

  # Enable telemetry for AVM
  enable_telemetry = true

  tags = var.tags
}

# ==============================================================================
# Azure Kubernetes Service - Using Azure Verified Module (AVM)
# https://github.com/Azure/terraform-azurerm-avm-res-containerservice-managedcluster
# ==============================================================================
module "aks" {
  source  = "Azure/avm-res-containerservice-managedcluster/azurerm"
  version = "0.3.3"

  name                = var.cluster_name
  resource_group_name = azurerm_resource_group.aks.name
  location            = azurerm_resource_group.aks.location

  # Default node pool configuration
  default_node_pool = {
    name       = "default"
    vm_size    = var.vm_size
    node_count = var.node_count

    upgrade_settings = {
      max_surge = "10%"
    }
  }

  # DNS prefix for the cluster
  dns_prefix = var.cluster_name

  # Use system-assigned managed identity
  managed_identities = {
    system_assigned = true
  }

  # Enable Azure AD RBAC
  azure_active_directory_role_based_access_control = {
    azure_rbac_enabled = true
    tenant_id          = data.azurerm_client_config.current.tenant_id
  }

  # Network configuration
  network_profile = {
    network_plugin = "azure"
  }

  # Enable telemetry for AVM
  enable_telemetry = true

  tags = var.tags
}

# ==============================================================================
# Role Assignment: Allow AKS to pull images from ACR
# ==============================================================================
resource "azurerm_role_assignment" "aks_acr" {
  principal_id                     = module.aks.kubelet_identity_id
  role_definition_name             = "AcrPull"
  scope                            = module.acr.resource_id
  skip_service_principal_aad_check = true
}
