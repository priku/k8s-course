# AI-Enhanced Todo Application
## Kubernetes + Azure AI Services Integration Project

### Project Vision
Build a production-grade, AI-powered Todo application that demonstrates both Kubernetes orchestration skills and Azure AI capabilities from AI-102 certification.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Frontend   │  │   API        │  │  AI Services │      │
│  │   (React)    │→ │  Gateway     │→ │  Container   │      │
│  │   Pod(s)     │  │   (Go)       │  │   Pod(s)     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│         ↓                  ↓                  ↓             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Ingress    │  │   Service    │  │   Service    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │         Azure AI Services (External)               │     │
│  │  - Azure OpenAI (GPT-4)                           │     │
│  │  - Language Service (Sentiment, Key Phrases)      │     │
│  │  - Computer Vision (Image Analysis)               │     │
│  │  - Speech Services (Optional)                     │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

---

## Learning Path Integration

### Phase 1: Kubernetes Foundations (Weeks 1-2)
**Current K8s Course: Exercises 1.1-1.10**
- Basic deployments and pods
- Networking and services
- ConfigMaps and Secrets
- Persistent storage
- **Deliverable**: Simple todo-app running in K8s

### Phase 2: AI Services Integration (Weeks 3-4)
**AI-102 Skills: Implement and Deploy**
- Azure OpenAI integration
- Language services (sentiment analysis)
- Computer Vision (image attachments)
- Secrets management for API keys
- **Deliverable**: AI-enhanced microservices

### Phase 3: Advanced K8s (Weeks 5-6)
**K8s Course: Chapter 2-3**
- Service mesh
- Monitoring and logging
- Scaling strategies
- Health checks
- **Deliverable**: Production-ready cluster

### Phase 4: CI/CD & Production (Weeks 7-8)
**Combined Skills**
- GitHub Actions for builds
- Automated deployments
- A/B testing AI models
- Cost optimization
- **Deliverable**: Full DevOps pipeline

---

## Feature Roadmap

### MVP Features (Phase 1)
- [x] Basic HTTP server (Exercise 1.5)
- [ ] CRUD operations for todos
- [ ] Persistent storage (PostgreSQL in K8s)
- [ ] RESTful API
- [ ] Simple web UI

### AI Features (Phase 2)

#### 1. Smart Todo Creation
**AI-102 Skills: Azure OpenAI**
```
User types: "Buy groceries tomorrow at 3pm"
↓
AI extracts:
- Task: "Buy groceries"
- Due date: Tomorrow 15:00
- Priority: Medium
- Category: Shopping
```

#### 2. Sentiment-Based Priority
**AI-102 Skills: Language Service**
```
User types: "URGENT! Fix production bug!!!"
↓
Sentiment Analysis:
- Sentiment: Negative (0.95)
- Priority: Automatically set to HIGH
- Notifications triggered
```

#### 3. Image Attachments
**AI-102 Skills: Computer Vision**
```
User uploads image of a receipt
↓
Vision API:
- Extracts text (OCR)
- Identifies objects
- Generates tags
- Creates todo: "Process receipt - $45.99 at Store X"
```

#### 4. AI Assistant
**AI-102 Skills: Azure OpenAI + Prompt Engineering**
```
User asks: "What should I focus on today?"
↓
AI analyzes todos and responds:
"Based on your 12 pending tasks, I recommend:
1. Complete the production bug fix (high priority, overdue)
2. Finish the presentation (due in 2 hours)
3. Review PRs (3 waiting for your input)"
```

### Advanced Features (Phase 3-4)
- [ ] Voice input (Speech Services)
- [ ] Multi-language support (Translator)
- [ ] Anomaly detection (unusual task patterns)
- [ ] Collaboration features
- [ ] Mobile app
- [ ] Real-time notifications

---

## Technology Stack

### Frontend
- **Framework**: React with TypeScript
- **Styling**: Tailwind CSS
- **State**: React Query + Context
- **Container**: Nginx serving static files

### Backend Services

#### API Gateway (Go)
- **Purpose**: Main entry point, routing
- **Framework**: Gin or Chi
- **Features**:
  - Authentication/Authorization
  - Rate limiting
  - Request logging
  - API versioning

#### AI Service (Python)
- **Purpose**: Azure AI SDK integration
- **Framework**: FastAPI
- **SDKs**:
  - `azure-ai-openai`
  - `azure-ai-textanalytics`
  - `azure-ai-vision`
  - `azure-cognitiveservices-speech`

#### Data Service (Go)
- **Purpose**: Database operations
- **Framework**: GORM
- **Database**: PostgreSQL

### Infrastructure
- **Orchestration**: Kubernetes (k3d locally, AKS for production)
- **Storage**: Persistent Volumes
- **Secrets**: Kubernetes Secrets + Azure Key Vault
- **Monitoring**: Prometheus + Grafana
- **Logging**: Loki + FluentBit

---

## Directory Structure

```
k8s-course/
├── frontend/
│   ├── src/
│   ├── Dockerfile
│   └── package.json
│
├── services/
│   ├── api-gateway/          # Go
│   │   ├── main.go
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── ai-service/           # Python + Azure AI
│   │   ├── main.py
│   │   ├── requirements.txt
│   │   ├── openai_handler.py
│   │   ├── vision_handler.py
│   │   ├── language_handler.py
│   │   └── Dockerfile
│   │
│   └── data-service/         # Go + PostgreSQL
│       ├── main.go
│       ├── models/
│       ├── Dockerfile
│       └── go.mod
│
├── manifests/
│   ├── namespace.yaml
│   ├── frontend/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   ├── api-gateway/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── configmap.yaml
│   ├── ai-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── secret.yaml        # Azure AI keys
│   ├── data-service/
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── database/
│       ├── statefulset.yaml
│       ├── service.yaml
│       └── pvc.yaml
│
├── scripts/
│   ├── build-all.sh
│   ├── deploy-local.sh
│   └── deploy-azure.sh
│
├── .github/
│   └── workflows/
│       ├── ci.yaml
│       └── cd.yaml
│
├── docs/
│   ├── API.md
│   ├── ARCHITECTURE.md
│   └── AI_FEATURES.md
│
├── log-output/              # Original exercise
├── todo-project/            # Original exercise
├── PROJECT_PLAN.md          # This file
└── README.md
```

---

## AI-102 Skills Coverage

### ✅ Plan and Manage an Azure AI Solution
- [ ] Create and manage Azure AI resources
- [ ] Secure Azure AI services
- [ ] Monitor Azure AI services
- [ ] Deploy AI models to containers

### ✅ Implement Computer Vision Solutions
- [ ] Analyze images using Vision API
- [ ] Extract text from images (OCR)
- [ ] Generate image descriptions
- [ ] Detect objects and faces

### ✅ Implement Natural Language Processing
- [ ] Analyze text with Language Service
- [ ] Extract key phrases
- [ ] Perform sentiment analysis
- [ ] Entity recognition
- [ ] Custom text classification

### ✅ Implement Knowledge Mining
- [ ] Create enrichment pipelines
- [ ] Index searchable content

### ✅ Implement Generative AI Solutions
- [ ] Use Azure OpenAI for completions
- [ ] Implement chat interfaces
- [ ] Prompt engineering techniques
- [ ] Function calling with GPT-4

---

## Kubernetes Skills Coverage

### ✅ Chapter 1: Introduction
- [x] Basic deployments (1.1-1.4)
- [x] HTTP services (1.5)
- [ ] Persistence (1.6-1.10)

### ✅ Chapter 2: Basics
- [ ] Services and networking
- [ ] ConfigMaps and environment
- [ ] Secrets management
- [ ] Health checks

### ✅ Chapter 3: Advanced Concepts
- [ ] Ingress controllers
- [ ] StatefulSets
- [ ] Jobs and CronJobs
- [ ] Resource limits

### ✅ Chapter 4: Update Strategies
- [ ] Rolling updates
- [ ] Blue-green deployments
- [ ] Canary releases

### ✅ Chapter 5: Monitoring
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Logging aggregation
- [ ] Distributed tracing

---

## Development Workflow

### Daily/Weekly Pattern
1. **Morning**: Study K8s course chapter
2. **Midday**: Implement feature in code
3. **Afternoon**: Study AI-102 concept
4. **Evening**: Integrate AI feature into K8s app

### Example Week
**Monday**
- K8s: Learn about Services
- Code: Create API gateway service
- AI-102: Study Azure OpenAI basics
- Integration: Plan AI service architecture

**Tuesday**
- K8s: Learn about ConfigMaps/Secrets
- Code: Implement configuration
- AI-102: Set up Azure OpenAI resource
- Integration: Store API keys in K8s secrets

**Wednesday**
- K8s: Learn about Ingress
- Code: Set up routing
- AI-102: Implement GPT-4 completions
- Integration: Deploy AI service to cluster

**Thursday**
- K8s: Learn about Persistence
- Code: Add PostgreSQL
- AI-102: Study Language Service
- Integration: Store AI responses in DB

**Friday**
- Review week's progress
- Integration testing
- Documentation
- Plan next week

---

## Success Metrics

### Technical Metrics
- [ ] All K8s course exercises completed
- [ ] AI-102 certification passed
- [ ] 90%+ test coverage
- [ ] < 200ms API response time
- [ ] Zero-downtime deployments
- [ ] Automated CI/CD pipeline

### Learning Metrics
- [ ] Understand K8s networking
- [ ] Master Azure AI SDKs
- [ ] Can troubleshoot production issues
- [ ] Can explain architecture decisions
- [ ] Portfolio-worthy project

---

## Cost Management

### Development (Local)
- k3d cluster: **FREE**
- Local Docker: **FREE**
- Azure Free Tier: **$200 credit**

### Azure AI Costs (Estimated Monthly)
- Azure OpenAI (GPT-4): ~$20-50
- Language Service: ~$5-10
- Computer Vision: ~$5-10
- **Total**: ~$30-70/month

### Cost Optimization
- Use free tier limits
- Implement caching
- Batch AI requests
- Auto-shutdown dev resources
- Monitor usage with alerts

---

## Next Steps (Immediate)

### This Week
1. ✅ Complete Exercise 1.5 (DONE)
2. [ ] Set up Azure account and AI resources
3. [ ] Continue K8s exercises 1.6-1.10
4. [ ] Design database schema
5. [ ] Create API specifications

### Next Week
1. [ ] Build API gateway (Go)
2. [ ] Create React frontend skeleton
3. [ ] Integrate first AI feature (OpenAI)
4. [ ] Deploy multi-service app to K8s

---

## Resources

### Kubernetes
- Course: https://courses.mooc.fi/org/uh-cs/courses/devops-with-kubernetes
- Docs: https://kubernetes.io/docs/
- k3d: https://k3d.io/

### AI-102
- Learning Path: https://learn.microsoft.com/en-us/credentials/certifications/azure-ai-engineer/
- Azure AI Docs: https://learn.microsoft.com/en-us/azure/ai-services/
- OpenAI SDK: https://learn.microsoft.com/en-us/azure/ai-services/openai/

### Tools
- Docker: https://docs.docker.com/
- Go: https://go.dev/doc/
- React: https://react.dev/
- Python: https://docs.python.org/

---

## Questions to Consider

1. **Should we use AKS (Azure Kubernetes Service) or stick with local k3d?**
   - Start local, migrate to AKS in Phase 3

2. **How to handle AI API costs during development?**
   - Implement aggressive caching
   - Use cheaper models for testing
   - Mock responses in tests

3. **Monorepo vs separate repos?**
   - Monorepo for easier learning/development
   - Can split later if needed

4. **Which AI feature to implement first?**
   - Start with sentiment analysis (simplest)
   - Then OpenAI (most impressive)
   - Vision last (most complex)

---

**Last Updated**: 2025-11-29
**Status**: Planning Phase
**Next Review**: After Exercise 1.10
