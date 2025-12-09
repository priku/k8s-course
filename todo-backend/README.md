# Todo Backend Application

A RESTful API backend service built with Go and PostgreSQL that manages todo items. This service provides CRUD operations for todos with persistent storage.

## Exercises

This application fulfills the requirements for:
- **Exercise 2.2:** REST API service with GET/POST endpoints for todos
- **Exercise 2.4:** Move to "project" namespace
- **Exercise 2.8:** PostgreSQL StatefulSet for todo persistence with Kubernetes Secrets

## Features

- RESTful API for todo management
- PostgreSQL database integration for persistent storage
- CORS enabled for frontend integration
- UUID-based todo IDs
- 140-character limit validation
- Health check endpoint
- Automatic database schema initialization
- Secure credential management with Kubernetes Secrets

## API Endpoints

### GET /todos
Retrieves all todos, ordered by creation date (newest first).

**Response:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "text": "Buy groceries",
    "done": false,
    "created_at": "2025-12-10T10:30:00Z"
  }
]
```

### POST /todos
Creates a new todo item.

**Request:**
```json
{
  "text": "Buy groceries"
}
```

**Response:** (201 Created)
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "text": "Buy groceries",
  "done": false,
  "created_at": "2025-12-10T10:30:00Z"
}
```

**Validations:**
- Text is required (400 Bad Request)
- Maximum 140 characters (400 Bad Request)

### GET /healthz
Health check endpoint for liveness/readiness probes.

**Response:** (200 OK)
```
OK
```

## Running Locally

### Using Docker

Build the image:
```bash
docker build -t todo-backend:v2.10 .
```

Run with PostgreSQL:
```bash
# Start PostgreSQL
docker run -d --name postgres-todos \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=example \
  -e POSTGRES_DB=todos \
  -p 5432:5432 \
  postgres:13.0

# Run todo-backend
docker run -p 3000:3000 \
  -e POSTGRES_URL="postgres://postgres:example@host.docker.internal:5432/todos?sslmode=disable" \
  todo-backend:v2.10
```

Test the API:
```bash
# Get all todos
curl http://localhost:3000/todos

# Create a todo
curl -X POST http://localhost:3000/todos \
  -H "Content-Type: application/json" \
  -d '{"text":"Buy groceries"}'

# Health check
curl http://localhost:3000/healthz
```

## Deploying to Kubernetes

### Local Deployment (k3d)

#### Prerequisites
- k3d cluster running
- kubectl configured
- PostgreSQL StatefulSet deployed

#### Deploy PostgreSQL StatefulSet
```bash
kubectl apply -f manifests/secret.yaml
kubectl apply -f manifests/statefulset.yaml
```

Wait for PostgreSQL to be ready:
```bash
kubectl wait --for=condition=ready pod -l app=postgres-todos -n project --timeout=300s
```

#### Import image to k3d
```bash
k3d image import todo-backend:v2.10 -c k3s-default
```

#### Deploy Todo Backend Application
```bash
kubectl apply -f manifests/deployment.yaml
kubectl apply -f manifests/service.yaml
kubectl apply -f manifests/ingress.yaml
kubectl apply -f manifests/middleware.yaml
```

#### Check deployment
```bash
kubectl get deployments -n project
kubectl get pods -n project
kubectl get services -n project
```

#### View logs
```bash
kubectl logs -f deployment/todo-backend-dep -n project
```

#### Test the application
If using Ingress with path prefix stripping:
```bash
# Get all todos
curl http://localhost:8081/api/todos

# Create a todo
curl -X POST http://localhost:8081/api/todos \
  -H "Content-Type: application/json" \
  -d '{"text":"Deploy to Kubernetes"}'
```

Or using port-forward:
```bash
kubectl port-forward -n project deployment/todo-backend-dep 3000:3000
curl http://localhost:3000/todos
```

#### Delete deployment
```bash
kubectl delete -f manifests/deployment.yaml
kubectl delete -f manifests/service.yaml
kubectl delete -f manifests/ingress.yaml
kubectl delete -f manifests/middleware.yaml
kubectl delete -f manifests/statefulset.yaml
kubectl delete -f manifests/secret.yaml
```

### Cloud Deployment (Azure AKS)

#### Prerequisites
- Azure CLI installed and authenticated (`az login`)
- AKS cluster created (see [terraform/README.md](../terraform/README.md))
- kubectl configured for AKS (`az aks get-credentials`)
- Azure Container Registry (ACR) access

#### Build and Push to ACR
```bash
# Get ACR name from Terraform output
ACR_NAME=$(cd ../terraform/bootstrap && terraform output -raw acr_name)

# Login to ACR
az acr login --name $ACR_NAME

# Build and tag image
docker build -t todo-backend:v2.10 .
docker tag todo-backend:v2.10 $ACR_NAME.azurecr.io/todo-backend:v2.10

# Push to ACR
docker push $ACR_NAME.azurecr.io/todo-backend:v2.10
```

#### Update manifests for AKS
Update the image reference in `manifests/deployment.yaml`:
```yaml
spec:
  containers:
    - name: todo-backend
      image: <acr-name>.azurecr.io/todo-backend:v2.10
```

#### Deploy to AKS
```bash
# Deploy secrets and PostgreSQL
kubectl apply -f manifests/secret.yaml
kubectl apply -f manifests/statefulset.yaml

# Wait for PostgreSQL
kubectl wait --for=condition=ready pod -l app=postgres-todos -n project --timeout=300s

# Deploy application
kubectl apply -f manifests/deployment.yaml
kubectl apply -f manifests/service.yaml
```

**Note:** For AKS, you may want to use LoadBalancer or configure Azure Ingress Controller for production deployments.

#### Check AKS deployment
```bash
kubectl get deployments -n project
kubectl get pods -n project
kubectl get services -n project
```

#### View logs in AKS
```bash
kubectl logs -f deployment/todo-backend-dep -n project
```

#### Delete AKS deployment
```bash
kubectl delete -f manifests/deployment.yaml
kubectl delete -f manifests/service.yaml
kubectl delete -f manifests/statefulset.yaml
kubectl delete -f manifests/secret.yaml
```

## Configuration

The application uses the following environment variables:

- `PORT` - The port number the server listens on (default: 3000)
- `POSTGRES_URL` - PostgreSQL connection string (stored in Kubernetes Secret)

## Database Schema

The application automatically creates the following table:

```sql
CREATE TABLE IF NOT EXISTS todos (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Files

- `main.go` - Main application code with REST API and database integration
- `go.mod` - Go module definition
- `Dockerfile` - Multi-stage Docker build
- `manifests/deployment.yaml` - Kubernetes deployment in "project" namespace
- `manifests/service.yaml` - ClusterIP service
- `manifests/ingress.yaml` - Ingress configuration for /api path
- `manifests/middleware.yaml` - Traefik middleware to strip /api prefix
- `manifests/secret.yaml` - Kubernetes Secret for database credentials
- `manifests/statefulset.yaml` - PostgreSQL StatefulSet with persistent storage

## Technologies

- **Language:** Go 1.21
- **Database:** PostgreSQL 13.0
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes
- **Storage:** StatefulSet with PersistentVolumeClaim
- **Security:** Kubernetes Secrets for sensitive data

## Integration with Todo-Project

The todo-backend service provides the REST API consumed by the todo-project frontend. The frontend sends requests to `/api/todos`, which are routed through the Ingress with path prefix stripping middleware to this backend service.

## Security Best Practices

- Database credentials stored in Kubernetes Secrets (not hard-coded)
- Secrets referenced via `secretKeyRef` in deployment
- CORS configured for controlled access
- Input validation (length limits, required fields)
- SQL injection protection via parameterized queries

## Request Logging

The application logs all incoming requests with the following information:
- HTTP method and path
- Remote address
- Request validation results
- Todo creation success/failure
- Database operation results

Example log output:
```
Todo-backend server started on port 3000
Database connected: postgres://postgres:***@postgres-svc:5432/todos?sslmode=disable
POST /todos - Request from 10.42.0.1:52134
POST /todos - Received todo text (length: 15): Buy groceries
POST /todos - SUCCESS: Created todo 123e4567-e89b-12d3-a456-426614174000 - Buy groceries
```
