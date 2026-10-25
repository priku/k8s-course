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

## Gateway API on AKS (exercise 3.3)

In 3.2 the AKS cluster ran Traefik, installed with Helm, as its Ingress controller. Since 3.3 the apps use the Gateway API instead, with the implementation that AKS manages itself: the **application routing Gateway API implementation**, GatewayClass `approuting-istio`. It is enabled in [`../terraform/gateway.tf`](../terraform/gateway.tf) together with the AKS-managed Gateway API CRDs, and it is the AKS counterpart of the GKE gateway classes in the course material. Traefik is no longer needed:

```bash
helm uninstall traefik --namespace traefik   # if it is still installed from 3.2
kubectl get gatewayclass                     # approuting-istio  ACCEPTED True
kubectl get pods -n aks-istio-system         # istiod: the control plane, managed by AKS
```

For every `Gateway`, AKS creates an Envoy proxy Deployment (two replicas with an autoscaler), a PodDisruptionBudget and a `LoadBalancer` Service with a public IP. Its address is in the Gateway's status (`kubectl get gateway`).
