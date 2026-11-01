# Ping-pong Application

A Go HTTP server that responds to `GET /` with `pong N`, where N is the number of earlier requests. `GET /count` returns the current count as JSON (`{"count":N}`), which the [Log output application](../log-output) fetches over HTTP. In the clusters the app is reached at `/pingpong`: the route rewrites that path to `/` (since exercise 3.4). The count is stored in a PostgreSQL database that runs as a StatefulSet, so it survives restarts of both ping-pong and the database.

## Exercises

- **1.9:** Responds with "pong 0", "pong 1", …; shares an Ingress with the [Log output application](../log-output) and is routed at `/pingpong`
- **1.11:** Saves the request count to a PersistentVolume shared with Log output
- **2.1:** `GET /count` endpoint for Log output instead of the shared volume
- **2.2:** The shared Ingress routes Log output at `/log-output`, leaving `/` to the todo project
- **2.3:** Runs in the `exercises` namespace together with Log output
- **2.7:** Saves the counter in Postgres, run as a StatefulSet with one replica
- **3.1:** Runs in Azure Kubernetes Service, exposed with a `LoadBalancer` Service
- **3.2:** In AKS together with Log output behind an Ingress, answering at `/pingpong`; `GET /` answered 200 for Ingress health checks (until 3.4, where `/` became the pong itself)
- **3.3:** The Ingress is replaced with the Gateway API: a Gateway and an HTTPRoute
- **3.4:** The app answers at `/`; the HTTPRoute (and in k3d the Ingress) rewrites `/pingpong` to `/`
- **4.1:** Readiness probe on `/healthz`: the pod is ready only when the database answers
- **4.4:** In k3d the Deployment is replaced with an Argo Rollouts `Rollout`: a canary release, reverted if the CPU usage of the namespace goes over a limit
- **5.7:** A serverless version with [Knative](../knative): `manifests/knative/service.yaml`

## Running locally

```bash
docker network create pingpong
docker run -d --name postgres --network pingpong -e POSTGRES_PASSWORD=local postgres:18.6-alpine3.24
docker build -t ping-pong:v4.1 .
docker run -p 3000:3000 --network pingpong -e PGHOST=postgres -e PGUSER=postgres \
  -e PGPASSWORD=local -e PGSSLMODE=disable ping-pong:v4.1
curl http://localhost:3000/   # pong 0
curl http://localhost:3000/   # pong 1
curl http://localhost:3000/count      # {"count":2}
```

## Deploying to Kubernetes (k3d)

The cluster needs host port 8081 mapped to the load balancer (Traefik Ingress):

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Deploy into the namespace `exercises` together with Log output, which is routed at `/log-output`. The database password is a Secret encrypted with SOPS (see [Database password](#database-password)), so decrypt it into the cluster first:

```bash
kubectl apply -f ../manifests/exercises-namespace.yaml
sops --decrypt secrets/postgres-secret.enc.yaml | kubectl apply -f -
docker build -t ping-pong:v4.1 .
k3d image import ping-pong:v4.1
kubectl apply -f manifests/               # ping-pong, Postgres StatefulSet, Ingress (needs Argo Rollouts, see below)
kubectl apply -f ../log-output/manifests/
curl http://localhost:8081/pingpong   # pong 0
kubectl get all -n exercises
```

The app answers at `/`, so the Ingress removes the `/pingpong` prefix with the Traefik Middleware `strip-pingpong` (in `manifests/ingress.yaml`).

Traffic flow: `localhost:8081/pingpong` → Traefik Ingress `app-ingress` (namespace `exercises`), `/pingpong` stripped → ClusterIP Service `ping-pong-svc:2345` → pod port `3000`. Log output calls `http://ping-pong-svc:2345/count` inside the cluster; `/count` is not routed by the Ingress. An Ingress can only point to Services in its own namespace, so `app-ingress` lives in `exercises` with both apps.

## Deploying to Azure Kubernetes Service (exercises 3.1–3.4)

The course material uses Google Kubernetes Engine; these exercises use AKS instead (see [`../terraform`](../terraform) for the cluster). `manifests/aks/` holds the files that differ from the k3d ones:

- `deployment.yaml` - the image comes from the course registry, `acrdwkakssc.azurecr.io/ping-pong:v4.1`
- `service.yaml` - ClusterIP (port 2345 → 3000). In 3.1 it was `type: LoadBalancer`, which made AKS add a rule and a public IP for ping-pong to the cluster's Azure Load Balancer; since 3.2 the traffic comes in through a shared entry point
- `gateway.yaml` - Gateway `app-gateway` (GatewayClass `approuting-istio`, HTTP on port 80) and HTTPRoute `app-route`: `/pingpong` → ping-pong with a `URLRewrite` filter (`ReplacePrefixMatch: /`), everything else → [Log output](../log-output). In 3.2 this was an Ingress for Traefik
- `statefulset.yaml` - no `storageClassName`, so AKS creates an Azure Disk (StandardSSD) from its default storage class; 1Gi is the smallest disk Azure makes

```bash
# Image for the AKS nodes (amd64), built in Azure by ACR Tasks
az acr build --registry acrdwkakssc --image ping-pong:v4.1 --platform linux/amd64 .

kubectl config use-context aks-dwk-sc
kubectl apply -f ../manifests/exercises-namespace.yaml
sops --decrypt secrets/postgres-secret.enc.yaml | kubectl apply -f -
kubectl apply -f manifests/aks/                      # ping-pong, Postgres, Gateway and HTTPRoute
kubectl apply -f ../log-output/manifests/aks/        # Log output behind the same Gateway
kubectl wait --for=condition=programmed gateway/app-gateway -n exercises
kubectl get gateway app-gateway -n exercises         # ADDRESS: the public IP of the Gateway
curl http://<ADDRESS>/pingpong                       # pong 0
curl http://<ADDRESS>/                               # Log output with Ping / Pongs: 1
```

The Gateway needs the GatewayClass `approuting-istio` in the cluster; see [`../manifests`](../manifests/README.md#gateway-api-on-aks-exercise-33). Since exercise 3.4 the app answers at `/`, and the route rewrites `/pingpong` to `/`, so the app does not need to know the URL structure of the cluster. The app itself answers 404 to `/pingpong`, and to any other path such as `/favicon.ico`, so stray requests do not add to the counter.

The first start can end in one restart: attaching the Azure Disk to Postgres may take longer than the minute ping-pong waits for the database, and Kubernetes then restarts ping-pong.

## Readiness (exercise 4.1)

`GET /healthz` answers `200 ok` only when the database answers, otherwise `503`. The readiness probe calls it every 5 seconds, so without a database the pod runs but is not ready (`0/1`), and the Service sends it no traffic.

For that, ping-pong no longer gives up when the database is missing at startup. It starts serving right away and connects in the background every 5 seconds until the database answers and the `counter` table exists; until then `/`, `/count` and `/healthz` answer `503`. (Before, it waited a minute and then crashed, so a missing database meant a crash loop instead of a not-ready pod.)

Test: apply everything but the database StatefulSet, then add it:

```bash
kubectl delete statefulset postgres-ss -n exercises   # the claim and the data stay
kubectl apply -f manifests/rollout.yaml -f ../log-output/manifests/deployment.yaml
kubectl get pods -n exercises    # ping-pong 0/1, log-output 1/2 (Running, no restarts)
kubectl apply -f manifests/statefulset.yaml
kubectl get pods -n exercises    # a moment later ping-pong 1/1, log-output 2/2
```


## Canary release (exercise 4.4)

In k3d, [`manifests/rollout.yaml`](manifests/rollout.yaml) is an [Argo Rollouts](https://argoproj.github.io/rollouts/) `Rollout` instead of a Deployment. The AKS manifests in `manifests/aks/` are still a plain Deployment.

```bash
kubectl create namespace argo-rollouts
kubectl apply --server-side -n argo-rollouts \
  -f https://github.com/argoproj/argo-rollouts/releases/download/v1.10.1/install.yaml
```

An update first moves one of the two pods to the new version. Then the analysis in [`manifests/analysistemplate.yaml`](manifests/analysistemplate.yaml) asks [Prometheus](../monitoring) every 30 seconds, ten times, for the CPU usage of all containers in the namespace:

```promql
scalar(sum(rate(container_cpu_usage_seconds_total{namespace="exercises", container!=""}[2m])))
```

The limit is `result < 0.1` CPU cores; the namespace normally uses about 0.001. A single measurement over the limit fails the analysis, and the Rollout aborts and goes back to the old version. If all ten pass, the rest of the pods are updated.

Tested in k3d, using two tags of the same image, since only the update itself matters here:

| Test | Measured | Result |
|---|---|---|
| Limit set too low (`result < 0.00001`), `v4.1` → `v4.4` | 0.0005 | Failed at the first measurement, aborted, both pods stayed on `v4.1` |
| Limit `0.1`, `v4.1` → `v4.4` | 0.0009–0.0015, ten times | Passed after 4 min 44 s, both pods on `v4.4` |
| Limit `0.1`, a pod burning CPU in the namespace, `v4.4` → `v4.1` | 1.0 | Failed at the first measurement, aborted, both pods stayed on `v4.4` |

To try the first one:

```bash
kubectl patch analysistemplate cpu-usage -n exercises --type=json \
  -p '[{"op":"replace","path":"/spec/metrics/0/successCondition","value":"result < 0.00001"}]'
docker tag ping-pong:v4.1 ping-pong:v4.4 && k3d image import ping-pong:v4.4
kubectl patch rollout ping-pong-dep -n exercises --type=json \
  -p '[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"ping-pong:v4.4"}]'
kubectl get rollout,analysisrun -n exercises --watch
kubectl get rollout ping-pong-dep -n exercises -o jsonpath='{.status.message}'   # RolloutAborted: ...
```

An aborted Rollout stays `Degraded` until the image changes or the update is tried again. The [kubectl plugin](https://argoproj.github.io/argo-rollouts/installation/#kubectl-plugin-installation) does that with `kubectl argo rollouts retry rollout ping-pong-dep`, and without the plugin:

```bash
kubectl patch rollout ping-pong-dep -n exercises --subresource=status --type=merge -p '{"status":{"abort":false}}'
```

## Database

`manifests/statefulset.yaml` holds the database in two resources:

- StatefulSet `postgres-ss` with one replica (`postgres:18.6-alpine3.24`). Its pod is always called `postgres-ss-0`, and the `volumeClaimTemplates` give it its own PersistentVolumeClaim `data-postgres-ss-0` from the k3s `local-path` storage class, so no PersistentVolume has to be created by hand. Postgres 18 keeps its data in `/var/lib/postgresql/18/docker`, so the volume is mounted at `/var/lib/postgresql`.
- Headless Service `postgres-svc` (`clusterIP: None`), which gives the pod a stable DNS name. It only resolves once the pod is ready, so ping-pong retries the connection for a minute at startup.

The counter is a single row in the table `counter`, and `GET /` adds one with an atomic `UPDATE ... RETURNING`.

The data outlives the pods and even the StatefulSet: deleting it leaves the claim, which binds again when the StatefulSet is applied.

```bash
kubectl delete pod postgres-ss-0 -n exercises       # recreated with the same name and volume
kubectl delete statefulset postgres-ss -n exercises
kubectl get pvc -n exercises                        # data-postgres-ss-0 is still there
kubectl apply -f manifests/statefulset.yaml         # the counter continues where it was
```

To look inside the database, start a stand-alone psql pod with the password from the Secret:

```bash
kubectl run psql-debug -n exercises --rm -it --restart=Never --image postgres:18.6-alpine3.24 \
  --overrides='{"spec":{"containers":[{"name":"psql-debug","image":"postgres:18.6-alpine3.24","stdin":true,"tty":true,"command":["psql"],"env":[{"name":"PGHOST","value":"postgres-svc"},{"name":"PGUSER","value":"postgres"},{"name":"PGPASSWORD","valueFrom":{"secretKeyRef":{"name":"postgres-secret","key":"POSTGRES_PASSWORD"}}}]}]}}'
postgres=# SELECT * FROM counter;
```

## Database password

The password is in the Secret `postgres-secret` (key `POSTGRES_PASSWORD`), which both Postgres and ping-pong read. `secrets/postgres-secret.enc.yaml` is encrypted with [SOPS](https://github.com/getsops/sops) and [age](https://github.com/FiloSottile/age), so it is safe in version control: only the `stringData` values are encrypted, and the age public key and rules are in [`../.sops.yaml`](../.sops.yaml). The file is kept out of `manifests/`, so that `kubectl apply -f manifests/` does not try to apply the encrypted version.

Decrypting needs the matching private key, which is not in the repository. SOPS finds it in `~/Library/Application Support/sops/age/keys.txt` on macOS (or the file in `SOPS_AGE_KEY_FILE`).

Without that key, make your own:

```bash
age-keygen -o key.txt                       # prints the public key
export SOPS_AGE_KEY_FILE=$PWD/key.txt
# Put the public key in ../.sops.yaml, then write and encrypt a new Secret:
printf 'apiVersion: v1\nkind: Secret\nmetadata:\n  namespace: exercises\n  name: postgres-secret\nstringData:\n  POSTGRES_PASSWORD: "%s"\n' "$(openssl rand -hex 24)" \
  | sops --encrypt --filename-override secrets/postgres-secret.enc.yaml --input-type yaml --output-type yaml /dev/stdin \
  > secrets/postgres-secret.enc.yaml
```

Postgres only reads `POSTGRES_PASSWORD` when it creates an empty database, so after changing the password delete the claim `data-postgres-ss-0` too (this deletes the data).

## Configuration

- `PORT` - HTTP server port (default: 3000)
- `PGHOST`, `PGPORT`, `PGUSER`, `PGDATABASE`, `PGSSLMODE` - database connection, the standard Postgres client variables that [pgx](https://github.com/jackc/pgx) reads; set in `manifests/deployment.yaml`
- `PGPASSWORD` - from the Secret `postgres-secret`

## Files

- `main.go` - HTTP server with the counter in Postgres (`/`, `/count`, `/healthz`)
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/rollout.yaml` - Argo Rollouts `Rollout` with `PORT=3000` and the database connection
- `manifests/analysistemplate.yaml` - The CPU analysis of the canary release
- `manifests/statefulset.yaml` - Postgres StatefulSet and its headless Service `postgres-svc`
- `manifests/aks/` - AKS versions of the Deployment (image from ACR), the Service and the StatefulSet (Azure Disk), and the Gateway with its HTTPRoute
- `secrets/postgres-secret.enc.yaml` - SOPS-encrypted Secret with the database password
- `manifests/service.yaml` - ClusterIP Service (port 2345 → 3000)
- `manifests/ingress.yaml` - Shared Ingress: `/pingpong` → ping-pong (prefix stripped by the Middleware `strip-pingpong`), `/log-output` → Log output

## Technologies

- **Language:** Go 1.27, [pgx](https://github.com/jackc/pgx) v5
- **Database:** PostgreSQL 18
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes (k3d locally, Azure Kubernetes Service in the cloud), Traefik Ingress (k3s default)
