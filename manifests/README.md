# Cluster storage

Cluster-level storage that the [Log output](../log-output) and [ping-pong](../ping-pong) applications shared in exercise 1.11. Since exercise 2.1 they talk over HTTP instead, so nothing uses these for the time being. These are kept apart from the apps, since PersistentVolumes are usually managed by cluster administrators rather than application developers.

- `persistentvolume.yaml` - `shared-pv`: a 1Gi local PersistentVolume at `/tmp/kube` on node `k3d-k3s-default-agent-0`
- `persistentvolumeclaim.yaml` - `shared-pvc`: claims `shared-pv` through the shared storage class name `local-storage`

## Setup (k3d)

A local PersistentVolume uses a directory on the node, so create it first:

```bash
docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/kube
kubectl apply -f manifests/
kubectl get pv,pvc   # shared-pvc should be Bound to shared-pv
```

The volume is tied to `agent-0`, so every pod that uses it is scheduled there. Local volumes are fine for development but not for production: if the node goes away, so does the data.
