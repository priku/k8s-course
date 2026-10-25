# ==============================================================================
# Container insights (exercise 3.12): the AKS counterpart of GKE's built-in monitoring.
# The Azure Monitor agent (ama-logs) on the nodes sends container logs to a Log Analytics
# workspace, where they are searched with KQL in the Azure portal (Monitor → Logs).
#
# To keep the cost down: only the namespace `project`, only the tables needed to find the
# logs (ContainerLogV2, KubePodInventory, KubeEvents), 30 days of retention (free) and a
# daily cap on ingestion (2.63 EUR/GB in Sweden Central).
# ==============================================================================
resource "azurerm_log_analytics_workspace" "aks" {
  name                = "log-dwk-aks-sc"
  resource_group_name = data.azurerm_resource_group.course.name
  location            = data.azurerm_resource_group.course.location
  sku                 = "PerGB2018"
  retention_in_days   = 30
  daily_quota_gb      = 0.2
  tags                = var.tags
}

# The Container insights data collection rule; Azure expects the name MSCI-<cluster>-<region>
resource "azurerm_monitor_data_collection_rule" "container_insights" {
  name                = "MSCI-${azurerm_kubernetes_cluster.aks.name}-${azurerm_kubernetes_cluster.aks.location}"
  resource_group_name = data.azurerm_resource_group.course.name
  location            = data.azurerm_resource_group.course.location
  tags                = var.tags

  destinations {
    log_analytics {
      name                  = "ciworkspace"
      workspace_resource_id = azurerm_log_analytics_workspace.aks.id
    }
  }

  data_flow {
    streams      = local.container_insights_streams
    destinations = ["ciworkspace"]
  }

  data_sources {
    extension {
      name           = "ContainerInsightsExtension"
      extension_name = "ContainerInsights"
      streams        = local.container_insights_streams
      extension_json = jsonencode({
        dataCollectionSettings = {
          interval               = "1m"
          namespaceFilteringMode = "Include"
          namespaces             = ["project"]
          enableContainerLogV2   = true
        }
      })
    }
  }
}

locals {
  container_insights_streams = [
    "Microsoft-ContainerLogV2", # stdout and stderr of the containers
    "Microsoft-KubePodInventory",
    "Microsoft-KubeEvents",
  ]
}

resource "azurerm_monitor_data_collection_rule_association" "container_insights" {
  name                    = "ContainerInsightsExtension"
  target_resource_id      = azurerm_kubernetes_cluster.aks.id
  data_collection_rule_id = azurerm_monitor_data_collection_rule.container_insights.id
}
