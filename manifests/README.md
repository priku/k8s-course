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
