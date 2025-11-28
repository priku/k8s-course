# Log Output Application

A simple Go application that generates a random UUID on startup and outputs it with a timestamp every 5 seconds.

## Exercise 1.1 & 1.3

This application fulfills the requirements for:
- **Exercise 1.1:** Creating and deploying an application that outputs a random string with timestamp
- **Exercise 1.3:** Moving to declarative approach with Kubernetes manifests

## Features

- Generates a random UUID on application startup
- Outputs the UUID with RFC3339 timestamp every 5 seconds
- Containerized with multi-stage Docker build for small image size

## Running Locally

### Using Docker

Build the image:
```bash
docker build -t log-output:v1.0 .
```

Run the container:
```bash
docker run log-output:v1.0
```

Expected output:
```
Log output application started
2025-11-28T22:13:21.991206053Z: 89e43ff3-dd56-4e1b-b92a-1cb0e1ae02fd
2025-11-28T22:13:26.991325185Z: 89e43ff3-dd56-4e1b-b92a-1cb0e1ae02fd
...
```

## Deploying to Kubernetes

### Prerequisites
- Kubernetes cluster (k3d, minikube, or similar)
- kubectl configured

### Import image to k3d (if using k3d)
```bash
k3d image import log-output:v1.0
```

### Deploy
```bash
kubectl apply -f manifests/deployment.yaml
```

### Check deployment
```bash
kubectl get deployments
kubectl get pods
```

### View logs
```bash
kubectl logs -f deployment/log-output-dep
```

### Delete deployment
```bash
kubectl delete -f manifests/deployment.yaml
```

## Files

- `main.go` - Main application code
- `go.mod` - Go module definition
- `Dockerfile` - Multi-stage Docker build
- `manifests/deployment.yaml` - Kubernetes deployment configuration

## Technologies

- **Language:** Go 1.21
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes
