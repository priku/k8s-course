# DevOps with Kubernetes 2026

Exercise submissions for [DevOps with Kubernetes](https://courses.mooc.fi/org/uh-cs/courses/devops-with-kubernetes-2026) (University of Helsinki). Each app lives in its own directory with its own README, and every exercise has a GitHub release tagged with the exercise number.

## Exercises

### Chapter 2

- [1.1.](https://github.com/priku/k8s-course/tree/1.1/log-output)
- [1.2.](https://github.com/priku/k8s-course/tree/1.2/todo-project)
- [1.3.](https://github.com/priku/k8s-course/tree/1.3/log-output)
- [1.4.](https://github.com/priku/k8s-course/tree/1.4/todo-project)
- [1.5.](https://github.com/priku/k8s-course/tree/1.5/todo-project)
- [1.6.](https://github.com/priku/k8s-course/tree/1.6/todo-project)
- [1.7.](https://github.com/priku/k8s-course/tree/1.7/log-output)
- [1.8.](https://github.com/priku/k8s-course/tree/1.8/todo-project)
- [1.9.](https://github.com/priku/k8s-course/tree/1.9/ping-pong)
- [1.10.](https://github.com/priku/k8s-course/tree/1.10/log-output)
- [1.11.](https://github.com/priku/k8s-course/tree/1.11/log-output) ([ping-pong](https://github.com/priku/k8s-course/tree/1.11/ping-pong))
- [1.12.](https://github.com/priku/k8s-course/tree/1.12/todo-project)
- [1.13.](https://github.com/priku/k8s-course/tree/1.13/todo-project)

### Chapter 3

- [2.1.](https://github.com/priku/k8s-course/tree/2.1/log-output) ([ping-pong](https://github.com/priku/k8s-course/tree/2.1/ping-pong))
- [2.2.](https://github.com/priku/k8s-course/tree/2.2/todo-backend)
- [2.3.](https://github.com/priku/k8s-course/tree/2.3/log-output) ([ping-pong](https://github.com/priku/k8s-course/tree/2.3/ping-pong))
- [2.4.](https://github.com/priku/k8s-course/tree/2.4/todo-project) ([todo-backend](https://github.com/priku/k8s-course/tree/2.4/todo-backend))
- [2.5.](https://github.com/priku/k8s-course/tree/2.5/log-output)
- [2.6.](https://github.com/priku/k8s-course/tree/2.6/todo-project)
- [2.7.](https://github.com/priku/k8s-course/tree/2.7/ping-pong)
- [2.8.](https://github.com/priku/k8s-course/tree/2.8/todo-backend)
- [2.9.](https://github.com/priku/k8s-course/tree/2.9/wiki-todo-generator)
- [2.10.](https://github.com/priku/k8s-course/tree/2.10/todo-backend)

### Chapter 4

On Azure Kubernetes Service instead of Google Kubernetes Engine (chapter 1 of the course allows other cloud providers). The cluster is in [`terraform`](https://github.com/priku/k8s-course/tree/3.4/terraform).

- [3.1.](https://github.com/priku/k8s-course/tree/3.1/ping-pong) ([terraform](https://github.com/priku/k8s-course/tree/3.1/terraform))
- [3.2.](https://github.com/priku/k8s-course/tree/3.2/ping-pong) ([log-output](https://github.com/priku/k8s-course/tree/3.2/log-output))
- [3.3.](https://github.com/priku/k8s-course/tree/3.3/ping-pong) ([terraform](https://github.com/priku/k8s-course/tree/3.3/terraform))
- [3.4.](https://github.com/priku/k8s-course/tree/3.4/ping-pong)
- [3.5.](https://github.com/priku/k8s-course/tree/3.5/project)
- [3.6.](https://github.com/priku/k8s-course/tree/3.6/project) ([workflow](https://github.com/priku/k8s-course/blob/3.6/.github/workflows/project.yaml))
- [3.7.](https://github.com/priku/k8s-course/tree/3.7/project) ([workflow](https://github.com/priku/k8s-course/blob/3.7/.github/workflows/project.yaml))

---

## About the project: AI-Enhanced Todo Application
### DevOps with Kubernetes + Azure AI Services (AI-102)

[![Terraform](https://github.com/priku/k8s-course/actions/workflows/terraform.yml/badge.svg)](https://github.com/priku/k8s-course/actions/workflows/terraform.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A production-grade, cloud-native todo application demonstrating Kubernetes orchestration and Azure AI capabilities. This project combines learning from:
- **DevOps with Kubernetes** course (University of Helsinki, 5 ECTS)
- **AI-102: Azure AI Engineer** certification path

**🔧 Built with Professional DevOps Practices:**
- ✅ GitOps workflow with Pull Requests
- ✅ CI/CD with GitHub Actions
- ✅ Infrastructure as Code (Terraform)
- ✅ Remote state management (Azure Storage)
- ✅ Branch protection and code review
- ✅ Automated testing and validation

---

## 🎯 Project Goals

1. Master Kubernetes concepts through hands-on development
2. Integrate Azure AI services (OpenAI, Language, Computer Vision)
3. Build a portfolio-worthy, production-ready application
4. Practice DevOps best practices (CI/CD, monitoring, security)

---

## 📋 Documentation

- **[PROJECT_PLAN.md](PROJECT_PLAN.md)** - Complete architecture and roadmap
- **[terraform/README.md](terraform/README.md)** - Azure infrastructure setup (AKS, ACR, AI Services)

---

## 🏗️ Architecture (Planned)

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Frontend   │→ │   API        │→ │  AI Services │      │
│  │   (React)    │  │  Gateway     │  │  Container   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                              ↓               │
│                                    ┌──────────────────┐     │
│                                    │  Azure AI APIs   │     │
│                                    │  - OpenAI GPT-4  │     │
│                                    │  - Language      │     │
│                                    │  - Vision        │     │
│                                    └──────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

---

## 💡 AI Features (Planned)

### 1. Smart Todo Creation
Using Azure OpenAI to parse natural language:
```
"Buy groceries tomorrow at 3pm"
→ Task: "Buy groceries", Due: Tomorrow 15:00, Priority: Medium
```

### 2. Sentiment-Based Priority
Using Language Service for automatic prioritization:
```
"URGENT! Fix production bug!!!"
→ Sentiment: Negative → Priority: HIGH
```

### 3. Image Analysis
Using Computer Vision for attachments:
```
Upload receipt image
→ OCR + Object Detection → Auto-create todo with extracted info
```

### 4. AI Assistant
Using GPT-4 for intelligent task recommendations

---

## 📁 Current Structure

```
k8s-course/
├── log-output/              # Exercises 1.1-2.5
│   ├── main.go
│   ├── Dockerfile
│   └── manifests/
│       └── deployment.yaml
│
├── todo-project/            # Exercises 1.2-2.10
│   ├── main.go
│   ├── Dockerfile
│   └── manifests/
│       └── deployment.yaml
│
├── ping-pong/               # Exercises 1.9-2.7
│   ├── main.go
│   ├── Dockerfile
│   └── manifests/
│
├── todo-backend/            # Exercises 2.2-2.10
│   ├── main.go
│   ├── Dockerfile
│   └── manifests/
│
├── terraform/               # Azure AKS Infrastructure (Part 3+)
│   ├── main.tf              # AKS & ACR using Azure Verified Modules
│   ├── providers.tf         # Azure provider configuration
│   ├── variables.tf         # Configurable parameters
│   ├── outputs.tf           # Useful outputs
│   ├── README.md            # Azure setup guide (AKS + AI services)
│   └── .gitignore           # Exclude state files
│
├── PROJECT_PLAN.md          # Complete project roadmap
├── SUBMISSION_GUIDE.md      # Course submission guide
└── README.md                # This file
```

---

## ☁️ Azure Infrastructure (Part 3+)

Chapter 4 runs on **Azure Kubernetes Service (AKS)** instead of GKE, created with **Terraform** in the resource group `rg-dwk-aks-sc` (Sweden Central):

- AKS 1.37, Free tier, 2 × `Standard_D2as_v7`, Entra ID sign-in with Azure RBAC
- Azure Container Registry `acrdwkakssc`, images built with `az acr build`
- Gateway API with the AKS-managed GatewayClass `approuting-istio`
- A Logic App that stops the cluster every night at 23:59 Finnish time, and a 30 EUR/month budget alert

Setup, costs and the GKE → AKS mapping are in [`terraform/README.md`](https://github.com/priku/k8s-course/tree/3.4/terraform).

Exercises 3.1–3.8 were first done on AKS in 2025. 3.1–3.7 are redone for the 2026 course and their tags point to the new versions; the tag 3.8 still holds the 2025 versions until they are redone.

---

## 🛠️ Technology Stack

### Current
- **Language**: Go 1.27
- **Container**: Docker
- **Orchestration**: Kubernetes (k3d local, AKS cloud)
- **Infrastructure**: Terraform with Azure Verified Modules
- **Cloud**: Azure (AKS, ACR)
- **Tools**: kubectl, k3d, Azure CLI

### Planned
- **Frontend**: React + TypeScript + Tailwind CSS
- **Backend**: Go (API Gateway), Python (AI Services)
- **Database**: PostgreSQL
- **AI**: Azure OpenAI, Language Service, Computer Vision
- **Monitoring**: Prometheus + Grafana
- **CI/CD**: GitHub Actions

---

## 🎓 Learning Path

### Phase 1: Kubernetes Foundations ✅
- [x] Basic deployments and pods
- [x] Services and networking
- [x] ConfigMaps and Secrets
- [x] Persistent storage (PersistentVolumeClaims)
- [x] StatefulSets (PostgreSQL)
- [x] CronJobs

### Phase 2: Cloud Deployment (Current)
- [x] Azure infrastructure with Terraform
- [ ] Deploy to Azure Kubernetes Service
- [ ] Azure Container Registry integration
- [ ] Cloud-native services

### Phase 3: AI Integration (Planned)
- [ ] Set up Azure AI resources
- [ ] Build Python AI service
- [ ] Integrate with API gateway
- [ ] Deploy to Kubernetes

### Phase 4: Production Ready
- [ ] Monitoring and logging
- [ ] CI/CD pipeline
- [ ] Security hardening
- [ ] Performance optimization

---

## 🚦 Quick Start

### Prerequisites
```bash
# Check installations
docker --version
kubectl version
k3d --version
```

### Run Current Applications

#### Log Output
```bash
# Build and deploy
docker build -t log-output:v1.0 log-output/
k3d image import log-output:v1.0 -c k3s-default
kubectl apply -f log-output/manifests/deployment.yaml

# View logs
kubectl logs -f deployment/log-output-dep
```

#### Todo Project
```bash
# Build and deploy
docker build -t todo-project:v1.5 todo-project/
k3d image import todo-project:v1.5 -c k3s-default
kubectl apply -f todo-project/manifests/deployment.yaml

# Port forward and test
kubectl port-forward deployment/todo-project-dep 3003:3000
# Open http://localhost:3003
```

---

## 💰 Cost Considerations

### Development
- **k3d (local K8s)**: FREE
- **Azure Free Tier**: $200 credit for new accounts

### Azure AI (Estimated)
- OpenAI (GPT-3.5): ~$10-20/month (development)
- Language Service: ~$5/month (within free tier)
- Computer Vision: ~$5/month (within free tier)
- **Total**: ~$20-30/month for development

---

## 📚 Resources

### Kubernetes
- [DevOps with Kubernetes Course](https://devopswithkubernetes.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [k3d Documentation](https://k3d.io/)

### Azure
- [Azure Kubernetes Service](https://learn.microsoft.com/en-us/azure/aks/)
- [Azure Verified Modules](https://azure.github.io/Azure-Verified-Modules/)
- [Terraform AzureRM Provider](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs)

### Azure AI
- [AI-102 Learning Path](https://learn.microsoft.com/en-us/credentials/certifications/azure-ai-engineer/)
- [Azure AI Services Docs](https://learn.microsoft.com/en-us/azure/ai-services/)
- [Azure OpenAI Documentation](https://learn.microsoft.com/en-us/azure/ai-services/openai/)

---

## 📜 Course Information

- **Course:** DevOps with Kubernetes
- **University:** University of Helsinki
- **Credits:** 5 ECTS
- **Submission:** Tags per exercise (e.g., `1.1`, `2.1`, `3.1`)

