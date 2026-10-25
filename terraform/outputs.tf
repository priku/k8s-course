output "cluster_name" {
  value = azurerm_kubernetes_cluster.aks.name
}

output "acr_login_server" {
  value = azurerm_container_registry.acr.login_server
}

output "get_credentials_command" {
  value = "az aks get-credentials --resource-group ${data.azurerm_resource_group.course.name} --name ${azurerm_kubernetes_cluster.aks.name} --overwrite-existing && kubelogin convert-kubeconfig -l azurecli"
}

output "oidc_issuer_url" {
  value = azurerm_kubernetes_cluster.aks.oidc_issuer_url
}
