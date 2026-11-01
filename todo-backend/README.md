# Todo Backend

REST API that stores the todos of the [todo project](../todo-project). The todos are stored in a PostgreSQL database that runs as a StatefulSet, so they survive restarts of both the backend and the database.

## Exercises

- **2.2:** `GET /todos` lists the todos, `POST /todos` creates one
- **2.4:** Runs in the namespace `project` with the todo project
- **2.6:** Port and maximum todo length come from env variables (Deployment and ConfigMap)
- **2.8:** Todos are saved in Postgres, run as a StatefulSet with one replica; the backend reaches it through a ConfigMap and a Secret
- **2.10:** Request logging as JSON lines, including every created and every rejected todo, searchable in Grafana
- **4.2:** Readiness probe on `/healthz` (the database), liveness probe on `/livez`; starts without a database and connects in the background
- **4.5:** `PUT /todos/<id>` marks a todo as done or not done
- **4.6:** Publishes a message to NATS when a todo is saved or updated, for [broadcaster](../broadcaster)

## API

| Method | Path | Body | Response |
|---|---|---|---|
| `GET` | `/todos` | | `200` with a JSON array of todos |
| `POST` | `/todos` | `{"text": "Buy milk"}` | `201` with the new todo; `400` if the text is empty or over `TODO_MAX_LENGTH` (140) characters |
| `PUT` | `/todos/<id>` | `{"done": true}` | `200` with the updated todo; `400` if the body has no `done`; `404` if there is no such todo |
| `GET` | `/healthz` | | `200 OK` when the database is connected and answers, `503` otherwise (readiness probe) |
| `GET` | `/livez` | | `200 OK` while the process runs (liveness probe) |

A todo looks like this:

```json
{"id": "6f1c…", "text": "Buy milk", "done": false, "created_at": "2026-10-25T18:40:00Z"}
```

The text is trimmed, and the limit counts characters, not bytes.

## Running locally

```bash
docker network create todo
docker run -d --name postgres --network todo -e POSTGRES_USER=todo -e POSTGRES_DB=todos \
  -e POSTGRES_PASSWORD=local postgres:18.6-alpine3.24
docker build -t todo-backend:v4.6 .
docker run -p 3000:3000 --network todo -e PORT=3000 -e TODO_MAX_LENGTH=140 -e NATS_URL=nats://nats:4222 \
  -e PGHOST=postgres -e PGPORT=5432 -e PGDATABASE=todos -e PGUSER=todo -e PGPASSWORD=local -e PGSSLMODE=disable \
  todo-backend:v4.6
curl -X POST -H 'Content-Type: application/json' -d '{"text":"Buy milk"}' http://localhost:3000/todos
curl http://localhost:3000/todos
curl -X PUT -H 'Content-Type: application/json' -d '{"done":true}' http://localhost:3000/todos/<id>
```

## Deploying to Kubernetes (k3d)

Deploy together with the todo project (see [its README](../todo-project/README.md#deploying-to-kubernetes-k3d)):

```bash
kubectl apply -f ../manifests/project-namespace.yaml
sops --decrypt ../../k8s-course-config/apps/todo-backend/secrets/postgres-secret.enc.yaml | kubectl apply -f -   # database password
docker build -t todo-backend:v4.6 .
k3d image import todo-backend:v4.6
kubectl apply -k ../../k8s-course-config/overlays/k3d                  # the whole project with Kustomize
curl http://localhost:8081/api/todos
```

Traffic flow in k3d: `localhost:8081/api/todos` → Traefik Ingress `todo-backend-ingress` (`/api`) → Traefik Middleware `strip-api` removes `/api` → ClusterIP Service `todo-backend-svc:2345` → pod port `3000` → `/todos`.

The Ingress refers to the Middleware as `project-strip-api@kubernetescrd` (`<namespace>-<name>`), so the annotation has to change whenever the namespace does. In AKS the HTTPRoute of [`../project/aks`](https://github.com/priku/k8s-course-config/tree/main/components/gateway/gateway.yaml) strips `/api` with a `URLRewrite` filter instead.

## Logging

The backend logs JSON lines to stdout with Go's `log/slog`, so that Loki can parse the fields (LogQL `| json`):

- every request: `msg="request"` with `method`, `path`, `status`, `duration_ms`, `remote_addr` and `user_agent` (successful `/healthz` calls are left out)
- every new todo: `msg="todo created"` with `id`, `length` and `text`
- every refused todo: `msg="todo rejected"`, level `WARN`, with `reason` (`too_long`, `empty` or `invalid_body`), `length`, `max_length` and `text`

```
{"level":"WARN","msg":"todo rejected","reason":"too_long","length":162,"max_length":140,"text":"This todo is far too long ..."}
{"level":"INFO","msg":"request","method":"POST","path":"/todos","status":400,"duration_ms":0,...}
```

To test the 140 character limit and find the refused todo in Grafana (see [`../monitoring`](../monitoring)):

```bash
curl -i -H 'Content-Type: application/json' \
  -d "{\"text\":\"$(printf 'x%.0s' {1..141})\"}" http://localhost:8081/api/todos   # 400 Bad Request
kubectl logs -n project deployment/todo-backend-dep | grep 'todo rejected'
# Grafana Explore, Loki: {namespace="project", container="todo-backend"} |= "todo rejected"
```

## Database

[`apps/todo-backend/statefulset.yaml`](https://github.com/priku/k8s-course-config/tree/main/apps/todo-backend/statefulset.yaml) in the configuration repository runs Postgres (`postgres:18.6-alpine3.24`) as the StatefulSet `postgres-ss` with one replica and the headless Service `postgres-svc`, the same way as the [ping-pong database](../ping-pong/README.md#database) in the `exercises` namespace. The two databases are separate: same names, different namespaces.

- The pod `postgres-ss-0` gets its own PersistentVolumeClaim `data-postgres-ss-0` from `volumeClaimTemplates` (the cluster's default storage class: `local-path` in k3d, an Azure Disk in AKS), mounted at `/var/lib/postgresql` as Postgres 18 expects.
- The backend creates the table `todos` at startup. The database fills in `id` (`gen_random_uuid()`), `done` and `created_at`, and `INSERT ... RETURNING` gives the new row back. The text is passed as a query parameter, never put into the SQL.
- The backend retries the connection for a minute at startup, since the headless Service only resolves once Postgres is ready. `GET /healthz` answers `503` if the database is not reachable.

How the backend reaches the database, without any connection string in the code:

| Setting | Where | Used by |
|---|---|---|
| `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, `PGSSLMODE` | ConfigMap `todo-backend-config` | the backend (pgx reads the standard `PG*` variables); Postgres creates `POSTGRES_USER` and `POSTGRES_DB` from the same keys |
| `POSTGRES_PASSWORD` | Secret `postgres-secret`, encrypted with SOPS in `apps/todo-backend/secrets/postgres-secret.enc.yaml` of the configuration repository | Postgres, and the backend as `PGPASSWORD` |

The Secret is encrypted for the age key in [`../.sops.yaml`](../.sops.yaml); see the [ping-pong README](../ping-pong/README.md#database-password) for how SOPS is set up and how to make your own key and password.

To look at the todos in the database:

```bash
kubectl run psql-debug -n project --rm -it --restart=Never --image postgres:18.6-alpine3.24 \
  --overrides='{"spec":{"containers":[{"name":"psql-debug","image":"postgres:18.6-alpine3.24","stdin":true,"tty":true,"command":["psql"],"envFrom":[{"configMapRef":{"name":"todo-backend-config"}}],"env":[{"name":"PGPASSWORD","valueFrom":{"secretKeyRef":{"name":"postgres-secret","key":"POSTGRES_PASSWORD"}}}]}]}}'
todos=# SELECT * FROM todos;
```

## Configuration

There are no defaults in the code; the app stops at startup if a variable is missing.

- `PORT` - HTTP server port, `3000` in the Deployment
- `TODO_MAX_LENGTH` - maximum todo length in characters, `140` in the ConfigMap `todo-backend-config`. The [todo project](../todo-project) reads the same key for its input, so the limit is changed in one place:

```bash
kubectl edit configmap todo-backend-config -n project
kubectl rollout restart deployment/todo-backend-dep deployment/todo-project-dep -n project
```

- `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, `PGSSLMODE` - database connection, in the ConfigMap `todo-backend-config`
- `PGPASSWORD` - database password, from the Secret `postgres-secret`

## Files

- `main.go` - HTTP API with the todos in Postgres
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- The Kubernetes configuration is in [`apps/todo-backend`](https://github.com/priku/k8s-course-config/tree/main/apps/todo-backend) of the configuration repository

## Technologies

- **Language:** Go 1.27, [pgx](https://github.com/jackc/pgx) v5
- **Database:** PostgreSQL 18
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress and Middleware (k3s default)
- **Monitoring:** Loki and Grafana (see [`../monitoring`](../monitoring))
