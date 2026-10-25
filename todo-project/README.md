# Todo Project Application

The course project: a Go web server that will grow into a todo application during the DevOps with Kubernetes course.

## Exercises

- **1.2:** Web server that prints "Server started in port NNNN", port chosen with the `PORT` environment variable
- **1.4:** Declarative deployment (`manifests/deployment.yaml`)
- **1.5:** Responds to `GET /` with an HTML page

## Running locally

```bash
docker build -t todo-project:v1.5 .
docker run -p 3000:3000 todo-project:v1.5
# Open http://localhost:3000
```

The port can be changed with `PORT`, e.g. `docker run -e PORT=8080 -p 8080:8080 todo-project:v1.5`.

## Deploying to Kubernetes (k3d)

```bash
docker build -t todo-project:v1.5 .
k3d image import todo-project:v1.5
kubectl apply -f manifests/deployment.yaml
kubectl logs -f deployment/todo-project-dep   # Server started in port 3000
```

### Accessing the app with port-forward

```bash
kubectl port-forward deployment/todo-project-dep 3003:3000
# Open http://localhost:3003
```

## Files

- `main.go` - HTTP server
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with `PORT=3000`

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes
