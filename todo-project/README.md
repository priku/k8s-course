# Todo Project Application

The course project: a Go web server that will grow into a todo application during the DevOps with Kubernetes course.

## Exercises

- **1.2:** Web server that prints "Server started in port NNNN", port chosen with the `PORT` environment variable
- **1.4:** Declarative deployment (`manifests/deployment.yaml`)
- **1.5:** Responds to `GET /` with an HTML page
- **1.6:** NodePort Service for access from outside the cluster
- **1.8:** Ingress instead of NodePort
- **1.12:** Random image from Lorem Picsum, cached on a PersistentVolume and refreshed every 10 minutes
- **1.13:** Todo input (max 140 characters), send button and a list of hardcoded todos

## Running locally

```bash
docker build -t todo-project:v1.13 .
docker run -p 3000:3000 -v "$PWD/files:/usr/src/app/files" todo-project:v1.13
# Open http://localhost:3000
```

## Deploying to Kubernetes (k3d)

The cluster needs host port 8081 mapped to the load balancer (Traefik Ingress):

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Build, import and deploy:

```bash
docker build -t todo-project:v1.13 .
k3d image import todo-project:v1.13
kubectl apply -f manifests/
# Open http://localhost:8081
```

If the Log output / ping-pong Ingress (`app-ingress`, exercise 1.9) is in the cluster, delete it first so the two don't compete for `/`:

```bash
kubectl delete ingress app-ingress
```

Traffic flow: `localhost:8081` → Traefik Ingress `todo-project-ingress` (`/`) → ClusterIP Service `todo-project-svc:2345` → pod port `3000`.

## Todo list

The page has an input for a new todo and a **Send** button, plus a list of hardcoded todos. The input accepts at most 140 characters (`maxlength="140"`), shows a character counter and the button refuses empty or too long todos. Sending does not store anything yet; the backend comes in a later exercise.

## Image caching

The page shows a random 1200px image from [Lorem Picsum](https://picsum.photos/1200), served by the app at `/image`:

- The image and its download time are stored in `/usr/src/app/files` (`daily-image.jpg`, `image-timestamp.txt`).
- The directory is a PersistentVolumeClaim (`todo-project-image-pvc`, k3s `local-path` storage class), so the cached image survives container crashes and pod restarts without calling the API again.
- When the image is older than 10 minutes, the next request downloads a new one. Downloads go to a temporary file that is renamed into place, and a lock keeps concurrent requests from downloading at the same time.

To test that the cache survives a crash, kill the container and check that the same image is still shown:

```bash
kubectl exec deployment/todo-project-dep -- kill 1
kubectl logs deployment/todo-project-dep   # "Image is ... old, still valid"
```

## Files

- `main.go` - HTTP server, image download and caching
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with `PORT=3000`, image cache volume at `/usr/src/app/files`
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)
- `manifests/ingress.yaml` - Ingress routing `/` to the Service
- `manifests/persistentvolumeclaim.yaml` - 100Mi claim for the image cache

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress (k3s default)
