# Technology Decisions and Implementation Plan

**Last Updated**: 2025-12-07
**Status**: Active Planning
**Strategy**: Follow Course Progression (K8s-first → Azure PaaS migration)

---

## Executive Summary

We've decided to **follow the DevOps with Kubernetes course exercises sequentially**, implementing everything in Kubernetes first to learn the fundamentals, then migrating to Azure PaaS services in later phases. This approach maximizes learning outcomes by experiencing both the operational complexity of running stateful services in K8s and the benefits of managed cloud services.

---

## Implementation Phases

### Phase 1: Kubernetes-First (Exercises 1.14-1.16, Part 2)
**Goal**: Learn K8s fundamentals by running everything in-cluster

**Timeline**: Current → Part 2 completion
**Focus**: Kubernetes patterns and operations

### Phase 2: Hybrid Approach (Part 3-4)
**Goal**: Migrate stateful services to Azure PaaS

**Timeline**: Part 3-4
**Focus**: Cloud integration and migration strategies

### Phase 3: Production Ready (Part 4-5)
**Goal**: Add advanced Azure features

**Timeline**: Part 4-5
**Focus**: Enterprise features and optimization

---

## Phase 1: Technology Stack (Current)

### Application Services (Running in Kubernetes)

| Component | Technology | Version | Rationale |
|-----------|-----------|---------|-----------|
| **Backend API** | Go + Chi | Go 1.21 | Lightweight router, stdlib-compatible, excellent middleware support |
| **Frontend** | React + TypeScript + Vite | React 18 | Modern, type-safe, fast builds, production-ready |
| **Styling** | Tailwind CSS | Latest | Utility-first, matches current design aesthetic |
| **Database** | PostgreSQL (StatefulSet) | 15 | Learn stateful workloads, ACID compliance |
| **Cache** | Redis (Deployment) | 7 | Learn in-cluster caching, simple deployment |
| **ORM** | GORM | Latest | Type-safe, excellent PostgreSQL support for Go |
| **API Design** | RESTful JSON | - | Standard, course-aligned, simple to learn |

### Kubernetes Resources

| Resource Type | Purpose | Learning Objective |
|--------------|---------|-------------------|
| **Deployments** | Frontend, Backend, Redis | Stateless application management |
| **StatefulSet** | PostgreSQL | Stateful application with stable network identity |
| **Services** | ClusterIP for inter-service communication | Service discovery |
| **Ingress** | External routing | HTTP routing, path-based routing |
| **ConfigMaps** | Application configuration | Configuration management |
| **Secrets** | Database credentials, API keys | Secret management |
| **PersistentVolumes** | Database storage | Persistent storage |
| **Init Containers** | Database migrations | Initialization patterns |

### Infrastructure

| Component | Technology | Environment |
|-----------|-----------|-------------|
| **Local K8s** | k3d | Development |
| **Container Registry** | Local (k3d import) | Development |
| **Container Runtime** | Docker | Development |

---

## Phase 2: Migration to Azure PaaS (Planned)

### Services to Migrate

| Service | From | To | Migration Trigger |
|---------|------|----|--------------------|
| **Database** | K8s StatefulSet | Azure Database for PostgreSQL Flexible Server | Part 3-4 exercises |
| **Cache** | K8s Redis | Azure Cache for Redis | Part 3-4 exercises |
| **Secrets** | K8s Secrets | Azure Key Vault | Part 2 (Secrets chapter) |
| **Storage** | PersistentVolumes | Azure Blob Storage | Part 3-4 exercises |
| **Container Registry** | Local | Azure Container Registry | Part 4 (Cloud chapter) |

### Benefits of Migration

**Azure Database for PostgreSQL:**
- Automatic backups and point-in-time restore
- Built-in high availability
- Automatic patching and updates
- Performance insights and monitoring
- Reduced operational burden

**Azure Cache for Redis:**
- Enterprise-grade availability
- Geo-replication
- Advanced data persistence
- Zone redundancy
- Managed updates

**Azure Key Vault:**
- Centralized secret management
- Automatic rotation
- Access policies and RBAC
- Audit logging
- HSM-backed keys

**Azure Blob Storage:**
- Infinite scalability
- CDN integration
- Lifecycle management
- Multiple access tiers
- Cost optimization

---

## Phase 3: Advanced Azure Integration (Planned)

### Additional Services

| Service | Purpose | Phase |
|---------|---------|-------|
| **Azure API Management** | API gateway with developer portal, rate limiting | Part 4-5 |
| **Azure OpenAI** | Direct API calls (no container needed) | Part 4 |
| **Azure Cognitive Services** | Language, Vision services | Part 4 |
| **Azure Monitor** | Centralized monitoring | Part 5 |
| **Application Insights** | APM and telemetry | Part 5 |
| **Managed Identity** | Passwordless authentication | Part 4 |

---

## Detailed Technology Decisions

### Decision 1: Follow Course Progression
- **Status**: ✅ Decided (2025-12-07)
- **Decision**: Implement in K8s first, migrate to Azure PaaS later
- **Alternatives Considered**:
  - Jump to Azure PaaS immediately
  - Stay K8s-only throughout
- **Rationale**:
  - Learn fundamentals before optimization
  - Understand why PaaS is better through experience
  - Course likely expects K8s implementation first
  - Better learning outcomes from doing it both ways
- **Trade-offs**:
  - ✅ Pros: Deeper K8s knowledge, understand operations burden
  - ❌ Cons: More initial complexity, temporary technical debt
- **Cost Impact**: Minimal (k3d is free for Phase 1)

### Decision 2: Backend Framework - Chi
- **Status**: ✅ Decided (2025-12-07)
- **Decision**: Use Chi router for Go backend
- **Alternatives Considered**:
  - **Gin**: More batteries-included, less idiomatic
  - **Echo**: Similar to Gin, heavier
  - **stdlib only**: Too minimal for real API
- **Rationale**:
  - Lightweight and stdlib-compatible
  - Excellent middleware ecosystem
  - Idiomatic Go patterns
  - Production-proven
  - Easy to learn
- **Trade-offs**:
  - ✅ Pros: Clean code, flexibility, performance
  - ❌ Cons: Less built-in features than Gin
- **References**: [Chi GitHub](https://github.com/go-chi/chi)

### Decision 3: Frontend - React with Separate Container
- **Status**: ✅ Decided (2025-12-07)
- **Decision**: Build separate React frontend, serve with nginx
- **Alternatives Considered**:
  - Keep embedded HTML in Go (current)
  - Vue.js
  - Svelte
- **Rationale**:
  - Separation of concerns (frontend/backend)
  - Industry standard technology
  - Type safety with TypeScript
  - Reusable for mobile/other clients
  - Matches PROJECT_PLAN.md goals
- **Trade-offs**:
  - ✅ Pros: Production-ready architecture, scalable, modern
  - ❌ Cons: More complexity, extra container
- **Implementation**:
  - Vite for fast builds
  - Tailwind CSS for styling (consistent with current design)
  - Axios for API calls
  - Build to static files
  - Nginx container to serve

### Decision 4: Database - PostgreSQL in K8s (Phase 1)
- **Status**: ✅ Decided (2025-12-07, temporary)
- **Decision**: PostgreSQL StatefulSet → Migrate to Azure DB later
- **Alternatives Considered**:
  - Azure Database immediately
  - MySQL
  - MongoDB
  - SQLite (too simple)
- **Rationale**:
  - **Phase 1**: Learn K8s StatefulSets and stateful workload management
  - **Phase 2**: Experience migration to managed service
  - PostgreSQL is production-grade and ACID-compliant
  - Excellent GORM support
- **Trade-offs**:
  - ✅ Pros: Learn K8s deeply, understand operations
  - ❌ Cons: Operational burden (backups, HA, scaling)
- **Migration Plan**: Part 3-4 exercises
- **K8s Patterns to Learn**:
  - StatefulSet with stable network identity
  - PersistentVolumeClaim
  - Init containers for schema migrations
  - Headless service
  - Backup strategies

### Decision 5: AI Service Architecture
- **Status**: ✅ Decided (2025-12-07)
- **Decision**: Direct Azure API calls from Go (no Python container in Phase 1)
- **Alternatives Considered**:
  - Python FastAPI container
  - Go service calling Python service
  - Separate microservice
- **Rationale**:
  - Simpler architecture for Phase 1
  - Azure SDK available for Go
  - No need for extra container initially
  - Can add Python service later if needed
- **Trade-offs**:
  - ✅ Pros: Simpler, faster, fewer containers
  - ❌ Cons: Less microservices practice
- **Future**: Add Python AI service in Phase 4 if complex AI workflows needed

### Decision 6: API Design - RESTful
- **Status**: ✅ Decided (2025-12-07)
- **Decision**: RESTful JSON API
- **Alternatives Considered**:
  - GraphQL (flexible but complex)
  - gRPC (fast but overkill for web)
- **Rationale**:
  - Industry standard
  - Course likely expects REST
  - Simple to learn and implement
  - Good tooling support
- **Trade-offs**:
  - ✅ Pros: Standard, simple, well-understood
  - ❌ Cons: Less flexible than GraphQL
- **API Versioning**: `/api/v1/` prefix for future compatibility

### Decision 7: Cache - Redis in K8s (Phase 1)
- **Status**: ✅ Decided (2025-12-07)
- **Decision**: Redis Deployment in K8s → Migrate to Azure Cache later
- **Alternatives Considered**:
  - In-memory (too simple)
  - Azure Cache immediately
- **Rationale**:
  - Learn K8s deployment patterns
  - Simple to set up (single pod initially)
  - Experience before migrating to managed service
- **Migration Plan**: Part 3-4 exercises

---

## Architecture Diagrams

### Phase 1: Kubernetes-First (Current → Part 2)

```
                    Internet
                       │
                       ▼
              ┌────────────────┐
              │  K8s Ingress   │
              │  /     → Frontend
              │  /api  → Backend
              └───────┬────────┘
                      │
         ┌────────────┴────────────┐
         │                         │
         ▼                         ▼
  ┌─────────────┐          ┌─────────────┐
  │  Frontend   │          │   Backend   │
  │   (nginx)   │          │   (Go+Chi)  │
  │ React build │          │  REST API   │
  │             │          │             │
  │  Port: 80   │          │  Port: 8080 │
  └─────────────┘          └──────┬──────┘
                                  │
                    ┌─────────────┼──────────────┐
                    ▼             ▼              ▼
            ┌───────────┐  ┌──────────┐  ┌──────────┐
            │PostgreSQL │  │  Redis   │  │ConfigMap │
            │StatefulSet│  │Deployment│  │ Secrets  │
            │   + PVC   │  │          │  │          │
            │Port: 5432 │  │Port: 6379│  └──────────┘
            └───────────┘  └──────────┘
```

### Phase 2: Hybrid Architecture (Part 3-4)

```
                    Internet
                       │
                       ▼
              ┌────────────────┐
              │  K8s Ingress   │
              └───────┬────────┘
                      │
         ┌────────────┴────────────┐
         │                         │
         ▼                         ▼
  ┌─────────────┐          ┌─────────────┐
  │  Frontend   │          │   Backend   │
  │    (K8s)    │          │    (K8s)    │
  └─────────────┘          └──────┬──────┘
                                  │
                    ┌─────────────┼──────────────┐
                    ▼             ▼              ▼
            ┌───────────┐  ┌──────────┐  ┌──────────┐
            │  Azure    │  │  Azure   │  │  Azure   │
            │ Database  │  │  Cache   │  │   Key    │
            │PostgreSQL │  │  Redis   │  │  Vault   │
            └───────────┘  └──────────┘  └──────────┘
                 (PaaS)        (PaaS)        (PaaS)
```

### Phase 3: Full Azure Integration (Part 4-5)

```
                    Internet
                       │
                       ▼
         ┌──────────────────────────┐
         │ Azure API Management     │
         │ - Developer Portal       │
         │ - Rate Limiting          │
         │ - Analytics              │
         └────────────┬─────────────┘
                      │
                      ▼
              ┌────────────────┐
              │  K8s Ingress   │
              └───────┬────────┘
                      │
         ┌────────────┴────────────┐
         │                         │
         ▼                         ▼
  ┌─────────────┐          ┌─────────────┐
  │  Frontend   │          │   Backend   │
  │   (AKS)     │          │   (AKS)     │
  └─────────────┘          └──────┬──────┘
                                  │
         ┌────────────────────────┼────────────────────┐
         │                        │                    │
         ▼                        ▼                    ▼
  ┌────────────┐         ┌────────────┐       ┌──────────────┐
  │   Azure    │         │   Azure    │       │  Azure AI    │
  │  Database  │         │   Cache    │       │   Services   │
  │            │         │            │       │ - OpenAI     │
  └────────────┘         └────────────┘       │ - Language   │
         ▲                                    │ - Vision     │
         │                                    └──────────────┘
         │
  ┌────────────┐         ┌────────────┐       ┌──────────────┐
  │   Azure    │         │   Azure    │       │    Azure     │
  │    Blob    │         │    Key     │       │   Monitor    │
  │  Storage   │         │   Vault    │       │  + Insights  │
  └────────────┘         └────────────┘       └──────────────┘
```

---

## Immediate Next Steps (Exercise 1.14+)

### Exercise 1.14: Project v0.4 - Backend for Todos

**Expected Requirements** (based on typical course progression):
- Create a backend service that stores todos
- RESTful API endpoints
- Connect frontend to backend

**Implementation Plan:**

#### 1. Create Backend Service (`todo-backend`)
**Location**: `c:\k8s-course\services\todo-backend\`

**Structure**:
```
services/todo-backend/
├── main.go              # Entry point
├── handlers/            # HTTP handlers
│   └── todos.go
├── models/              # Data models
│   └── todo.go
├── storage/             # Storage layer (in-memory for 1.14)
│   └── memory.go
├── Dockerfile           # Multi-stage build
├── go.mod
└── go.sum
```

**API Endpoints**:
- `GET /api/todos` - List all todos
- `POST /api/todos` - Create todo
- `GET /api/todos/:id` - Get single todo
- `PUT /api/todos/:id` - Update todo
- `DELETE /api/todos/:id` - Delete todo

**Technology**:
- Chi router
- In-memory storage (slice) for Exercise 1.14
- JSON responses
- CORS middleware

#### 2. Create Frontend Service (`todo-frontend`)
**Location**: `c:\k8s-course\services\todo-frontend\`

**Structure**:
```
services/todo-frontend/
├── src/
│   ├── App.tsx
│   ├── components/
│   │   ├── TodoList.tsx
│   │   ├── TodoItem.tsx
│   │   └── TodoForm.tsx
│   ├── api/
│   │   └── todos.ts
│   └── main.tsx
├── public/
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
├── tailwind.config.js
├── Dockerfile           # Multi-stage: build + nginx
└── nginx.conf
```

**Features**:
- List todos from API
- Add new todo
- Mark as complete
- Delete todo
- Loading states
- Error handling

#### 3. Kubernetes Manifests
**Location**: `c:\k8s-course\manifests\todo-app\`

**Structure**:
```
manifests/todo-app/
├── backend/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── configmap.yaml
├── frontend/
│   ├── deployment.yaml
│   └── service.yaml
└── ingress.yaml
```

**Ingress Routing**:
```yaml
- path: /api
  pathType: Prefix
  backend:
    service:
      name: todo-backend-svc
      port: 8080
- path: /
  pathType: Prefix
  backend:
    service:
      name: todo-frontend-svc
      port: 80
```

---

### Exercise 1.15: Database Integration

**Expected Requirements**:
- Add PostgreSQL for persistent storage
- Database migrations
- Connect backend to database

**Implementation Plan**:

#### 1. PostgreSQL StatefulSet
**Location**: `c:\k8s-course\manifests\database\`

**Components**:
- StatefulSet for PostgreSQL
- PersistentVolumeClaim
- Service (headless)
- ConfigMap (init SQL)
- Secret (credentials)

#### 2. Update Backend
**Changes**:
- Add GORM
- Create database models
- Database migrations
- Connection pooling
- Environment-based config

#### 3. Kubernetes Patterns
**Learning**:
- StatefulSet configuration
- Init containers for migrations
- Secret management
- Service discovery

---

### Exercise 1.16: Project v0.5 - Complete CRUD

**Expected Requirements**:
- Full CRUD operations
- Update todos (edit text, toggle complete)
- Better error handling
- UI polish

**Implementation Plan**:
- Add PUT endpoint for updates
- Edit functionality in frontend
- Mark complete/incomplete
- Validation
- Better UI/UX

---

## Data Models

### Todo Model (Phase 1)

```go
// models/todo.go
type Todo struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Text        string    `json:"text" gorm:"not null"`
    Completed   bool      `json:"completed" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Database Schema (Exercise 1.15+)

```sql
CREATE TABLE todos (
    id VARCHAR(36) PRIMARY KEY,
    text VARCHAR(255) NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_todos_completed ON todos(completed);
CREATE INDEX idx_todos_created_at ON todos(created_at);
```

---

## Cost Estimates

### Phase 1: K8s-First (Local Development)
- k3d cluster: **FREE**
- Docker Desktop: **FREE**
- Total: **$0/month**

### Phase 2: Hybrid (Azure PaaS - Development)
- AKS (2 nodes, B2s): ~$70/month
- Azure Database PostgreSQL (Burstable B1ms): ~$12/month
- Azure Cache for Redis (Basic C0): ~$16/month
- Azure Blob Storage: <$1/month
- Azure Key Vault: ~$0.03/month
- **Total: ~$100/month**

### Phase 3: Full Azure (Development)
- Phase 2 costs: ~$100/month
- Azure API Management (Developer): ~$50/month
- Azure OpenAI (moderate usage): ~$20/month
- Azure Monitor + App Insights: ~$10/month
- **Total: ~$180/month**

### Cost Optimization Strategies
- Use free tier credits ($200 for new accounts)
- Auto-shutdown dev resources nights/weekends
- Use cheaper tiers for development
- Implement caching aggressively
- Set budget alerts
- Use spot instances for AKS nodes

---

## Testing Strategy

### Phase 1
- **Unit Tests**: stdlib `testing` package
- **Integration Tests**: Test against real containers (docker-compose)
- **E2E Tests**: Manual testing in k3d cluster

### Phase 2+
- Add API integration tests
- Database migration tests
- Performance testing
- Load testing with k6

---

## Migration Checklist (Phase 1 → Phase 2)

### Database Migration
- [ ] Export data from K8s PostgreSQL
- [ ] Provision Azure Database for PostgreSQL
- [ ] Configure VNet integration
- [ ] Update connection strings in ConfigMap
- [ ] Test connection from K8s pods
- [ ] Migrate data
- [ ] Update backup strategy
- [ ] Decommission K8s PostgreSQL

### Cache Migration
- [ ] Provision Azure Cache for Redis
- [ ] Update connection strings
- [ ] Test connection
- [ ] Warm cache if needed
- [ ] Switch traffic
- [ ] Decommission K8s Redis

### Secrets Migration
- [ ] Create Azure Key Vault
- [ ] Configure Managed Identity
- [ ] Migrate secrets
- [ ] Update pods to use Key Vault
- [ ] Test secret rotation
- [ ] Remove K8s secrets

---

## Review Schedule

- **After Exercise 1.14**: Review backend/frontend architecture
- **After Exercise 1.16**: Review Phase 1 completion, plan Phase 2
- **After Part 2**: Review K8s patterns learned, plan Azure migration
- **After Part 4**: Review Azure integration, plan production deployment

---

## References

### Courses
- [DevOps with Kubernetes](https://devopswithkubernetes.com/)
- [AI-102 Learning Path](https://learn.microsoft.com/en-us/credentials/certifications/azure-ai-engineer/)

### Documentation
- [Kubernetes Official Docs](https://kubernetes.io/docs/)
- [Go Chi Router](https://github.com/go-chi/chi)
- [GORM](https://gorm.io/)
- [React + TypeScript](https://react.dev/)
- [Vite](https://vitejs.dev/)

### Azure
- [Azure Database for PostgreSQL](https://learn.microsoft.com/en-us/azure/postgresql/)
- [Azure Cache for Redis](https://learn.microsoft.com/en-us/azure/azure-cache-for-redis/)
- [Azure Key Vault](https://learn.microsoft.com/en-us/azure/key-vault/)
- [Azure API Management](https://learn.microsoft.com/en-us/azure/api-management/)

---

## Change Log

| Date | Change | Reason |
|------|--------|--------|
| 2025-12-07 | Initial document created | Formalize technology decisions |
| 2025-12-07 | Decided on K8s-first approach | Follow course progression |
| 2025-12-07 | Selected Chi, React, PostgreSQL | Best fit for learning + production |

---

**Next Document Update**: After completing Exercise 1.14
