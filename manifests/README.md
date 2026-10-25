# Cluster resources

Cluster-level resources shared by several applications. These are kept apart from the apps, since they are usually managed by cluster administrators rather than application developers.

- `exercises-namespace.yaml` - Namespace `exercises` for the [Log output](../log-output) and [ping-pong](../ping-pong) applications (exercise 2.3)

Create the namespace before deploying the apps into it:

```bash
kubectl apply -f manifests/exercises-namespace.yaml
```

The shared PersistentVolume of exercise 1.11 was removed in 2.3: since 2.1 the apps talk over HTTP instead.
