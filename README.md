# AI-Enhanced Todo Application
## DevOps with Kubernetes + Azure AI Services (AI-102)

A production-grade, cloud-native todo application demonstrating Kubernetes orchestration and Azure AI capabilities. This project combines learning from:
- **DevOps with Kubernetes** course (University of Helsinki, 5 ECTS)
- **AI-102: Azure AI Engineer** certification path

---

## 🎯 Project Goals

1. Master Kubernetes concepts through hands-on development
2. Integrate Azure AI services (OpenAI, Language, Computer Vision)
3. Build a portfolio-worthy, production-ready application
4. Practice DevOps best practices (CI/CD, monitoring, security)

---

## 📋 Documentation

- **[PROJECT_PLAN.md](PROJECT_PLAN.md)** - Complete architecture and roadmap
- **[AZURE_SETUP.md](AZURE_SETUP.md)** - Azure AI services setup guide
- **[SUBMISSION_GUIDE.md](SUBMISSION_GUIDE.md)** - Course submission instructions

---

## 🚀 Current Status

### Completed Exercises

#### Part 1
- [x] **1.1** - Log Output: UUID generation with timestamps
- [x] **1.2** - Todo Project: Basic HTTP server
- [x] **1.3** - Log Output: Declarative Kubernetes deployment
- [x] **1.4** - Todo Project: Declarative Kubernetes deployment
- [x] **1.5** - Todo Project: HTTP server with HTML response
- [x] **1.6** - Todo Project: NodePort Service
- [x] **1.7** - Log Output: HTTP endpoint with Ingress
- [x] **1.8** - Todo Project: Ingress configuration
- [x] **1.9** - Ping-pong application with shared Ingress
- [x] **1.10** - Log Output: Split into two containers
- [x] **1.11** - Share data between applications with PersistentVolume
- [x] **1.12** - Todo Project: Cached daily image from Lorem Picsum
- [x] **1.13** - Todo Project: Todo list with input field, validation, and hardcoded todos

#### Part 2
- [x] **[2.1](https://github.com/priku/k8s-course/tree/2.1)** - Connecting pods: HTTP communication between log-output and ping-pong

### In Progress
- [ ] Azure AI services integration
- [ ] Microservices architecture

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
├── log-output/              # Exercises 1.1, 1.3
│   ├── main.go
│   ├── Dockerfile
│   └── manifests/
│       └── deployment.yaml
│
├── todo-project/            # Exercises 1.2, 1.4, 1.5
│   ├── main.go
│   ├── Dockerfile
│   └── manifests/
│       └── deployment.yaml
│
├── PROJECT_PLAN.md          # Complete project roadmap
├── AZURE_SETUP.md           # Azure AI setup guide
├── SUBMISSION_GUIDE.md      # Course submission guide
└── README.md                # This file
```

---

## 🛠️ Technology Stack

### Current
- **Language**: Go 1.21
- **Container**: Docker
- **Orchestration**: Kubernetes (k3d)
- **Tools**: kubectl, k3d

### Planned
- **Frontend**: React + TypeScript + Tailwind CSS
- **Backend**: Go (API Gateway), Python (AI Services)
- **Database**: PostgreSQL
- **AI**: Azure OpenAI, Language Service, Computer Vision
- **Monitoring**: Prometheus + Grafana
- **CI/CD**: GitHub Actions

---

## 🎓 Learning Path

### Phase 1: Kubernetes Foundations (Current)
- [x] Basic deployments and pods
- [ ] Services and networking
- [ ] ConfigMaps and Secrets
- [ ] Persistent storage

### Phase 2: AI Integration
- [ ] Set up Azure AI resources
- [ ] Build Python AI service
- [ ] Integrate with API gateway
- [ ] Deploy to Kubernetes

### Phase 3: Production Ready
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

### Azure AI
- [AI-102 Learning Path](https://learn.microsoft.com/en-us/credentials/certifications/azure-ai-engineer/)
- [Azure AI Services Docs](https://learn.microsoft.com/en-us/azure/ai-services/)
- [Azure OpenAI Documentation](https://learn.microsoft.com/en-us/azure/ai-services/openai/)

---

## 📝 Exercises Progress

### Chapter 1: Getting Started
- [1.1](../../tree/1.1/log-output) - Log Output: Application that generates and outputs a random string with timestamp every 5 seconds
- [1.2](../../tree/1.2/todo-project) - Todo Project: Web server with configurable PORT environment variable
- [1.3](../../tree/1.3/log-output) - Log Output: Declarative deployment with Kubernetes
- [1.4](../../tree/1.4/todo-project) - Todo Project: Declarative deployment with Kubernetes
- [1.5](../../tree/1.5/todo-project) - Todo Project: HTTP server responding to GET requests with HTML page
- [1.6](../../tree/1.6/todo-project) - Todo Project: NodePort Service for external access

## Project Structure

```
.
├── log-output/          # Exercise 1.1, 1.3 - Log output application
│   ├── main.go
│   ├── Dockerfile
│   ├── go.mod
│   └── manifests/
│       └── deployment.yaml
│
└── todo-project/        # Exercise 1.2, 1.4 - Todo project application
    ├── main.go
    ├── Dockerfile
    ├── go.mod
    └── manifests/
        └── deployment.yaml
```

## Applications

### Log Output (Exercises 1.1, 1.3)
A Go application that generates a random UUID on startup and outputs it with a timestamp every 5 seconds.

**Running locally:**
```bash
docker build -t log-output:v1.0 log-output/
docker run log-output:v1.0
```

**Deploying to Kubernetes:**
```bash
kubectl apply -f log-output/manifests/deployment.yaml
kubectl logs -f deployment/log-output-dep
```

### Todo Project (Exercises 1.2, 1.4)
A Go web server that can be configured via the PORT environment variable.

**Running locally:**
```bash
docker build -t todo-project:v1.0 todo-project/
docker run -e PORT=3000 -p 3000:3000 todo-project:v1.0
```

**Deploying to Kubernetes:**
```bash
kubectl apply -f todo-project/manifests/deployment.yaml
kubectl logs -f deployment/todo-project-dep
```

## Tools Used
- Kubernetes (k3d cluster)
- Docker
- Go 1.21
- kubectl

## Course Information
- **Course:** DevOps with Kubernetes
- **University:** University of Helsinki
- **Credits:** 5 ECTS
