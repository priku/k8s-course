# Cluster resources

Cluster-level resources shared by several applications. These are kept apart from the apps, since they are usually managed by cluster administrators rather than application developers.

- `exercises-namespace.yaml` - Namespace `exercises` for the [Log output](../log-output) and [ping-pong](../ping-pong) applications (exercise 2.3)
- `project-namespace.yaml` - Namespace `project` for the [todo project](../todo-project) and [todo-backend](../todo-backend) (exercise 2.4)

Create the namespaces before deploying the apps into them:

```bash
kubectl apply -f manifests/
kubectl get namespaces
```

The shared PersistentVolume of exercise 1.11 was removed in 2.3: since 2.1 the apps talk over HTTP instead.

## Ingress controller on AKS (exercise 3.2)

k3s comes with Traefik, but an AKS cluster has no Ingress controller of its own. The application routing add-on's managed NGINX is based on the retired ingress-nginx project, and Azure supports it only until November 2026, so the AKS cluster runs Traefik too, installed with Helm. The same Ingress manifests then work in k3d and in AKS.

```bash
helm repo add traefik https://traefik.github.io/charts
helm upgrade --install traefik traefik/traefik --version 42.1.3 \
  --namespace traefik --create-namespace --values manifests/aks/traefik-values.yaml
kubectl get svc traefik -n traefik   # EXTERNAL-IP: the public IP for all Ingresses
```

`aks/traefik-values.yaml` makes Traefik the default IngressClass `traefik` and exposes it with a `LoadBalancer` Service, so all Ingress traffic goes through one Azure public IP.
