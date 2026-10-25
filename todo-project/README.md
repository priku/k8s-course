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

## Running locally

```bash
docker build -t todo-project:v2.2 .
docker run -p 3000:3000 -v "$PWD/files:/usr/src/app/files" todo-project:v2.2
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
docker build -t todo-project:v2.2 .
docker build -t todo-backend:v2.2 ../todo-backend
k3d image import todo-project:v2.2 todo-backend:v2.2
kubectl apply -f manifests/
kubectl apply -f ../todo-backend/manifests/
# Open http://localhost:8081, add a todo and reload: it stays in the list
kubectl get all,ingress,pvc -n project
```

Traffic flow:

- `localhost:8081/` → Traefik Ingress `todo-project-ingress` (`/`) → ClusterIP Service `todo-project-svc:2345` → pod port `3000` (HTML page, `/image`)
- `localhost:8081/api/todos` → Ingress `todo-backend-ingress` (`/api`, prefix stripped) → `todo-backend-svc:2345` → todo-backend pod port `3000` (`GET`/`POST /todos`)

The browser talks to both apps through the same Ingress and origin, so no CORS is needed. Log output and ping-pong live at `/log-output` and `/pingpong` in the namespace `exercises`. Ingresses from different namespaces share the same Traefik, so the paths must not overlap.

## Todo list

The page has an input for a new todo and a **Send** button. The input accepts at most 140 characters (`maxlength="140"`), shows a character counter and the button refuses empty or too long todos.

On load the page calls `GET /api/todos` and renders the list; **Send** posts `{"text": "..."}` to `POST /api/todos` and reloads the list. Todo texts are inserted with `textContent`, so HTML in a todo is shown as text instead of being run. The todos are kept in todo-backend's memory, so they disappear when its pod restarts; a database comes in a later exercise.

## Image caching

The page shows a random 1200px image from [Lorem Picsum](https://picsum.photos/1200), served by the app at `/image`:

- The image and its download time are stored in `/usr/src/app/files` (`daily-image.jpg`, `image-timestamp.txt`).
- The directory is a PersistentVolumeClaim (`todo-project-image-pvc`, k3s `local-path` storage class), so the cached image survives container crashes and pod restarts without calling the API again.
- When the image is older than 10 minutes, the next request downloads a new one. Downloads go to a temporary file that is renamed into place, and a lock keeps concurrent requests from downloading at the same time.

To test that the cache survives a crash, kill the container and check that the same image is still shown:

```bash
kubectl exec deployment/todo-project-dep -n project -- kill 1
kubectl logs deployment/todo-project-dep -n project   # "Image is ... old, still valid"
```

## Files

- `main.go` - HTTP server, image download and caching, page with the todo list JavaScript
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with `PORT=3000`, image cache volume at `/usr/src/app/files`
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)
- `manifests/ingress.yaml` - Ingress routing `/` to the Service
- `manifests/persistentvolumeclaim.yaml` - 100Mi claim for the image cache

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress (k3s default)
