# The todo project with Kustomize

The todo project is four applications that are always deployed together: the frontend [todo-project](../todo-project), [todo-backend](../todo-backend) with its Postgres, the hourly CronJob [wiki-todo-generator](../wiki-todo-generator), and [broadcaster](../broadcaster) with NATS (exercise 4.6). Since exercise 3.5 they are deployed with [Kustomize](https://kubectl.docs.kubernetes.io/references/kustomize/), which is built into `kubectl` (`kubectl apply -k`).

```
todo-project/manifests/kustomization.yaml          ┐ each app's manifests are a base:
todo-backend/manifests/kustomization.yaml          │ what is the same in every cluster
wiki-todo-generator/manifests/kustomization.yaml   │
broadcaster/manifests/kustomization.yaml           ┘
project/base/   → the three apps together
project/k3d/    → base + Traefik Ingresses (local cluster, images from `k3d image import`)
project/aks/    → production: base + Gateway API routing and backups, images from the registry acrdwkakssc.azurecr.io
project/staging/ → staging: base + Gateway API routing, images from the registry
```

Only what differs between the clusters is in the overlays:

| | k3d (`project/k3d`) | AKS (`project/aks`) |
|---|---|---|
| Images | `todo-project:v2.9` etc., imported into the cluster | `acrdwkakssc.azurecr.io/...`, set with the `images` field |
| Routing | Traefik Ingresses; a Middleware strips `/api` | Gateway `project-gateway` (`approuting-istio`) and HTTPRoute `project-route`; a `URLRewrite` filter strips `/api` |
| Storage | default storage class `local-path` | default storage class, an Azure Disk |

The bases leave out `storageClassName`, so every cluster uses its own default storage class.

## Deploy

The namespace and the database password are applied first: the password is a SOPS-encrypted Secret, which Kustomize cannot decrypt by itself.

```bash
kubectl apply -f manifests/project-namespace.yaml
sops --decrypt todo-backend/secrets/postgres-secret.enc.yaml | kubectl apply -f -

# Local k3d cluster
docker build -t todo-project:v2.9 todo-project
docker build -t todo-backend:v2.10 todo-backend
docker build -t wiki-todo-generator:v2.9 wiki-todo-generator
k3d image import todo-project:v2.9 todo-backend:v2.10 wiki-todo-generator:v2.9
kubectl apply -k project/k3d          # http://localhost:8081

# Azure Kubernetes Service (exercise 3.5)
az acr build --registry acrdwkakssc --image todo-project:v2.9 --platform linux/amd64 todo-project
az acr build --registry acrdwkakssc --image todo-backend:v2.10 --platform linux/amd64 todo-backend
az acr build --registry acrdwkakssc --image wiki-todo-generator:v2.9 --platform linux/amd64 wiki-todo-generator
kubectl config use-context aks-dwk-sc
kubectl apply -k project/aks
kubectl get gateway project-gateway -n project   # ADDRESS: the public IP of the project
```

`kubectl kustomize project/aks` prints the resulting manifests without applying them.

Changing the storage class of an existing claim is not allowed, so a cluster that still has the claims from before exercise 3.5 (with `storageClassName: local-path`) needs them deleted first, together with the Postgres StatefulSet. This deletes the data.

## GitOps and environments (exercises 4.8 and 4.9)

Argo CD deploys the project to AKS. Nothing outside the cluster has access to it. There are two environments:

| | Staging | Production |
|---|---|---|
| Goes there | Every commit to `main` that changes the project | Every tagged commit |
| Workflow | [`project.yaml`](../.github/workflows/project.yaml) | [`project-production.yaml`](../.github/workflows/project-production.yaml) |
| Namespace | `staging` | `project` |
| Overlay | [`staging/`](staging) | [`aks/`](aks) |
| Branch that Argo CD follows | `staging` | `production` |
| Images | `<app>:main-<sha>` | `<app>:production-<sha>` |
| [broadcaster](../broadcaster) | Only logs the messages, one replica | Sends them to the webhook, six replicas |
| Database backups | No | Yes, every day |

1. The workflow signs in to Azure with GitHub's OIDC token, builds the images and pushes them to ACR.
2. It merges the commit into the branch of the environment and writes the new tags to the overlay there.
3. [Argo CD](../argocd) notices the change within three minutes and applies the overlay from that branch.

Production has a branch of its own, so a commit to `main` changes nothing in production until it is tagged. The Applications have `prune` and `selfHeal`, so a resource removed from Git is removed from the cluster, and a manual change in the cluster is reverted. The database password is the exception: it is a SOPS-encrypted Secret, applied to each namespace outside Argo CD. When the nightly stop has stopped the cluster nothing starts it any more, and Argo CD catches up after `az aks start`.

The branch environments of exercises 3.7 and 3.8 are gone: they would need the workflow to have access to the cluster. Their implementation is in the tags [3.7](https://github.com/priku/k8s-course/tree/3.7) and [3.8](https://github.com/priku/k8s-course/tree/3.8).

```bash
kubectl get application -n argocd
kubectl get gateway project-gateway -n staging   # ADDRESS: the public IP of staging
kubectl logs -n staging deployment/broadcaster-dep | grep 'message logged'
```

| Course material (GCP) | Here (Azure) |
|---|---|
| Workload Identity Pool and OIDC provider | Federated identity credential on the managed identity `id-dwk-github-actions`, subject `repo:priku/k8s-course:environment:aks` ([`../terraform/github.tf`](../terraform/github.tf)) |
| `roles/artifactregistry.writer` | `AcrPush` on the registry |
| Artifact Registry | Azure Container Registry |

The frontend's Deployment uses the `Recreate` strategy: its image cache is a `ReadWriteOnce` Azure Disk that only one node can attach, and a rolling update would start the new pod, possibly on another node, before the old one lets go of the disk.

Staging images pile up, a set for every push to `main`. The ACR task `purge-old-project-images` ([`../terraform/acr-purge.tf`](../terraform/acr-purge.tf)) runs `acr purge` every Monday: it deletes `main-*` images older than a week, keeping the five newest of each repository. `production-*` images are kept, so that the purge can never delete what production runs.

GitHub environment `aks` (Settings → Environments) holds the secrets `AZURE_CLIENT_ID`, `AZURE_TENANT_ID` and `AZURE_SUBSCRIPTION_ID`.

## Health checks (exercise 4.2)

| | Readiness probe (gets traffic?) | Liveness probe (restart?) |
|---|---|---|
| todo-backend | `GET /healthz`: the database is connected and answers | `GET /livez`: the process answers |
| todo-project | `GET /readyz`: not broken, and todo-backend's `/healthz` answers 200 | `GET /healthz`: not broken |

The liveness probes leave the database out on purpose: a database outage should make the apps unready, so the Service stops sending them traffic, but restarting them would not help. Both apps now start without a database and connect in the background, so without one they run as `0/1` instead of crashing in a loop.

The page has a **Break the app** button (`POST /break`). After it the frontend answers every request with `500`, and its liveness probe fails; after three failures in a row (about 15 seconds) Kubernetes restarts the container, and the new process works again. Tested in k3d: broken, restarted (`RESTARTS 1`) and working again after 32 seconds. Without the database both apps turned `0/1` without restarts and came back to `1/1` when it returned.

```bash
kubectl get pods -n project -w                                         # RESTARTS goes up after the button
kubectl get events -n project --field-selector reason=Unhealthy        # the failing probes
```

The frontend's in-cluster address of the backend for the readiness probe is `BACKEND_HEALTH_URL` in the ConfigMap `todo-project-config`.

## Logs in Azure Monitor (exercise 3.12)

GKE sends container logs to Cloud Logging by itself. The AKS counterpart is **Azure Monitor Container insights**, enabled in [`../terraform/monitoring.tf`](../terraform/monitoring.tf): the agent `ama-logs` on every node sends the stdout and stderr of the containers in the namespace `project` to the Log Analytics workspace `log-dwk-aks-sc`, table `ContainerLogV2`. todo-backend logs JSON (exercise 2.10), so the fields of each line are columns in `LogMessage`.

To find the log of a new todo:

1. Create a todo in the app.
2. In the Azure portal open **Kubernetes services → aks-dwk-sc → Monitoring → Logs** (or the workspace **log-dwk-aks-sc → Logs**), switch the editor to **KQL mode** and run:

```kusto
ContainerLogV2
| where PodNamespace == "project" and ContainerName == "todo-backend"
| where LogMessage has "todo created"
| project TimeGenerated, PodName, Message = tostring(LogMessage.msg), Todo = tostring(LogMessage.text)
| order by TimeGenerated desc
```

The same from the command line:

```bash
az monitor log-analytics query --workspace "$(az monitor log-analytics workspace show -g rg-dwk-aks-sc -n log-dwk-aks-sc --query customerId -o tsv)" \
  --analytics-query 'ContainerLogV2 | where PodNamespace == "project" and LogMessage has "todo created" | take 10' -o table
```

Logs arrive within a few minutes. To keep the cost down, only the namespace `project` is collected (`namespaceFilteringMode: Include`), only three tables (ContainerLogV2, KubePodInventory, KubeEvents), retention is the free 30 days, and ingestion stops at 0.2 GB a day.

![A new todo in the logs of todo-backend in Azure Monitor](images/azure-monitor-todo-created.png)

## Resource requests and limits (exercise 3.11)

Every container of the project has resource requests and limits. The values come from `kubectl top pods --containers` while `hey` loaded the AKS deployment through the Gateway (about 50 concurrent users on `/` and `/api/todos`, roughly 1 300 requests per second):

| Container | Peak under load | Requests | Limits |
|---|---|---|---|
| todo-project | 162m CPU, 7Mi | 20m, 16Mi | 250m, 64Mi |
| todo-backend | 109m CPU, 7Mi | 20m, 16Mi | 250m, 64Mi |
| postgres | 62m CPU, 37Mi | 50m, 64Mi | 500m, 512Mi |
| wiki-todo-generator (CronJob) | one short run | 10m, 16Mi | 100m, 32Mi |
| db-backup (CronJob) | one short run | 20m, 32Mi | 250m, 128Mi |

- **Requests** are what the scheduler reserves on a node, so they follow normal use: small, so that the project and its branch environments still fit next to the cluster's own pods (the AKS-managed `istiod` alone requests 2 × 500m CPU and 2 × 2Gi).
- **Limits** leave room above the measured peaks. Going over the CPU limit only slows a container down (throttling), but going over the memory limit kills it (`OOMKilled`), so the memory limits have the larger margin. Postgres keeps up to 128MB of shared buffers by default, hence its 512Mi.

```bash
kubectl top pods -n project --containers
kubectl describe nodes | grep -A6 'Allocated resources'   # what the requests take of each node
```

The Gateway's Envoy proxies are managed by AKS and come with their own requests (100m, 128Mi). A `ResourceQuota` per namespace could cap what a branch environment may take in total; it is not used here because every pod in the namespace, including the managed Gateway pods, would then need limits.

## Database backups (exercise 3.10)

Production (`main`, namespace `project`) dumps its database every day with `pg_dump` and uploads it to Azure Blob Storage, authenticated with workload identity: no keys or Secrets for Azure. The CronJob and its ServiceAccount are the Kustomize component `components/db-backup`, which the deployment workflow adds for `main` only, so branch environments do not back up. See [`../db-backup`](../db-backup) for how it works and how to restore.

## DBaaS vs DIY (exercise 3.9)

The project's Postgres runs in the cluster as a StatefulSet with an Azure Disk ("do it yourself"). The managed alternative in Azure is **Azure Database for PostgreSQL – Flexible Server** ("database as a service"; Cloud SQL in the course material). Prices are Sweden Central list prices in October 2026.

| | DBaaS: Azure Database for PostgreSQL Flexible Server | DIY: Postgres StatefulSet in AKS (current) |
|---|---|---|
| **Getting started** | A Terraform resource (or a few clicks): server, version, size, firewall or private network. The app needs the server's host name and a login; with Entra ID authentication even the password can go, using workload identity. | A StatefulSet, a headless Service and a Secret that we write and keep up to date ourselves (`todo-backend/manifests/statefulset.yaml`). Needs Kubernetes knowledge: volume claim templates, storage classes, the Postgres 18 data directory change. |
| **Cost** | Paid separately from the cluster. The smallest server (Burstable B1ms, 1 vCore, 2 GiB) is about 0.0175 €/h, so about 12.8 €/month if it runs all the time, plus storage from 32 GB upwards at 0.12 €/GB/month (about 3.9 €/month). It can be stopped for up to 7 days at a time. | Runs on nodes we pay for anyway, so the extra cost is the disk: 1 GiB StandardSSD (E1) is about 0.26 €/month. It also takes some CPU and memory from the nodes, which here means fewer branch environments fit. |
| **Maintenance** | Azure patches the OS and minor versions in a maintenance window that we can choose; major upgrades are an in-place operation. Monitoring, metrics and alerts are built in. | All ours: image updates (we bumped 13 → 18 by hand), tuning, monitoring, and making sure the volume follows the pod. A wrong mount path (as in the original 2.7) silently loses data. |
| **Backups** | Automatic daily backups with point-in-time restore to any second within the retention period (7–35 days), included up to the size of the provisioned storage. Restoring is one command or a few clicks and creates a new server. Geo-redundant backups as an option. | Nothing by default. We have to build it: a CronJob with `pg_dump` and storage for the dumps (exercise 3.10), retention, and a tested restore procedure. Only as fresh as the last dump; no point-in-time restore without WAL archiving, which is more work again. |
| **Availability and scaling** | Zone-redundant high availability with automatic failover as an option (doubles the compute cost). Size and storage can be changed with little downtime; read replicas are available. | One pod and one disk in one zone. If the node fails, the pod and the disk move to another node, with downtime meanwhile. Scaling Postgres in Kubernetes properly needs an operator (e.g. CloudNativePG). |
| **Security** | TLS, encryption at rest, Entra ID authentication, private endpoints and Defender for Cloud are built in. | Our job: the Secret (SOPS here), network policies, TLS between the app and the database (currently `sslmode=disable` inside the cluster). |
| **Portability** | Tied to Azure, although it is standard Postgres and `pg_dump` works for moving out. | The same manifests run in k3d, AKS or any other cluster. |
| **Branch environments** | One server per branch would be slow and expensive to create; branches would rather share one server with a database each. | A new StatefulSet per branch is cheap and quick, which is what makes the per-branch environments of exercise 3.7 easy. |

**Conclusion.** For this course project DIY is the better fit: it costs almost nothing extra, it works the same locally and in every branch environment, and the data is not valuable. For anything with real users I would choose the managed Flexible Server: automatic backups with point-in-time restore, patching and high availability are a lot of work to build and keep working on our own, and its cost is small compared to that work. A common middle ground is to run the apps in Kubernetes and the production database as a managed service, while development and branch environments keep a database in the cluster.

## Files

- `base/kustomization.yaml` - the apps and NATS
- `k3d/` - the Ingresses for `/` and `/api` and the Traefik Middleware `strip-api`
- `aks/kustomization.yaml` - production: the images from ACR, Gateway and database backups
- `staging/kustomization.yaml` - staging: the images from ACR and Gateway; broadcaster only logs
- `components/gateway/` - Gateway and HTTPRoute for `/` and `/api`
- `components/db-backup/` - the daily backup CronJob and its workload identity ServiceAccount, production only
