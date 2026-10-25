# Log Output Application

Two containers in one pod that share an `emptyDir` volume:

- **writer** generates a random UUID on startup and every 5 seconds writes `<timestamp>: <uuid>` to `/usr/src/app/files/output.txt`
- **reader** serves that line over HTTP together with the request count it fetches from the [ping-pong application](../ping-pong) with `GET http://ping-pong-svc:2345/count`, preceded by the file `information.txt` and the env variable `MESSAGE` from the ConfigMap `log-output-config`:

```
file content: this text is from file
env variable: MESSAGE=hello world
2026-10-25T18:01:02+02:00: f35b159b-1c83-40b1-ada3-21c3029c87c2.
Ping / Pongs: 3
```

## Exercises

- **1.1:** Random string generated on startup, printed with a timestamp every 5 seconds
- **1.3:** Declarative deployment (`manifests/deployment.yaml`)
- **1.7:** HTTP endpoint for the current status and an Ingress for browser access
- **1.9:** Shares an Ingress with the [ping-pong application](../ping-pong)
- **1.10:** Split into two containers in a single pod
- **1.11:** Shares a PersistentVolume with ping-pong and shows its request count
- **2.1:** Gets the ping-pong count over HTTP instead of a shared volume
- **2.2:** Moved to the Ingress path `/log-output` so the todo project can have `/`
- **2.3:** Runs in the `exercises` namespace together with ping-pong
- **2.5:** ConfigMap with the file `information.txt` (mounted as a volume) and the env variable `MESSAGE`
- **3.2:** Runs in Azure Kubernetes Service with ping-pong, behind an Ingress at `/`
- **3.3:** Behind a Gateway API Gateway and HTTPRoute instead of the Ingress

## Deploying to Kubernetes (k3d)

The cluster needs host port 8081 mapped to the load balancer (Traefik Ingress):

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Everything runs in the namespace `exercises` (set in each manifest's `metadata.namespace`). Create it first:

```bash
kubectl apply -f ../manifests/exercises-namespace.yaml
```

Build, import and deploy together with ping-pong:

```bash
docker build -t log-writer:v1.11 writer
docker build -t log-reader:v2.5 reader
k3d image import log-writer:v1.11 log-reader:v2.5
kubectl apply -f manifests/
kubectl apply -f ../ping-pong/manifests/   # ping-pong + shared Ingress app-ingress
# Open http://localhost:8081/log-output, call http://localhost:8081/pingpong and reload
kubectl get all -n exercises
```

To skip `-n exercises` in the commands below, make it the default namespace with `kubens exercises` or `kubectl config set-context --current --namespace=exercises`.

With more than one container in the pod, pick the container with `-c`:

```bash
kubectl logs -f deployment/log-output-dep -c log-writer -n exercises
kubectl logs -f deployment/log-output-dep -c log-reader -n exercises
```

Traffic flow: `localhost:8081/log-output` → Traefik Ingress `app-ingress` (`/log-output`) → ClusterIP Service `log-output-svc:2345` → reader container port `3000` → reads `output.txt` from the `emptyDir` and calls `ping-pong-svc:2345/count` through cluster DNS. Both apps are in the same namespace, so the short Service name is enough.

Pod-to-pod requests can be tested by hand from a debugging pod. From another namespace, such as `default`, the Service is found at `<service>.<namespace>`:

```bash
kubectl run busybox --image=busybox --rm -it --restart=Never -- wget -qO - http://ping-pong-svc.exercises:2345/count
# {"count":3}
```

## Deploying to Azure Kubernetes Service (exercises 3.2–3.3)

`manifests/aks/` has the AKS versions of the manifests: the same ConfigMap and Service, and a Deployment whose images come from the course registry `acrdwkakssc.azurecr.io`. The Gateway and the HTTPRoute that route `/` here are in [ping-pong](../ping-pong/manifests/aks/gateway.yaml), and the cluster is described in [`../terraform`](../terraform).

```bash
# Images for the AKS nodes (amd64), built in Azure by ACR Tasks
az acr build --registry acrdwkakssc --image log-writer:v1.11 --platform linux/amd64 writer
az acr build --registry acrdwkakssc --image log-reader:v2.5 --platform linux/amd64 reader

kubectl config use-context aks-dwk-sc
kubectl apply -f manifests/aks/
kubectl apply -f ../ping-pong/manifests/aks/   # ping-pong, Postgres, the Gateway and the HTTPRoute
curl http://<Gateway ADDRESS>/                 # file content, env variable, log line, Ping / Pongs
```

## ConfigMap

`manifests/configmap.yaml` defines two keys, which the Deployment uses in two ways:

- `information.txt` is mounted as a file at `/etc/config/information.txt` in the reader (a `configMap` volume with `items`, so only that key becomes a file)
- `MESSAGE` becomes an env variable with `valueFrom.configMapKeyRef`

The reader reads the file on every request, so a changed ConfigMap shows up without a restart once the kubelet has synced the volume (usually within a minute or two). Env variables are set when the container starts, so a new `MESSAGE` needs a restart:

```bash
kubectl edit configmap log-output-config -n exercises   # change both values
curl http://localhost:8081/log-output                   # new file content after the sync, old MESSAGE
kubectl rollout restart deployment/log-output-dep -n exercises
curl http://localhost:8081/log-output                   # new MESSAGE too
```

## Configuration

- `PORT` - reader HTTP server port (default: 3000)
- `PING_PONG_URL` - where the reader gets the pong count (default: `http://ping-pong-svc:2345/count`)
- `MESSAGE` - shown as `env variable: MESSAGE=...`, from the ConfigMap
- `INFORMATION_FILE` - file shown as `file content: ...` (default: `/etc/config/information.txt`)

## Files

- `writer/` - Writes the timestamped random string to the shared file (Go 1.27, Alpine 3.24 with tzdata, Europe/Helsinki time)
- `reader/` - Serves the ConfigMap values, the log line and the ping-pong count over HTTP (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - One pod, two containers sharing an `emptyDir` volume; the reader also gets the ConfigMap file and env variable
- `manifests/configmap.yaml` - ConfigMap `log-output-config` with `information.txt` and `MESSAGE`
- `manifests/aks/` - the same for Azure Kubernetes Service, images from ACR
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)

## Technologies

- **Language:** Go 1.27
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes, Traefik Ingress (k3s default)
