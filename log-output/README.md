# Log Output Application

Two containers in one pod that share a file through an `emptyDir` volume:

- **writer** generates a random UUID on startup and every 5 seconds writes `<timestamp>: <uuid>` to `/usr/src/app/files/output.txt`
- **reader** serves the contents of that file over HTTP

## Exercises

- **1.1:** Random string generated on startup, printed with a timestamp every 5 seconds
- **1.3:** Declarative deployment (`manifests/deployment.yaml`)
- **1.7:** HTTP endpoint for the current status and an Ingress for browser access
- **1.9:** Shares an Ingress with the [ping-pong application](../ping-pong)
- **1.10:** Split into two containers in a single pod

## Deploying to Kubernetes (k3d)

The cluster needs host port 8081 mapped to the load balancer (Traefik Ingress):

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Build, import and deploy:

```bash
docker build -t log-writer:v1.10 writer
docker build -t log-reader:v1.10 reader
k3d image import log-writer:v1.10 log-reader:v1.10
kubectl apply -f manifests/
kubectl apply -f ../ping-pong/manifests/   # shared Ingress app-ingress
# Open http://localhost:8081
```

With more than one container in the pod, pick the container with `-c`:

```bash
kubectl logs -f deployment/log-output-dep -c log-writer
kubectl logs -f deployment/log-output-dep -c log-reader
```

Traffic flow: `localhost:8081` → Traefik Ingress `app-ingress` (`/`) → ClusterIP Service `log-output-svc:2345` → reader container port `3000` → reads `output.txt` from the shared volume.

The volume is an `emptyDir`, so the file lives as long as the pod: a new pod starts with an empty volume and a new random string.

## Configuration

- `PORT` - reader HTTP server port (default: 3000)

## Files

- `writer/` - Writes the timestamped random string to the shared file (Go 1.27, Alpine 3.24 with tzdata, Europe/Helsinki time)
- `reader/` - Serves the shared file over HTTP (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - One pod, two containers, shared `emptyDir` volume
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress (k3s default)
