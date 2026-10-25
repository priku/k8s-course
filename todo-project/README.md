# Todo Project Application

The course project's frontend: a Go web server that serves the HTML page and the cached random image. The todos themselves live in [todo-backend](../todo-backend); the page fetches and creates them with JavaScript through the Ingress path `/api`.

## Exercises

- **1.2:** Web server that prints "Server started in port NNNN", port chosen with the `PORT` environment variable
- **1.4:** Declarative deployment (`manifests/deployment.yaml`)
- **1.5:** Responds to `GET /` with an HTML page
- **1.6:** NodePort Service for access from outside the cluster
- **1.8:** Ingress instead of NodePort
- **1.12:** Random image from Lorem Picsum, cached on a PersistentVolume and refreshed every 10 minutes
- **1.13:** Todo input (max 140 characters), send button and a list of hardcoded todos
- **2.2:** Todos are fetched from and created in [todo-backend](../todo-backend) instead of being hardcoded
- **2.4:** Everything of the project runs in its own namespace `project`
- **2.6:** No hardcoded ports, URLs or other settings: all configuration comes from env variables set in ConfigMaps and the Deployment
- **2.8:** The todos are saved in Postgres by [todo-backend](../todo-backend), so they survive restarts
- **2.9:** [wiki-todo-generator](../wiki-todo-generator) adds a `Read <random Wikipedia URL>` todo every hour; Enter in the input now sends the todo

## Running locally

```bash
docker build -t todo-project:v2.9 .
docker run -p 3000:3000 -e PORT=3000 -e IMAGE_DIR=/usr/src/app/files \
  -e IMAGE_URL=https://picsum.photos/1200 -e IMAGE_REFRESH_INTERVAL=10m \
  -e BACKEND_URL=/api -e TODO_MAX_LENGTH=140 todo-project:v2.9
# Open http://localhost:3000 (the todo list needs the backend at /api, so use the cluster for that)
```

## Deploying to Kubernetes (k3d)

The cluster needs host port 8081 mapped to the load balancer (Traefik Ingress):

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Everything of the project, including todo-backend, runs in the namespace `project` (set in each manifest's `metadata.namespace`). Create it first:

```bash
kubectl apply -f ../manifests/project-namespace.yaml
```

Build, import and deploy together with todo-backend:

```bash
sops --decrypt ../todo-backend/secrets/postgres-secret.enc.yaml | kubectl apply -f -   # database password
docker build -t todo-project:v2.9 .
docker build -t todo-backend:v2.10 ../todo-backend
docker build -t wiki-todo-generator:v2.9 ../wiki-todo-generator
k3d image import todo-project:v2.9 todo-backend:v2.10 wiki-todo-generator:v2.9
kubectl apply -f manifests/
kubectl apply -f ../todo-backend/manifests/          # backend + Postgres
kubectl apply -f ../wiki-todo-generator/manifests/   # hourly CronJob
# Open http://localhost:8081, add a todo and reload: it stays in the list
kubectl get all,ingress,pvc -n project
```

Traffic flow:

- `localhost:8081/` → Traefik Ingress `todo-project-ingress` (`/`) → ClusterIP Service `todo-project-svc:2345` → pod port `3000` (HTML page, `/image`)
- `localhost:8081/api/todos` → Ingress `todo-backend-ingress` (`/api`, prefix stripped) → `todo-backend-svc:2345` → todo-backend pod port `3000` (`GET`/`POST /todos`)

The browser talks to both apps through the same Ingress and origin, so no CORS is needed. Log output and ping-pong live at `/log-output` and `/pingpong` in the namespace `exercises`. Ingresses from different namespaces share the same Traefik, so the paths must not overlap.

## Todo list

The page has an input for a new todo and a **Send** button (or Enter). The input accepts at most `TODO_MAX_LENGTH` characters (140, set as `maxlength`), shows a character counter and the button refuses empty or too long todos.

On load the page calls `GET /api/todos` and renders the list; **Send** posts `{"text": "..."}` to `POST /api/todos` and reloads the list. Todo texts are inserted with `textContent`, so HTML in a todo is shown as text instead of being run. The todos are stored by todo-backend in Postgres, so they survive restarts (see [its README](../todo-backend/README.md#database)).

## Image caching

The page shows a random image from `IMAGE_URL` (a 1200px image from [Lorem Picsum](https://picsum.photos/1200)), served by the app at `/image`:

- The image and its download time are stored in `IMAGE_DIR`, `/usr/src/app/files` (`daily-image.jpg`, `image-timestamp.txt`).
- The directory is a PersistentVolumeClaim (`todo-project-image-pvc`, k3s `local-path` storage class), so the cached image survives container crashes and pod restarts without calling the API again.
- When the image is older than `IMAGE_REFRESH_INTERVAL` (10 minutes), the next request downloads a new one. Downloads go to a temporary file that is renamed into place, and a lock keeps concurrent requests from downloading at the same time.

To test that the cache survives a crash, kill the container and check that the same image is still shown:

```bash
kubectl exec deployment/todo-project-dep -n project -- kill 1
kubectl logs deployment/todo-project-dep -n project   # "Image is ... old, still valid"
```

## Configuration

The app has no defaults in the code: every setting is an env variable, and the app stops at startup with `Environment variable X is not set` if one is missing.

| Variable | Value | Set in |
|---|---|---|
| `PORT` | `3000` | `manifests/deployment.yaml` (belongs with the Service's `targetPort`) |
| `IMAGE_DIR` | `/usr/src/app/files` | `manifests/deployment.yaml` (belongs with the `volumeMounts` path) |
| `IMAGE_URL` | `https://picsum.photos/1200` | ConfigMap `todo-project-config` |
| `IMAGE_REFRESH_INTERVAL` | `10m` (a Go duration) | ConfigMap `todo-project-config` |
| `BACKEND_URL` | `/api` (todo-backend as the browser sees it) | ConfigMap `todo-project-config` |
| `TODO_MAX_LENGTH` | `140` | ConfigMap `todo-backend-config` of [todo-backend](../todo-backend), so both apps check the same limit |

`envFrom` brings in every key of `todo-project-config`, and `configMapKeyRef` picks the single key from the backend's ConfigMap. The page is an `html/template`, so the values are escaped correctly in HTML and JavaScript.

Env variables are read at startup, so restart after changing a ConfigMap:

```bash
kubectl edit configmap todo-project-config -n project
kubectl rollout restart deployment/todo-project-dep -n project
```

## Files

- `main.go` - HTTP server, image download and caching, page template with the todo list JavaScript
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with `PORT` and `IMAGE_DIR`, env from the ConfigMaps, image cache volume at `/usr/src/app/files`
- `manifests/configmap.yaml` - ConfigMap `todo-project-config` (image URL and interval, backend URL)
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)
- `manifests/ingress.yaml` - Ingress routing `/` to the Service
- `manifests/persistentvolumeclaim.yaml` - 100Mi claim for the image cache

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress (k3s default)
