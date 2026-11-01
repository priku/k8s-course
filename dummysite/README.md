# DummySite

A custom Kubernetes resource and its controller (exercise 5.1). A `DummySite` names a page, and the controller publishes a copy of it.

```yaml
apiVersion: stable.dwk/v1
kind: DummySite
metadata:
  namespace: default
  name: example
spec:
  website_url: https://example.com/
```

## How it works

The controller watches DummySites through the Kubernetes API with [client-go](https://github.com/kubernetes/client-go). For each one it creates:

| Resource | Name | Does |
|---|---|---|
| Deployment | `<name>-dep` | An init container fetches the page with `curl` into a shared volume, and nginx serves it |
| Service | `<name>-svc` | ClusterIP in front of nginx |
| Ingress | `<name>-ingress` | Publishes the site at `<name>.localhost` |

- **Reconciling:** server-side apply creates a resource or updates it to the wanted state, so reconciling again is safe. A failed reconcile is retried with a growing delay. After a restart the controller goes through the existing DummySites again.
- **A changed URL** changes the Pod template, so Kubernetes starts a new Pod that fetches the new page.
- **Deleting:** the resources have the DummySite as their owner, so Kubernetes deletes them with it. The controller does not need to delete anything.
- **Validation:** the CRD accepts only `http` and `https` URLs, and the controller checks them again.

The copy is one HTML file, so relative links and styles of a complex page may break, which the exercise allows.

## Trying it

```bash
docker build -t dummysite-controller:v5.1 .
k3d image import dummysite-controller:v5.1

kubectl apply -f manifests/customresourcedefinition.yaml
kubectl apply -f manifests/clusterrole.yaml -f manifests/serviceaccount.yaml -f manifests/clusterrolebinding.yaml
kubectl apply -f manifests/deployment.yaml
kubectl apply -f dummysite.yaml

kubectl get dummysites
curl -H 'Host: example.localhost' http://localhost:8081/   # or http://example.localhost:8081 in a browser
kubectl delete dummysite example                           # removes the Deployment, Service and Ingress
```

## Access

The ClusterRole has only what the controller uses: read DummySites, and create and patch Deployments, Services and Ingresses. Server-side apply needs `create` for a resource that does not exist yet and `patch` for the rest. It is a ClusterRole with a ClusterRoleBinding since a DummySite can be in any namespace.

## Test results (k3d)

| Test | Result |
|---|---|
| The workflow of the exercise from an empty cluster: role, account, binding, Deployment, DummySite | The site was up in seconds |
| `https://example.com/` | The copy is byte for byte the same as the original |
| `https://en.wikipedia.org/wiki/Kubernetes` | About 595 kB copied, title correct |
| A changed `website_url` | A new Pod fetched the new page |
| Deleting the DummySite | Its Deployment, Service and Ingress disappeared |
| `file:///etc/passwd` | Refused by the API server |
| Restarting the controller | The existing DummySite was picked up again |

## Files

- `main.go` - The controller
- `main_test.go` - Tests for the validation and reconciling, against a fake API
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `dummysite.yaml` - An example DummySite
- `manifests/customresourcedefinition.yaml` - The resource type `DummySite`
- `manifests/clusterrole.yaml`, `serviceaccount.yaml`, `clusterrolebinding.yaml` - The access of the controller
- `manifests/deployment.yaml` - The controller

## Technologies

Go 1.27, client-go 0.37, nginx 1.31, curl 8.23
