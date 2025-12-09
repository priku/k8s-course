# Ping-Pong Application

A simple Go web service that maintains a counter with persistent storage using PostgreSQL. Each request to `/pingpong` increments the counter and returns "pong N".

## Exercises

This application fulfills the requirements for:
- **Exercise 1.9:** Ping-pong application with shared Ingress
- **Exercise 2.1:** HTTP communication between pods (provides `/count` endpoint)
- **Exercise 2.7:** StatefulSet with PostgreSQL for counter persistence

## Features

- HTTP endpoint `/pingpong` - Increments counter and returns "pong N"
- HTTP endpoint `/count` - Returns current counter as JSON
- PostgreSQL database integration for persistent counter storage
- Configurable via environment variables
- Automatic database schema initialization

## Running Locally

### Using Docker

Build the image:
```bash
docker build -t ping-pong:v2.7 .
```

Run with PostgreSQL:
```bash
# Start PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=example \
  -p 5432:5432 \
  postgres:13.0

# Run ping-pong
docker run -p 3000:3000 \
  -e POSTGRES_URL="postgres://postgres:example@host.docker.internal:5432/postgres?sslmode=disable" \
  ping-pong:v2.7
```

Test the endpoints:
```bash
# Increment counter
curl http://localhost:3000/pingpong
# Output: pong 0

curl http://localhost:3000/pingpong
# Output: pong 1

# Get current count
curl http://localhost:3000/count
# Output: {"count":2}
```

## Deploying to Kubernetes

### Local Deployment (k3d)

#### Prerequisites
- k3d cluster running
- kubectl configured
- PostgreSQL StatefulSet deployed

#### Deploy PostgreSQL StatefulSet
```bash
kubectl apply -f manifests/statefulset.yaml
```

Wait for PostgreSQL to be ready:
```bash
kubectl wait --for=condition=ready pod -l app=postgres-ss -n exercises --timeout=300s
```

#### Import image to k3d
```bash
k3d image import ping-pong:v2.7 -c k3s-default
```

#### Deploy Ping-Pong Application
```bash
kubectl apply -f manifests/deployment.yaml
kubectl apply -f manifests/service.yaml
kubectl apply -f manifests/ingress.yaml
```

#### Check deployment
```bash
kubectl get deployments -n exercises
kubectl get pods -n exercises
kubectl get services -n exercises
```

#### View logs
```bash
kubectl logs -f deployment/ping-pong-dep -n exercises
```

#### Test the application
If using Ingress:
```bash
curl http://localhost:8081/pingpong
curl http://localhost:8081/count
```

Or using port-forward:
```bash
kubectl port-forward -n exercises deployment/ping-pong-dep 3000:3000
curl http://localhost:3000/pingpong
curl http://localhost:3000/count
```

#### Delete deployment
```bash
kubectl delete -f manifests/deployment.yaml
kubectl delete -f manifests/service.yaml
kubectl delete -f manifests/ingress.yaml
kubectl delete -f manifests/statefulset.yaml
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
docker build -t ping-pong:v2.7 .
docker tag ping-pong:v2.7 $ACR_NAME.azurecr.io/ping-pong:v2.7

# Push to ACR
docker push $ACR_NAME.azurecr.io/ping-pong:v2.7
```

#### Update manifests for AKS
Update the image reference in `manifests/deployment.yaml`:
```yaml
spec:
  containers:
    - name: ping-pong
      image: <acr-name>.azurecr.io/ping-pong:v2.7
```

#### Deploy to AKS
```bash
# Deploy PostgreSQL StatefulSet
kubectl apply -f manifests/statefulset.yaml

# Wait for PostgreSQL
kubectl wait --for=condition=ready pod -l app=postgres-ss -n exercises --timeout=300s

# Deploy application
kubectl apply -f manifests/deployment.yaml
kubectl apply -f manifests/service.yaml
```

#### Expose with LoadBalancer (Exercise 3.1)
For Exercise 3.1, use a LoadBalancer service instead of Ingress:

Create `manifests/service-loadbalancer.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  namespace: exercises
  name: ping-pong-lb
spec:
  type: LoadBalancer
  selector:
    app: ping-pong
  ports:
    - port: 80
      protocol: TCP
      targetPort: 3000
```

Deploy and get external IP:
```bash
kubectl apply -f manifests/service-loadbalancer.yaml

# Wait for external IP (may take 2-3 minutes)
kubectl get svc -n exercises --watch

# Test once you have the external IP
curl http://<EXTERNAL-IP>/pingpong
```

#### Check AKS deployment
```bash
kubectl get deployments -n exercises
kubectl get pods -n exercises
kubectl get services -n exercises
```

#### View logs in AKS
```bash
kubectl logs -f deployment/ping-pong-dep -n exercises
```

#### Delete AKS deployment
```bash
kubectl delete -f manifests/deployment.yaml
kubectl delete -f manifests/service-loadbalancer.yaml
kubectl delete -f manifests/statefulset.yaml
```

## Configuration

The application uses the following environment variables:

- `PORT` - The port number the server listens on (default: 3000)
- `POSTGRES_URL` - PostgreSQL connection string (default: postgres://postgres:example@postgres-svc:5432/postgres?sslmode=disable)

## Database Schema

The application automatically creates the following table:

```sql
CREATE TABLE IF NOT EXISTS counter (
    id INTEGER PRIMARY KEY,
    count INTEGER NOT NULL
);
```

## API Endpoints

### POST/GET /pingpong
Increments the counter and returns the previous count.

**Response:**
```
pong 0
```

### GET /count
Returns the current counter value as JSON.

**Response:**
```json
{
  "count": 1
}
```

## Files

- `main.go` - Main application code with database integration
- `go.mod` - Go module definition
- `Dockerfile` - Multi-stage Docker build
- `manifests/deployment.yaml` - Kubernetes deployment in "exercises" namespace
- `manifests/service.yaml` - ClusterIP service
- `manifests/ingress.yaml` - Ingress configuration for external access
- `manifests/statefulset.yaml` - PostgreSQL StatefulSet with persistent storage

## Technologies

- **Language:** Go 1.21
- **Database:** PostgreSQL 13.0
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes
- **Storage:** StatefulSet with PersistentVolumeClaim

## Integration with Log-Output

The ping-pong application provides a `/count` endpoint that is consumed by the log-output application (Exercise 2.1). The log-output service fetches the current count via HTTP and displays it in its status page.
