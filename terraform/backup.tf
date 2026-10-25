# ==============================================================================
# Database backups (exercise 3.10): a CronJob in the namespace `project` dumps the todo
# database every day and uploads it here. The material uses Google Cloud Storage; this is
# Azure Blob Storage.
#
# The CronJob authenticates with workload identity, the approach the material recommends:
# its Kubernetes ServiceAccount `project/db-backup` is trusted by a federated credential of the
# managed identity below, so no keys or connection strings exist anywhere. The account
# does not even accept shared keys.
# ==============================================================================
resource "azurerm_storage_account" "backups" {
  name                            = var.backup_storage_account_name
  resource_group_name             = data.azurerm_resource_group.course.name
  location                        = data.azurerm_resource_group.course.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  access_tier                     = "Cool" # written daily, read rarely
  min_tls_version                 = "TLS1_2"
  shared_access_key_enabled       = false
  allow_nested_items_to_be_public = false
  tags                            = var.tags
}

resource "azurerm_storage_container" "todo_backups" {
  name                  = "todo-backups"
  storage_account_id    = azurerm_storage_account.backups.id
  container_access_type = "private"
}

# Keep the backups for 30 days
resource "azurerm_storage_management_policy" "backups" {
  storage_account_id = azurerm_storage_account.backups.id

  rule {
    name    = "delete-old-backups"
    enabled = true
    filters {
      blob_types   = ["blockBlob"]
      prefix_match = ["${azurerm_storage_container.todo_backups.name}/"]
    }
    actions {
      base_blob {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
}

resource "azurerm_user_assigned_identity" "db_backup" {
  name                = "id-dwk-db-backup"
  resource_group_name = data.azurerm_resource_group.course.name
  location            = data.azurerm_resource_group.course.location
  tags                = var.tags
}

# Tokens of the ServiceAccount db-backup in the namespace project, issued by this cluster
resource "azurerm_federated_identity_credential" "db_backup" {
  name                      = "aks-project-db-backup"
  user_assigned_identity_id = azurerm_user_assigned_identity.db_backup.id
  audience                  = ["api://AzureADTokenExchange"]
  issuer                    = azurerm_kubernetes_cluster.aks.oidc_issuer_url
  subject                   = "system:serviceaccount:project:db-backup"
}

# Write and read blobs in the backup container only
resource "azurerm_role_assignment" "db_backup_blob" {
  principal_id         = azurerm_user_assigned_identity.db_backup.principal_id
  principal_type       = "ServicePrincipal"
  role_definition_name = "Storage Blob Data Contributor"
  scope                = azurerm_storage_container.todo_backups.id
}

# The user running Terraform can list and download the backups
resource "azurerm_role_assignment" "user_backup_reader" {
  principal_id         = data.azurerm_client_config.current.object_id
  role_definition_name = "Storage Blob Data Reader"
  scope                = azurerm_storage_container.todo_backups.id
}
