output "resource_group_name" {
  description = "Resource group name"
  value       = azurerm_resource_group.aks.name
}

output "kubernetes_cluster_name" {
  description = "AKS cluster name"
  value       = module.aks.name
}

output "kube_config_command" {
  description = "Command to get kubeconfig"
  value       = "az aks get-credentials --resource-group ${azurerm_resource_group.aks.name} --name ${module.aks.name} --overwrite-existing"
}

output "acr_login_server" {
  description = "ACR login server"
  value       = module.acr.resource.login_server
}

output "acr_name" {
  description = "ACR name"
  value       = module.acr.name
}

output "cluster_portal_url" {
  description = "Azure Portal URL for the cluster"
  value       = "https://portal.azure.com/#resource${module.aks.resource_id}/overview"
}

output "acr_portal_url" {
  description = "Azure Portal URL for the container registry"
  value       = "https://portal.azure.com/#resource${module.acr.resource_id}/overview"
}
