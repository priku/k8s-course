# Azure AI Services Setup Guide
## For AI-Enhanced Todo Application

This guide walks you through setting up all necessary Azure AI resources for the project.

---

## Prerequisites

- Azure account (create free account at https://azure.microsoft.com/free/)
- Azure CLI installed: `az --version`
- Azure subscription with AI-102 credits or free tier

---

## Quick Start (5 minutes)

### 1. Login to Azure
```bash
az login
az account show
```

### 2. Create Resource Group
```bash
# Set variables
RESOURCE_GROUP="rg-todo-ai-app"
LOCATION="eastus"  # or "westeurope" for Europe

# Create resource group
az group create \
  --name $RESOURCE_GROUP \
  --location $LOCATION
```

### 3. Create Azure AI Services (Multi-Service)
```bash
AI_SERVICE_NAME="ai-todo-app-$(date +%s)"

az cognitiveservices account create \
  --name $AI_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --kind CognitiveServices \
  --sku S0 \
  --location $LOCATION \
  --yes
```

### 4. Get Keys and Endpoint
```bash
# Get primary key
AI_KEY=$(az cognitiveservices account keys list \
  --name $AI_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "key1" -o tsv)

# Get endpoint
AI_ENDPOINT=$(az cognitiveservices account show \
  --name $AI_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "properties.endpoint" -o tsv)

echo "AI Service Key: $AI_KEY"
echo "AI Service Endpoint: $AI_ENDPOINT"
```

### 5. Create Azure OpenAI Resource
```bash
OPENAI_NAME="openai-todo-app-$(date +%s)"

az cognitiveservices account create \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP \
  --kind OpenAI \
  --sku S0 \
  --location $LOCATION \
  --yes

# Get OpenAI keys
OPENAI_KEY=$(az cognitiveservices account keys list \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "key1" -o tsv)

OPENAI_ENDPOINT=$(az cognitiveservices account show \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "properties.endpoint" -o tsv)

echo "OpenAI Key: $OPENAI_KEY"
echo "OpenAI Endpoint: $OPENAI_ENDPOINT"
```

### 6. Deploy GPT-4 Model
```bash
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

# Alternative: Deploy GPT-3.5-Turbo (cheaper for testing)
az cognitiveservices account deployment create \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP \
  --deployment-name gpt-35-turbo \
  --model-name gpt-35-turbo \
  --model-version "0613" \
  --model-format OpenAI \
  --sku-capacity 10 \
  --sku-name "Standard"
```

---

## Detailed Service Setup

### Azure Language Service

#### Create Language Resource
```bash
LANGUAGE_SERVICE_NAME="lang-todo-app-$(date +%s)"

az cognitiveservices account create \
  --name $LANGUAGE_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --kind TextAnalytics \
  --sku S \
  --location $LOCATION \
  --yes

# Get keys
LANGUAGE_KEY=$(az cognitiveservices account keys list \
  --name $LANGUAGE_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "key1" -o tsv)

LANGUAGE_ENDPOINT=$(az cognitiveservices account show \
  --name $LANGUAGE_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "properties.endpoint" -o tsv)
```

#### Features Available
- Sentiment Analysis
- Key Phrase Extraction
- Named Entity Recognition
- Language Detection
- Opinion Mining

### Azure Computer Vision

#### Create Vision Resource
```bash
VISION_SERVICE_NAME="vision-todo-app-$(date +%s)"

az cognitiveservices account create \
  --name $VISION_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --kind ComputerVision \
  --sku S1 \
  --location $LOCATION \
  --yes

# Get keys
VISION_KEY=$(az cognitiveservices account keys list \
  --name $VISION_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "key1" -o tsv)

VISION_ENDPOINT=$(az cognitiveservices account show \
  --name $VISION_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "properties.endpoint" -o tsv)
```

#### Features Available
- Image Analysis
- OCR (Optical Character Recognition)
- Object Detection
- Face Detection
- Image Tagging

### Azure Speech Service (Optional)

#### Create Speech Resource
```bash
SPEECH_SERVICE_NAME="speech-todo-app-$(date +%s)"

az cognitiveservices account create \
  --name $SPEECH_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --kind SpeechServices \
  --sku S0 \
  --location $LOCATION \
  --yes

# Get keys
SPEECH_KEY=$(az cognitiveservices account keys list \
  --name $SPEECH_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP \
  --query "key1" -o tsv)

SPEECH_REGION=$LOCATION
```

---

## Save Credentials to Kubernetes Secrets

### Create Secret YAML
```bash
cat > manifests/ai-service/secret.yaml <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: azure-ai-secrets
  namespace: default
type: Opaque
stringData:
  # Azure OpenAI
  AZURE_OPENAI_KEY: "$OPENAI_KEY"
  AZURE_OPENAI_ENDPOINT: "$OPENAI_ENDPOINT"
  AZURE_OPENAI_DEPLOYMENT: "gpt-4"

  # Azure Language Service
  AZURE_LANGUAGE_KEY: "$LANGUAGE_KEY"
  AZURE_LANGUAGE_ENDPOINT: "$LANGUAGE_ENDPOINT"

  # Azure Computer Vision
  AZURE_VISION_KEY: "$VISION_KEY"
  AZURE_VISION_ENDPOINT: "$VISION_ENDPOINT"

  # Azure Speech (optional)
  AZURE_SPEECH_KEY: "$SPEECH_KEY"
  AZURE_SPEECH_REGION: "$SPEECH_REGION"
EOF

# Apply to cluster
kubectl apply -f manifests/ai-service/secret.yaml
```

### Verify Secret
```bash
kubectl get secret azure-ai-secrets
kubectl describe secret azure-ai-secrets
```

---

## Environment Variables for Local Development

### Create .env file
```bash
cat > services/ai-service/.env <<EOF
# Azure OpenAI
AZURE_OPENAI_KEY=$OPENAI_KEY
AZURE_OPENAI_ENDPOINT=$OPENAI_ENDPOINT
AZURE_OPENAI_DEPLOYMENT=gpt-4

# Azure Language Service
AZURE_LANGUAGE_KEY=$LANGUAGE_KEY
AZURE_LANGUAGE_ENDPOINT=$LANGUAGE_ENDPOINT

# Azure Computer Vision
AZURE_VISION_KEY=$VISION_KEY
AZURE_VISION_ENDPOINT=$VISION_ENDPOINT

# Azure Speech
AZURE_SPEECH_KEY=$SPEECH_KEY
AZURE_SPEECH_REGION=$SPEECH_REGION
EOF

# Add to .gitignore
echo "**/.env" >> .gitignore
```

---

## Cost Management

### Set Up Budget Alerts
```bash
# Create budget (e.g., $50/month)
az consumption budget create \
  --budget-name "todo-app-monthly-budget" \
  --amount 50 \
  --time-grain Monthly \
  --start-date $(date +%Y-%m-01) \
  --end-date $(date -d "+1 year" +%Y-%m-01) \
  --resource-group $RESOURCE_GROUP
```

### Monitor Costs
```bash
# View current month costs
az consumption usage list \
  --start-date $(date +%Y-%m-01) \
  --end-date $(date +%Y-%m-%d)
```

### Cost Optimization Tips
1. **Use GPT-3.5-Turbo instead of GPT-4** for development (10x cheaper)
2. **Implement caching** to reduce API calls
3. **Set rate limits** in your application
4. **Use free tier** for Computer Vision and Language Service when possible
5. **Delete resources** when not in use

### Estimated Monthly Costs (Moderate Use)
| Service | Usage | Cost |
|---------|-------|------|
| Azure OpenAI (GPT-4) | 100k tokens/day | $30 |
| Language Service | 5k requests/day | $8 |
| Computer Vision | 1k images/day | $5 |
| **Total** | | **~$43/month** |

### Free Tier Limits
- Language Service: 5,000 text records/month FREE
- Computer Vision: 5,000 transactions/month FREE
- Speech: 5 hours/month FREE

---

## Testing the Setup

### Test Azure OpenAI
```bash
curl $OPENAI_ENDPOINT/openai/deployments/gpt-4/chat/completions?api-version=2023-05-15 \
  -H "Content-Type: application/json" \
  -H "api-key: $OPENAI_KEY" \
  -d '{
    "messages": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

### Test Language Service
```bash
curl -X POST "$LANGUAGE_ENDPOINT/text/analytics/v3.1/sentiment" \
  -H "Ocp-Apim-Subscription-Key: $LANGUAGE_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "id": "1",
        "language": "en",
        "text": "This todo app is amazing!"
      }
    ]
  }'
```

### Test Computer Vision
```bash
curl -X POST "$VISION_ENDPOINT/vision/v3.2/analyze?visualFeatures=Description,Tags" \
  -H "Ocp-Apim-Subscription-Key: $VISION_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://upload.wikimedia.org/wikipedia/commons/thumb/1/12/User_icon_2.svg/220px-User_icon_2.svg.png"
  }'
```

---

## Cleanup Resources

### Delete Individual Resources
```bash
# Delete OpenAI resource
az cognitiveservices account delete \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP

# Delete Language Service
az cognitiveservices account delete \
  --name $LANGUAGE_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP

# Delete Computer Vision
az cognitiveservices account delete \
  --name $VISION_SERVICE_NAME \
  --resource-group $RESOURCE_GROUP
```

### Delete Entire Resource Group (⚠️ Careful!)
```bash
az group delete --name $RESOURCE_GROUP --yes --no-wait
```

---

## Troubleshooting

### Issue: "Deployment not available in region"
**Solution**: Try different regions
```bash
# List available regions for OpenAI
az cognitiveservices account list-skus \
  --kind OpenAI \
  --query "[].locations" -o table

# Common regions with OpenAI: eastus, westeurope, southcentralus
```

### Issue: "Quota exceeded"
**Solution**: Request quota increase
```bash
# Check current quota
az cognitiveservices account list-usage \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP

# Request increase via Azure Portal
```

### Issue: "Authentication failed"
**Solution**: Regenerate keys
```bash
az cognitiveservices account keys regenerate \
  --name $OPENAI_NAME \
  --resource-group $RESOURCE_GROUP \
  --key-name key1
```

---

## Next Steps

1. ✅ Complete this Azure setup
2. [ ] Test all endpoints with curl
3. [ ] Create Kubernetes secrets
4. [ ] Build Python AI service
5. [ ] Integrate with API gateway

---

## Useful Links

- [Azure AI Services Documentation](https://learn.microsoft.com/en-us/azure/ai-services/)
- [Azure OpenAI Service](https://learn.microsoft.com/en-us/azure/ai-services/openai/)
- [Azure CLI Reference](https://learn.microsoft.com/en-us/cli/azure/)
- [AI-102 Certification](https://learn.microsoft.com/en-us/credentials/certifications/azure-ai-engineer/)
- [Pricing Calculator](https://azure.microsoft.com/en-us/pricing/calculator/)

---

**Created**: 2025-11-29
**Last Updated**: 2025-11-29
