# Azure infrastructure (Chapter 4)

The course material uses Google Kubernetes Engine in Chapter 4, and also allows other cloud providers. These exercises run on **Azure Kubernetes Service (AKS)** instead. This Terraform configuration creates the infrastructure in the existing resource group `rg-dwk-aks-sc` (Sweden Central).

| Resource | Name | Notes |
|---|---|---|
| AKS cluster | `aks-dwk-sc` | Kubernetes 1.37, Free tier, 2 × `Standard_D2as_v7` (2 vCPU, 8 GiB), Azure CNI Overlay |
| Node resource group | `rg-dwk-aks-nodes-sc` | Created by AKS for the VMs, disks, load balancer and public IPs |
| Container registry | `acrdwkakssc.azurecr.io` | Basic SKU, no admin user; the nodes pull with their kubelet identity (`AcrPull`) |
| Logic App | `logic-aks-dwk-sc-nightly-stop` | Stops the cluster every night at 23:59 Finnish time |
| Gateway API | GatewayClass `approuting-istio` | AKS-managed Gateway API CRDs and the application routing Gateway API implementation (`gateway.tf`, exercise 3.3) |

Sign-in to the cluster goes through Entra ID with Azure RBAC (`local_account_disabled`), so there are no static admin credentials. The user who runs Terraform gets `Azure Kubernetes Service RBAC Cluster Admin` on the cluster and `AcrPush` on the registry.

## GKE in the material vs. AKS here

| Course material (GKE) | Here (AKS) |
|---|---|
| `gcloud container clusters create` | `terraform apply` in this directory |
| `gcloud container clusters get-credentials` | `az aks get-credentials` + `kubelogin` |
| Google Artifact Registry | Azure Container Registry, images built with `az acr build` |
| Persistent disk from the default storage class | Azure Disk from the default storage class `default` |
| `Service` of type `LoadBalancer` | The same; AKS adds a rule and a public IP to the cluster's Azure Load Balancer |
| GKE Ingress (Google L7 load balancer) | Traefik installed with Helm (exercise 3.2); AKS's managed NGINX is being retired |
| `gcloud container clusters update --gateway-api=standard`, GatewayClass `gke-l7-global-external-managed` | `gateway.tf` (azapi), GatewayClass `approuting-istio`: an Envoy proxy in the cluster behind an Azure Load Balancer IP |
| Delete the cluster to save credits | Stop it: `az aks stop` (nightly automatically), `az aks start` when needed |

## One-time setup

These were created once with the Azure CLI, outside Terraform, so that `terraform destroy` never removes the budget or the Terraform state:

```bash
# Resource group for everything of the course
az group create --name rg-dwk-aks-sc --location swedencentral \
  --tags project=dwk-k8s-course purpose=course owner=priku

# Budget: 30 EUR/month for rg-dwk-aks-sc and the node group rg-dwk-aks-nodes-sc,
# alerts at 50/80/100 % actual and 100 % forecast (created with az rest, Microsoft.Consumption/budgets)

# Terraform state in Azure Storage, read with Entra ID (shared keys disabled)
az storage account create --resource-group rg-dwk-aks-sc --name stdwkakstfstatesc \
  --location swedencentral --sku Standard_LRS --min-tls-version TLS1_2 \
  --allow-blob-public-access false --allow-shared-key-access false
az role assignment create --assignee-object-id "$(az ad signed-in-user show --query id -o tsv)" \
  --assignee-principal-type User --role "Storage Blob Data Contributor" \
  --scope "$(az storage account show -g rg-dwk-aks-sc -n stdwkakstfstatesc --query id -o tsv)"
az storage container create --name tfstate --account-name stdwkakstfstatesc --auth-mode login
```

## Create the infrastructure

```bash
az login
export ARM_SUBSCRIPTION_ID=$(az account show --query id -o tsv)   # the subscription is not in the code
terraform init
terraform plan -out=tfplan
terraform apply tfplan

# kubeconfig for the cluster; kubelogin uses the Azure CLI sign-in
$(terraform output -raw get_credentials_command)
kubectl get nodes
```

Tools: Terraform 1.17, azurerm provider 5.11, azapi provider 2.13 (for the Gateway API settings that azurerm does not cover yet), Azure CLI 2.90, kubelogin 0.2.21.

## Stop and start

The cluster costs money only while its nodes run. A Logic App stops it every night at 23:59 Finnish time (`FLE Standard Time`, so daylight saving time is followed). Its managed identity has only the `Azure Kubernetes Service Contributor Role` on this cluster. Nothing starts the cluster automatically:

```bash
az aks start --resource-group rg-dwk-aks-sc --name aks-dwk-sc   # about 5 minutes
az aks stop  --resource-group rg-dwk-aks-sc --name aks-dwk-sc
az aks show  --resource-group rg-dwk-aks-sc --name aks-dwk-sc --query powerState.code -o tsv
```

A stopped cluster keeps its configuration, Kubernetes objects and persistent disks, so nothing has to be applied again after `az aks start`.

## Costs

| Item | Price (Sweden Central) |
|---|---|
| 2 × `Standard_D2as_v7` | 0.085 EUR/h each, only while running |
| Container registry Basic | about 0.15 EUR/day |
| Public IPs, persistent disks, state storage | a few euros per month |
| AKS control plane (Free tier) | free |

Cheaper B-series sizes are out: the subscription has no quota for the Bsv2 family, and AKS does not allow the older B-series in Sweden Central. `Standard_D2als_v7` (4 GiB) would be about 10 % cheaper, but the Istio gateway of the later exercises and two Postgres databases need the memory of `D2as_v7`.

## Remove

```bash
terraform destroy                                  # cluster, registry, Logic App
az group delete --name rg-dwk-aks-sc               # at the end of the course: state storage and the group
```
