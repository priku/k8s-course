# Ping-pong Application

A Go HTTP server that responds to `GET /pingpong` with `pong N`, where N is the number of earlier requests. `GET /count` returns the current count as JSON (`{"count":N}`), which the [Log output application](../log-output) fetches over HTTP. The count is kept in memory, so it resets when the pod restarts.

## Exercises

- **1.9:** Responds with "pong 0", "pong 1", …; shares an Ingress with the [Log output application](../log-output) and is routed at `/pingpong`
- **1.11:** Saves the request count to a PersistentVolume shared with Log output
- **2.1:** `GET /count` endpoint for Log output instead of the shared volume

## Running locally

```bash
docker build -t ping-pong:v2.1 .
docker run -p 3000:3000 ping-pong:v2.1
curl http://localhost:3000/pingpong   # pong 0
curl http://localhost:3000/pingpong   # pong 1
curl http://localhost:3000/count      # {"count":2}
```

## Deploying to Kubernetes (k3d)

The cluster needs host port 8081 mapped to the load balancer (Traefik Ingress):

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Deploy together with Log output, which serves `/`:

```bash
docker build -t ping-pong:v2.1 .
k3d image import ping-pong:v2.1
kubectl apply -f manifests/
kubectl apply -f ../log-output/manifests/
curl http://localhost:8081/pingpong   # pong 0
```

The app listens on `/pingpong` itself, so the Ingress needs no path rewriting.

Traffic flow: `localhost:8081/pingpong` → Traefik Ingress `app-ingress` → ClusterIP Service `ping-pong-svc:2345` → pod port `3000`. Log output calls `http://ping-pong-svc:2345/count` inside the cluster; `/count` is not routed by the Ingress.

## Configuration

- `PORT` - HTTP server port (default: 3000)

## Files

- `main.go` - HTTP server with an in-memory counter (`/pingpong`, `/count`)
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with `PORT=3000`
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)
- `manifests/ingress.yaml` - Shared Ingress: `/pingpong` → ping-pong, `/` → Log output

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress (k3s default)
