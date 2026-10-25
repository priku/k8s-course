# Log Output Application

A Go application that generates a random UUID on startup, keeps it in memory and prints it with a timestamp (Europe/Helsinki) every 5 seconds. It also serves the current status over HTTP.

## Exercises

- **1.1:** Random string generated on startup, printed with a timestamp every 5 seconds
- **1.3:** Declarative deployment (`manifests/deployment.yaml`)
- **1.7:** HTTP endpoint for the current status (timestamp and random string) and an Ingress for browser access
- **1.9:** Shares an Ingress with the [ping-pong application](../ping-pong)

## Running locally

```bash
docker build -t log-output:v1.8 .
docker run -p 3000:3000 log-output:v1.8
# Logs: 2026-10-25T17:11:44+02:00: e2f97d37-4032-4768-adb8-b065cb6acf1d
# Status page: http://localhost:3000
```

## Deploying to Kubernetes (k3d)

The cluster needs host port 8081 mapped to the load balancer (Traefik Ingress):

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Build, import and deploy:

```bash
docker build -t log-output:v1.8 .
k3d image import log-output:v1.8
kubectl apply -f manifests/
kubectl apply -f ../ping-pong/manifests/   # shared Ingress app-ingress
kubectl logs -f deployment/log-output-dep
# Open http://localhost:8081 (log output) and http://localhost:8081/pingpong
```

The Ingress is shared with the ping-pong application and defined in [`../ping-pong/manifests/ingress.yaml`](../ping-pong/manifests/ingress.yaml): `/pingpong` goes to ping-pong, everything else to Log output.

Traffic flow: `localhost:8081` → Traefik Ingress `app-ingress` (`/`) → ClusterIP Service `log-output-svc:2345` → pod port `3000`.

## Configuration

- `PORT` - HTTP server port (default: 3000)

## Files

- `main.go` - Log loop and HTTP status endpoint
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24 with tzdata)
- `manifests/deployment.yaml` - Deployment with `PORT=3000`
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress (k3s default)
