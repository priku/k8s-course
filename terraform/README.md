# Azure Infrastructure for DevOps with Kubernetes

This directory contains Terraform configuration for deploying Azure Kubernetes Service (AKS) and Azure Container Registry (ACR) using **Azure Verified Modules (AVM)**.

> **Note:** This course uses Azure AKS instead of GKE as suggested in the original course material.

## Table of Contents

- [Architecture](#architecture)
- [AKS Infrastructure](#aks-infrastructure)
- [Usage](#usage)
- [Pushing Images to ACR](#pushing-images-to-acr)
- [Configuration Variables](#configuration-variables)
- [Outputs](#outputs)
- [Azure AI Services (Optional)](#azure-ai-services-optional)
- [Cost Management](#cost-management)
- [Clean Up](#clean-up)
- [Troubleshooting](#troubleshooting)

---

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
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │         Azure AI Services (Optional/Future)          │   │
│  │  ┌─────────┐  ┌──────────┐  ┌─────────────────────┐ │   │
│  │  │ OpenAI  │  │ Language │  │  Computer Vision    │ │   │
│  │  │ GPT-4   │  │ Service  │  │  OCR, Analysis      │ │   │
│  │  └─────────┘  └──────────┘  └─────────────────────┘ │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## AKS Infrastructure

### Resources Created

| Resource | Name | Description |
|----------|------|-------------|
| Resource Group | `dwk-aks-rg` | Container for all resources |
| AKS Cluster | `dwk-aks-cluster` | Managed Kubernetes cluster |
| ACR | `dwkacr<random>` | Container registry for images |
| Role Assignment | AcrPull | Allows AKS to pull images from ACR |

### Azure Verified Modules

This configuration uses [Azure Verified Modules (AVM)](https://azure.github.io/Azure-Verified-Modules/) for best practices:

- **AKS Module**: `Azure/avm-res-containerservice-managedcluster/azurerm` v0.3.3
- **ACR Module**: `Azure/avm-res-containerregistry-registry/azurerm` v0.4.0

---

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.9
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) >= 2.50
- Azure subscription with sufficient permissions

---

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

---

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

---

## Configuration Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `resource_group_name` | `dwk-aks-rg` | Name of the resource group |
| `location` | `swedencentral` | Azure region |
| `cluster_name` | `dwk-aks-cluster` | Name of the AKS cluster |
| `node_count` | `2` | Number of nodes |
| `vm_size` | `Standard_B2s` | VM size for nodes |

---

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

---

## Azure AI Services (Optional)

> **Note:** This section is for future AI-102 integration. Not required for the DevOps with Kubernetes course.

### Quick Setup

```bash
# Set variables
RESOURCE_GROUP="rg-todo-ai-app"
LOCATION="eastus"

# Create resource group for AI services
az group create --name $RESOURCE_GROUP --location $LOCATION

# Create Azure OpenAI
OPENAI_NAME="openai-todo-app-$(date +%s)"
az cognitiveservices account create \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP \
  --kind OpenAI \
  --sku S0 \
  --location $LOCATION \
  --yes

# Deploy GPT-4 model
az cognitiveservices account deployment create \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP \
  --deployment-name gpt-4 \
  --model-name gpt-4 \
  --model-version "0613" \
  --model-format OpenAI \
  --sku-capacity 10 \
  --sku-name "Standard"
```

### Available AI Services

| Service | Use Case | Free Tier |
|---------|----------|-----------|
| Azure OpenAI | Smart todo parsing, AI assistant | No free tier |
| Language Service | Sentiment analysis, key phrases | 5,000 records/month |
| Computer Vision | OCR, image analysis | 5,000 transactions/month |
| Speech Service | Voice input/output | 5 hours/month |

### Create Kubernetes Secret for AI Services

```bash
cat > ai-secrets.yaml <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: azure-ai-secrets
type: Opaque
stringData:
  AZURE_OPENAI_KEY: "<your-key>"
  AZURE_OPENAI_ENDPOINT: "<your-endpoint>"
  AZURE_OPENAI_DEPLOYMENT: "gpt-4"
EOF

kubectl apply -f ai-secrets.yaml
```

---

## Cost Management

### AKS Infrastructure Costs

| Resource | Estimated Monthly Cost |
|----------|----------------------|
| AKS (2× Standard_B2s) | ~$60-80 |
| ACR (Basic) | ~$5 |
| Load Balancer | ~$20 |
| **Total** | **~$85-105/month** |

### AI Services Costs (Optional)

| Service | Usage | Cost |
|---------|-------|------|
| Azure OpenAI (GPT-4) | 100k tokens/day | ~$30/month |
| Language Service | 5k requests/day | ~$8/month |
| Computer Vision | 1k images/day | ~$5/month |

### Cost Optimization Tips

1. **Destroy when not in use**: `terraform destroy`
2. **Use GPT-3.5-Turbo** instead of GPT-4 (10x cheaper)
3. **Stay within free tiers** for Language/Vision services
4. **Set budget alerts** in Azure Portal

### Set Up Budget Alert

```bash
az consumption budget create \
  --budget-name "k8s-course-budget" \
  --amount 100 \
  --time-grain Monthly \
  --start-date $(date +%Y-%m-01) \
  --end-date $(date -d "+1 year" +%Y-%m-01) \
  --resource-group dwk-aks-rg
```

---

## Clean Up

### Destroy AKS Infrastructure

```bash
cd terraform
terraform destroy
```

### Delete AI Services (if created)

```bash
az group delete --name rg-todo-ai-app --yes --no-wait
```

---

## Troubleshooting

### "subscription ID could not be determined"
```bash
az login
az account set --subscription "<subscription-name>"
```

### "quota exceeded"
Try a different region or request quota increase in Azure Portal.

### Pods can't pull images from ACR
```bash
az aks check-acr --resource-group dwk-aks-rg --name dwk-aks-cluster --acr <acr-name>.azurecr.io
```

### "Deployment not available in region" (AI Services)
```bash
# List available regions for OpenAI
az cognitiveservices account list-skus --kind OpenAI --query "[].locations" -o table
```

---

## Useful Links

- [Azure Verified Modules](https://azure.github.io/Azure-Verified-Modules/)
- [AKS Documentation](https://learn.microsoft.com/en-us/azure/aks/)
- [Azure AI Services](https://learn.microsoft.com/en-us/azure/ai-services/)
- [Terraform AzureRM Provider](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs)
- [AI-102 Certification](https://learn.microsoft.com/en-us/credentials/certifications/azure-ai-engineer/)
