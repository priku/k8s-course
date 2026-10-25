# Todo Project Application

The course project: a Go web server that will grow into a todo application during the DevOps with Kubernetes course.

## Exercises

- **1.2:** Web server that prints "Server started in port NNNN", port chosen with the `PORT` environment variable
- **1.4:** Declarative deployment (`manifests/deployment.yaml`)
- **1.5:** Responds to `GET /` with an HTML page
- **1.6:** NodePort Service for access from outside the cluster

## Running locally

```bash
docker build -t todo-project:v1.5 .
docker run -p 3000:3000 todo-project:v1.5
# Open http://localhost:3000
```

## Deploying to Kubernetes (k3d)

The cluster needs host port 8082 mapped to NodePort 30080:

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Build, import and deploy:

```bash
docker build -t todo-project:v1.5 .
k3d image import todo-project:v1.5
kubectl apply -f manifests/
# Open http://localhost:8082
```

Traffic flow: `localhost:8082` → k3d → NodePort `30080` → Service `todo-project-svc:3000` → pod port `3000`.

## Files

- `main.go` - HTTP server
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with `PORT=3000`
- `manifests/service.yaml` - NodePort Service (`nodePort: 30080`)

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes
