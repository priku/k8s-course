# Todo Backend

REST API that stores the todos of the [todo project](../todo-project). The todos are kept in memory for now, so they disappear when the pod restarts; a database comes in a later exercise.

## Exercises

- **2.2:** `GET /todos` lists the todos, `POST /todos` creates one
- **2.4:** Runs in the namespace `project` with the todo project

## API

| Method | Path | Body | Response |
|---|---|---|---|
| `GET` | `/todos` | | `200` with a JSON array of todos |
| `POST` | `/todos` | `{"text": "Buy milk"}` | `201` with the new todo; `400` if the text is empty or over 140 characters |
| `GET` | `/healthz` | | `200 OK` |

A todo looks like this:

```json
{"id": "6f1c…", "text": "Buy milk", "done": false, "created_at": "2026-10-25T18:40:00Z"}
```

The text is trimmed, and the 140 limit counts characters, not bytes.

## Running locally

```bash
docker build -t todo-backend:v2.2 .
docker run -p 3000:3000 todo-backend:v2.2
curl -X POST -H 'Content-Type: application/json' -d '{"text":"Buy milk"}' http://localhost:3000/todos
curl http://localhost:3000/todos
```

## Deploying to Kubernetes (k3d)

Deploy together with the todo project (see [its README](../todo-project/README.md#deploying-to-kubernetes-k3d)):

```bash
kubectl apply -f ../manifests/project-namespace.yaml
docker build -t todo-backend:v2.2 .
k3d image import todo-backend:v2.2
kubectl apply -f manifests/
curl http://localhost:8081/api/todos
```

Traffic flow: `localhost:8081/api/todos` → Traefik Ingress `todo-backend-ingress` (`/api`) → Traefik Middleware `strip-api` removes `/api` → ClusterIP Service `todo-backend-svc:2345` → pod port `3000` → `/todos`.

The Ingress refers to the Middleware as `project-strip-api@kubernetescrd` (`<namespace>-<name>`), so the annotation has to change whenever the namespace does.

## Configuration

- `PORT` - HTTP server port (default: 3000)

## Files

- `main.go` - HTTP API with an in-memory todo list
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with `PORT=3000`
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)
- `manifests/ingress.yaml` - Ingress routing `/api` to the Service through the Middleware
- `manifests/middleware.yaml` - Traefik Middleware that strips the `/api` prefix

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress and Middleware (k3s default)
