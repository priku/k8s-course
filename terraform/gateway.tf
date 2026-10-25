# ==============================================================================
# Gateway API (exercise 3.3): the AKS-managed Gateway API CRDs and the application
# routing Gateway API implementation, which gives the cluster the GatewayClass
# `approuting-istio` (a lightweight Istio control plane that only handles Gateway API
# resources, no service mesh). This is the AKS counterpart of the GKE gateway classes
# in the course material. azurerm has no arguments for these yet, so azapi sets them.
# ==============================================================================
resource "azapi_update_resource" "gateway_api" {
  type        = "Microsoft.ContainerService/managedClusters@2026-08-01"
  resource_id = azurerm_kubernetes_cluster.aks.id

  body = {
    properties = {
      ingressProfile = {
        gatewayAPI = {
          installation = "Standard"
        }
        webAppRouting = {
          gatewayAPIImplementations = {
            appRoutingIstio = {
              mode = "Enabled"
            }
          }
        }
      }
    }
  }

  response_export_values = ["properties.ingressProfile"]
}
