# Todo Project Application

The course project: a Go web server that will grow into a todo application during the DevOps with Kubernetes course.

## Exercises

- **1.2:** Web server that prints "Server started in port NNNN", port chosen with the `PORT` environment variable
- **1.4:** Declarative deployment (`manifests/deployment.yaml`)
- **1.5:** Responds to `GET /` with an HTML page
- **1.6:** NodePort Service for access from outside the cluster
- **1.8:** Ingress instead of NodePort

## Running locally

```bash
docker build -t todo-project:v1.8 .
docker run -p 3000:3000 todo-project:v1.8
# Open http://localhost:3000
```

## Deploying to Kubernetes (k3d)

The cluster needs host port 8081 mapped to the load balancer (Traefik Ingress):

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Build, import and deploy:

```bash
docker build -t todo-project:v1.8 .
k3d image import todo-project:v1.8
kubectl apply -f manifests/
# Open http://localhost:8081
```

If the Log output / ping-pong Ingress (`app-ingress`, exercise 1.9) is in the cluster, delete it first so the two don't compete for `/`:

```bash
kubectl delete ingress app-ingress
```

Traffic flow: `localhost:8081` → Traefik Ingress `todo-project-ingress` (`/`) → ClusterIP Service `todo-project-svc:2345` → pod port `3000`.

## Files

- `main.go` - HTTP server
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with `PORT=3000`
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)
- `manifests/ingress.yaml` - Ingress routing `/` to the Service

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress (k3s default)
