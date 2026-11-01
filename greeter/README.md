# Greeter

A small Go service that answers `GET /` with a greeting, for [Log output](../log-output) in the [service mesh](../service-mesh) (exercise 5.3). There are two versions of it, which differ only by the env variable `VERSION`, and the mesh splits the traffic between them.

```bash
docker build -t greeter:v5.3 .
docker run -p 3000:3000 -e PORT=3000 -e VERSION=1 greeter:v5.3
curl http://localhost:3000/    # hello from version 1
```

## Files

- `main.go` - The HTTP server
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - A Deployment for each version
- `manifests/service.yaml` - `greeter-svc` for the callers, and a Service for each version
- `manifests/httproute.yaml` - The 75/25 split
