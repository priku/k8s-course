# Azure Infrastructure for DevOps with Kubernetes

This directory contains Terraform configuration for deploying Azure Kubernetes Service (AKS) and Azure Container Registry (ACR) using **Azure Verified Modules (AVM)**.

> **Note:** This course uses Azure AKS instead of GKE as suggested in the original course material.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 Azure (Sweden Central)                       │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Resource Group: dwk-aks-rg              │   │
│  │                                                      │   │
│  │  ┌─────────────────┐    ┌─────────────────────────┐ │   │
│  │  │  Azure Container │    │   Azure Kubernetes     │ │   │
│  │  │  Registry (ACR)  │◄───│   Service (AKS)        │ │   │
│  │  │  Basic SKU       │    │   2× Standard_B2s      │ │   │
│  │  └─────────────────┘    │   Azure CNI            │ │   │
│  │                          │   Azure AD RBAC        │ │   │
│  │                          └─────────────────────────┘ │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Resources Created

| Resource | Name | Description |
|----------|------|-------------|
| Resource Group | `dwk-aks-rg` | Container for all resources |
| AKS Cluster | `dwk-aks-cluster` | Managed Kubernetes cluster |
| ACR | `dwkacr<random>` | Container registry for images |
| Role Assignment | AcrPull | Allows AKS to pull images from ACR |

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.9
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) >= 2.50
- Azure subscription with sufficient permissions

## Usage

### 1. Login to Azure

```bash
az login
az account set --subscription "<your-subscription-name>"
```

### 2. Initialize Terraform

```bash
cd terraform
terraform init
```

### 3. Review the Plan

```bash
terraform plan -out=tfplan
```

### 4. Apply Infrastructure

```bash
terraform apply tfplan
```

This will take approximately **5-10 minutes**.

### 5. Configure kubectl

```bash
# Get the command from Terraform output
terraform output kube_config_command

# Or run directly:
az aks get-credentials --resource-group dwk-aks-rg --name dwk-aks-cluster --overwrite-existing
```

### 6. Verify Connection

```bash
kubectl get nodes
```

## Pushing Images to ACR

```bash
# Get ACR name
ACR_NAME=$(terraform output -raw acr_name)

# Login to ACR
az acr login --name $ACR_NAME

# Tag and push an image
docker tag myapp:latest $ACR_NAME.azurecr.io/myapp:latest
docker push $ACR_NAME.azurecr.io/myapp:latest
```

## Configuration Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `resource_group_name` | `dwk-aks-rg` | Name of the resource group |
| `location` | `swedencentral` | Azure region |
| `cluster_name` | `dwk-aks-cluster` | Name of the AKS cluster |
| `node_count` | `2` | Number of nodes |
| `vm_size` | `Standard_B2s` | VM size for nodes |

## Outputs

| Output | Description |
|--------|-------------|
| `resource_group_name` | Name of the created resource group |
| `kubernetes_cluster_name` | Name of the AKS cluster |
| `kube_config_command` | Command to get kubeconfig |
| `acr_login_server` | ACR login server URL |
| `acr_name` | ACR name |
| `cluster_portal_url` | Azure Portal URL for cluster |
| `acr_portal_url` | Azure Portal URL for ACR |

## Clean Up

To destroy all resources:

```bash
terraform destroy
```

## Azure Verified Modules

This configuration uses [Azure Verified Modules (AVM)](https://azure.github.io/Azure-Verified-Modules/) for best practices:

- **AKS Module**: `Azure/avm-res-containerservice-managedcluster/azurerm` v0.3.3
- **ACR Module**: `Azure/avm-res-containerregistry-registry/azurerm` v0.4.0

## Cost Estimate

| Resource | Estimated Monthly Cost |
|----------|----------------------|
| AKS (2× Standard_B2s) | ~$60-80 |
| ACR (Basic) | ~$5 |
| Load Balancer | ~$20 |
| **Total** | **~$85-105/month** |

> **Tip:** Remember to `terraform destroy` when not using the cluster to avoid charges!

## Troubleshooting

### "subscription ID could not be determined"
```bash
az login
az account set --subscription "<subscription-name>"
```

### "quota exceeded"
Try a different region or request quota increase in Azure Portal.

### Pods can't pull images from ACR
Verify the role assignment exists:
```bash
az aks check-acr --resource-group dwk-aks-rg --name dwk-aks-cluster --acr dwkacr<suffix>.azurecr.io
```
