# The resource group is created outside Terraform together with the course budget (see README.md),
# so that destroying the cluster never removes the Terraform state stored in the same group.
data "azurerm_resource_group" "course" {
  name = var.resource_group_name
}

# ==============================================================================
# Container registry for the course images (Basic: cheapest SKU, no admin user)
# ==============================================================================
resource "azurerm_container_registry" "acr" {
  name                = var.acr_name
  resource_group_name = data.azurerm_resource_group.course.name
  location            = data.azurerm_resource_group.course.location
  sku                 = "Basic"
  admin_enabled       = false
  tags                = var.tags
}

# ==============================================================================
# AKS cluster
# ==============================================================================
resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.cluster_name
  resource_group_name = data.azurerm_resource_group.course.name
  location            = data.azurerm_resource_group.course.location
  dns_prefix          = var.cluster_name
  kubernetes_version  = var.kubernetes_version
  sku_tier            = "Free"
  node_resource_group = var.node_resource_group_name

  default_node_pool {
    name                        = "system"
    vm_size                     = var.vm_size
    node_count                  = var.node_count
    os_disk_size_gb             = 32
    temporary_name_for_rotation = "systemtmp"

    upgrade_settings {
      max_surge = "10%"
    }
  }

  # Node pools are managed here, not provisioned automatically
  node_provisioning_profile {
    mode = "Manual"
  }

  identity {
    type = "SystemAssigned"
  }

  # Sign in with Entra ID (kubelogin) and authorize with Azure RBAC; no static admin credentials
  local_account_disabled = true
  azure_active_directory_role_based_access_control {
    azure_rbac_enabled = true
    tenant_id          = data.azurerm_client_config.current.tenant_id
  }

  # Workload identity lets pods and CI pipelines use Entra ID without secrets (later exercises)
  oidc_issuer_enabled       = true
  workload_identity_enabled = true

  network_profile {
    network_plugin      = "azure"
    network_plugin_mode = "overlay"
    load_balancer_sku   = "standard"
  }

  tags = var.tags
}

# The nodes pull images from ACR with their kubelet identity, so no image pull secrets are needed
resource "azurerm_role_assignment" "aks_acr_pull" {
  principal_id                     = azurerm_kubernetes_cluster.aks.kubelet_identity[0].object_id
  role_definition_name             = "AcrPull"
  scope                            = azurerm_container_registry.acr.id
  skip_service_principal_aad_check = true
}

# The user running Terraform gets cluster admin through Azure RBAC
resource "azurerm_role_assignment" "aks_user_admin" {
  principal_id         = data.azurerm_client_config.current.object_id
  role_definition_name = "Azure Kubernetes Service RBAC Cluster Admin"
  scope                = azurerm_kubernetes_cluster.aks.id
}

# ...and can push images to ACR
resource "azurerm_role_assignment" "user_acr_push" {
  principal_id         = data.azurerm_client_config.current.object_id
  role_definition_name = "AcrPush"
  scope                = azurerm_container_registry.acr.id
}
