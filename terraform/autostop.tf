# ==============================================================================
# Nightly stop: a Logic App stops the cluster every night at 23:59 Finnish time,
# so that the nodes do not run (and cost) when nobody is working.
# Nothing starts it again; run `az aks start` when needed.
# ==============================================================================
resource "azurerm_logic_app_workflow" "aks_nightly_stop" {
  name                = "logic-${var.cluster_name}-nightly-stop"
  resource_group_name = data.azurerm_resource_group.course.name
  location            = data.azurerm_resource_group.course.location
  tags                = var.tags

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_logic_app_trigger_recurrence" "nightly" {
  name         = "every-night-2359"
  logic_app_id = azurerm_logic_app_workflow.aks_nightly_stop.id
  frequency    = "Day"
  interval     = 1
  time_zone    = var.stop_time_zone

  schedule {
    at_these_hours   = [23]
    at_these_minutes = [59]
  }
}

# Calls the AKS stop operation of Azure Resource Manager with the Logic App's managed identity
resource "azurerm_logic_app_action_custom" "stop_cluster" {
  name         = "stop-aks-cluster"
  logic_app_id = azurerm_logic_app_workflow.aks_nightly_stop.id

  body = jsonencode({
    type = "Http"
    inputs = {
      method = "POST"
      uri    = "https://management.azure.com${azurerm_kubernetes_cluster.aks.id}/stop?api-version=2026-08-01"
      authentication = {
        type     = "ManagedServiceIdentity"
        audience = "https://management.azure.com/"
      }
    }
    runAfter = {}
  })

  depends_on = [azurerm_logic_app_trigger_recurrence.nightly]
}

# The Logic App may manage this one cluster only
resource "azurerm_role_assignment" "logic_app_aks" {
  principal_id         = azurerm_logic_app_workflow.aks_nightly_stop.identity[0].principal_id
  role_definition_name = "Azure Kubernetes Service Contributor Role"
  scope                = azurerm_kubernetes_cluster.aks.id
}
