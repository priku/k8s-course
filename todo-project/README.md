# Todo Project Application

A simple Go web server with configurable PORT via environment variable. This is the course project that will be extended throughout the DevOps with Kubernetes course.

## Exercise 1.2 & 1.4

This application fulfills the requirements for:
- **Exercise 1.2:** Creating a web server that outputs "Server started in port NNNN" with configurable PORT
- **Exercise 1.4:** Creating declarative deployment configuration for Kubernetes

## Features

- Simple HTTP web server
- Configurable port via `PORT` environment variable (defaults to 3000)
- Outputs startup message with the configured port
- Containerized with multi-stage Docker build

## Running Locally

### Using Docker

Build the image:
```bash
docker build -t todo-project:v1.0 .
```

Run with default port (3000):
```bash
docker run -p 3000:3000 todo-project:v1.0
```

Run with custom port:
```bash
docker run -e PORT=8080 -p 8080:8080 todo-project:v1.0
```

Expected output:
```
Server started in port 3000
```

## Deploying to Kubernetes

### Prerequisites
- Kubernetes cluster (k3d, minikube, or similar)
- kubectl configured

### Import image to k3d (if using k3d)
```bash
k3d image import todo-project:v1.0
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
kubectl logs -f deployment/todo-project-dep
```

### Delete deployment
```bash
kubectl delete -f manifests/deployment.yaml
```

## Configuration

The application uses the following environment variable:

- `PORT` - The port number the server listens on (default: 3000)

The Kubernetes deployment is configured with `PORT=3000` in the deployment manifest.

## Files

- `main.go` - Main application code
- `go.mod` - Go module definition
- `Dockerfile` - Multi-stage Docker build
- `manifests/deployment.yaml` - Kubernetes deployment configuration with environment variables

## Technologies

- **Language:** Go 1.21
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes

## Future Enhancements

This is the starting point for the course project. Future exercises will add:
- Todo CRUD functionality
- Database integration
- Additional microservices
- Monitoring and logging
- And more!
