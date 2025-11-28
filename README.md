# DevOps with Kubernetes - Exercises

This repository contains my solutions for the [DevOps with Kubernetes](https://devopswithkubernetes.com/) course.

## Exercises

### Chapter 1: Getting Started
- [1.1](https://github.com/priku/k8s-course/tree/1.1/log-output) - Log Output: Application that generates and outputs a random string with timestamp every 5 seconds
- [1.2](https://github.com/priku/k8s-course/tree/1.2/todo-project) - Todo Project: Web server with configurable PORT environment variable
- [1.3](https://github.com/priku/k8s-course/tree/1.3/log-output) - Log Output: Declarative deployment with Kubernetes
- [1.4](https://github.com/priku/k8s-course/tree/1.4/todo-project) - Todo Project: Declarative deployment with Kubernetes
- [1.5](https://github.com/priku/k8s-course/tree/1.5/todo-project) - Todo Project: HTTP server responding to GET requests with HTML page

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
