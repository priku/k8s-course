# ==============================================================================
# Identity for the GitHub Actions pipeline of the project (exercise 3.6).
# The workflow signs in to Azure with GitHub's OIDC token: Entra ID trusts tokens from
# jobs that run in the GitHub environment `aks` of this repository, and exchanges them for
# a token of this managed identity. No passwords or keys are stored in GitHub.
# This is the Azure counterpart of the GCP service account and Workload Identity
# Federation in the course material.
# ==============================================================================
resource "azurerm_user_assigned_identity" "github_actions" {
  name                = "id-dwk-github-actions"
  resource_group_name = data.azurerm_resource_group.course.name
  location            = data.azurerm_resource_group.course.location
  tags                = var.tags
}

resource "azurerm_federated_identity_credential" "github_environment" {
  name                      = "github-${replace(var.github_repository, "/", "-")}-env-${var.github_environment}"
  user_assigned_identity_id = azurerm_user_assigned_identity.github_actions.id
  audience                  = ["api://AzureADTokenExchange"]
  issuer                    = "https://token.actions.githubusercontent.com"
  # Any branch can deploy, but only through a job that declares `environment: aks`
  subject = "repo:${var.github_repository}:environment:${var.github_environment}"
}

# Push images to the registry
resource "azurerm_role_assignment" "github_acr_push" {
  principal_id         = azurerm_user_assigned_identity.github_actions.principal_id
  principal_type       = "ServicePrincipal"
  role_definition_name = "AcrPush"
  scope                = azurerm_container_registry.acr.id
}

# Start the cluster if the nightly stop has stopped it
resource "azurerm_role_assignment" "github_aks_contributor" {
  principal_id         = azurerm_user_assigned_identity.github_actions.principal_id
  principal_type       = "ServicePrincipal"
  role_definition_name = "Azure Kubernetes Service Contributor Role"
  scope                = azurerm_kubernetes_cluster.aks.id
}

# Get a kubeconfig for the cluster
resource "azurerm_role_assignment" "github_aks_user" {
  principal_id         = azurerm_user_assigned_identity.github_actions.principal_id
  principal_type       = "ServicePrincipal"
  role_definition_name = "Azure Kubernetes Service Cluster User Role"
  scope                = azurerm_kubernetes_cluster.aks.id
}

# Deploy into the cluster, including creating namespaces for branch environments (3.7)
resource "azurerm_role_assignment" "github_aks_rbac_admin" {
  principal_id         = azurerm_user_assigned_identity.github_actions.principal_id
  principal_type       = "ServicePrincipal"
  role_definition_name = "Azure Kubernetes Service RBAC Cluster Admin"
  scope                = azurerm_kubernetes_cluster.aks.id
}
