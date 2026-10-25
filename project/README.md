# The todo project with Kustomize

The todo project is three applications that are always deployed together: the frontend [todo-project](../todo-project), [todo-backend](../todo-backend) with its Postgres, and the hourly CronJob [wiki-todo-generator](../wiki-todo-generator). Since exercise 3.5 they are deployed with [Kustomize](https://kubectl.docs.kubernetes.io/references/kustomize/), which is built into `kubectl` (`kubectl apply -k`).

```
todo-project/manifests/kustomization.yaml          ┐ each app's manifests are a base:
todo-backend/manifests/kustomization.yaml          │ what is the same in every cluster
wiki-todo-generator/manifests/kustomization.yaml   ┘
project/base/   → the three apps together
project/k3d/    → base + Traefik Ingresses (local cluster, images from `k3d image import`)
project/aks/    → base + Gateway API routing, images from the registry acrdwkakssc.azurecr.io
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

## Deployment pipeline (exercise 3.6)

[`.github/workflows/project.yaml`](../.github/workflows/project.yaml) deploys the project to AKS on every push to `main` that changes it:

1. **Sign in to Azure** with `azure/login` and GitHub's OIDC token. Entra ID trusts the token because of a federated credential on the managed identity `id-dwk-github-actions` ([`../terraform/github.tf`](../terraform/github.tf)): only jobs of this repository that run in the GitHub environment `aks` can use the identity. There are no passwords or keys in GitHub, only the IDs of the identity, tenant and subscription.
2. **Start the cluster** if the nightly stop has stopped it.
3. **Build and push** the three images to ACR, tagged `<branch>-<commit SHA>`.
4. **Deploy** with Kustomize: `kustomize edit set image` puts the new tags into `project/aks`, and `kustomize build | kubectl apply` applies it. The database password is decrypted with SOPS using the workflow's own age key (the environment secret `SOPS_AGE_KEY`); `.sops.yaml` encrypts the project's Secret for it and for the personal key, and other Secrets for the personal key only.
5. **Wait** for the rollouts and write the address of the Gateway to the job summary.

| Course material (GCP) | Here (Azure) |
|---|---|
| Service account `github-actions-sa` | User-assigned managed identity `id-dwk-github-actions` |
| Workload Identity Pool and OIDC provider, `roles/iam.workloadIdentityUser` binding | Federated identity credential, subject `repo:priku/k8s-course:environment:aks` |
| `roles/artifactregistry.writer` | `AcrPush` on the registry |
| `roles/container.admin` | `Azure Kubernetes Service Cluster User Role` and `Azure Kubernetes Service RBAC Cluster Admin` on the cluster, `Azure Kubernetes Service Contributor Role` to start it |
| `google-github-actions/auth`, `get-gke-credentials` | `azure/login`, `azure/aks-set-context` with kubelogin |
| Artifact Registry | Azure Container Registry |

The frontend's Deployment uses the `Recreate` strategy: its image cache is a `ReadWriteOnce` Azure Disk that only one node can attach, and a rolling update would start the new pod, possibly on another node, before the old one lets go of the disk.

GitHub environment `aks` (Settings → Environments) holds the secrets `AZURE_CLIENT_ID`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID` and `SOPS_AGE_KEY`.

## Files

- `base/kustomization.yaml` - the three apps
- `k3d/` - the Ingresses for `/` and `/api` and the Traefik Middleware `strip-api`
- `aks/kustomization.yaml` - the images from ACR
- `aks/gateway.yaml` - Gateway and HTTPRoute for `/` and `/api`
