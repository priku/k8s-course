# Argo CD

[Argo CD](https://argo-cd.readthedocs.io/) keeps the AKS cluster in the state described in Git (exercises 4.7–4.10). This is a pull deployment: GitHub Actions builds the images and writes their tags to Git, and Argo CD, running inside the cluster, applies the change. CI needs no access to the cluster for this.

```
git push → GitHub Actions: build, push to ACR, write the tag to the branch release → Argo CD syncs the cluster
```

## Install

Argo CD is the one thing installed by hand, since it has to exist before it can manage anything:

```bash
kubectl create namespace argocd
kubectl apply --server-side -n argocd \
  -f https://raw.githubusercontent.com/argoproj/argo-cd/v3.5.6/manifests/install.yaml
```

The material opens the UI with a `LoadBalancer`. Here it has no public address, only a port-forward, so that the admin login is not exposed to the internet:

```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
# https://localhost:8080, user admin
kubectl get secret argocd-initial-admin-secret -n argocd -o jsonpath='{.data.password}' | base64 -d
```

## Applications

| File | App | Follows |
|---|---|---|
| [`log-output.yaml`](log-output.yaml) | Log output and ping-pong in the namespace `exercises` | `log-output/manifests/aks` on the branch `release` |
| [`argocd/project.yaml`](https://github.com/priku/k8s-course-config/tree/main/argocd/project.yaml) in the configuration repository | Production of the todo project in the namespace `project` | `overlays/production` of [priku/k8s-course-config](https://github.com/priku/k8s-course-config) |
| [`argocd/project-staging.yaml`](https://github.com/priku/k8s-course-config/tree/main/argocd/project-staging.yaml) in the configuration repository | Staging of the todo project in the namespace `staging` | `overlays/staging` of [priku/k8s-course-config](https://github.com/priku/k8s-course-config) |

The Applications are applied once, and after that Git decides. Secrets are the exception: they are SOPS-encrypted in the repository and applied separately, since Argo CD has no way to decrypt them here.

```bash
kubectl apply -f ../manifests/exercises-namespace.yaml
sops --decrypt ../ping-pong/secrets/postgres-secret.enc.yaml | kubectl apply -f -
kubectl apply -f log-output.yaml
kubectl get application -n argocd
```

## The branch release

Log output keeps its configuration in this repository. `main` is protected and takes changes only through pull requests, so its workflow cannot commit image tags to it: it merges `main` into the branch `release` and writes the tags there. Never commit to `release` by hand.

The todo project has a repository of its own for its configuration (exercise 4.10), so its workflows commit there instead.
