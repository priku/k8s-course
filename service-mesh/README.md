# Service mesh

[Istio](https://istio.io/) in ambient mode in the local k3d cluster (exercises 5.2 and 5.3). In ambient mode the pods have no sidecars: `ztunnel` on every node takes care of encrypted traffic and layer 4 policies, and a `waypoint` proxy takes care of layer 7 features such as HTTP policies and traffic splitting, only where it is asked for.

## Install (exercise 5.2)

Istio 1.31 follows the [getting started guide](https://istio.io/latest/docs/ambient/getting-started/), with the fixes below that this cluster needs.

```bash
brew install istioctl
colima ssh -- sudo sh -c 'printf "fs.inotify.max_user_instances=1024\nfs.inotify.max_user_watches=1048576\n" > /etc/sysctl.d/99-inotify.conf && sysctl -q -p /etc/sysctl.d/99-inotify.conf'
istioctl install --set profile=ambient --set values.global.platform=k3d \
  --set values.cni.cniBinDir=/var/lib/rancher/k3s/data/cni --skip-confirmation
```

| Problem | Symptom | Fix |
|---|---|---|
| Colima allows too few inotify instances for three nodes | `istio-cni-node` crashes with `too many open files` | Raise `fs.inotify.max_user_instances` in the VM |
| `global.platform=k3d` installs the CNI plugin to `/bin`, but k3s 1.37 looks for it in `/var/lib/rancher/k3s/data/cni` | **No new pod starts anywhere in the cluster**: `failed to find plugin "istio-cni"` | Set `values.cni.cniBinDir` |

The second one is worth knowing: the plugin is added to the CNI configuration of the whole node, so a wrong path breaks every namespace, not just the mesh. It is also a deadlock if it has already happened, since the corrected `istio-cni-node` cannot start either. Then copy the plugin by hand on each node once: `docker exec <node> cp /bin/istio-cni /var/lib/rancher/k3s/data/cni/istio-cni`.

The cluster keeps Traefik, which the guide turns off. The Gateway of the sample app is a `ClusterIP` and is reached with a port-forward, so the two do not compete for a port.

## The sample app

Bookinfo in the namespace `default`, following the guide until its clean up:

```bash
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml -f samples/bookinfo/platform/kube/bookinfo-versions.yaml
kubectl apply -f samples/bookinfo/gateway-api/bookinfo-gateway.yaml
kubectl annotate gateway bookinfo-gateway networking.istio.io/service-type=ClusterIP
kubectl label namespace default istio.io/dataplane-mode=ambient
kubectl port-forward svc/bookinfo-gateway-istio 8080:80
```

Kiali uses the Prometheus of [`../monitoring`](../monitoring), so `samples/addons/kiali.yaml` is changed as the exercise says, and Istio's own Prometheus addon is not installed:

```yaml
prometheus:
  enabled: true
  url: http://prom-prometheus-server.monitoring:80
```

```bash
kubectl apply -f samples/addons/kiali.yaml
kubectl port-forward svc/kiali -n istio-system 20001:20001   # http://localhost:20001/kiali
```

![Bookinfo in Kiali](images/bookinfo-kiali.png)

| Step of the guide | Tested | Result |
|---|---|---|
| Ambient on the namespace | `istioctl ztunnel-config workloads` | The pods use `HBONE`, without any restarts |
| Layer 4 policy: only the Gateway may call `productpage` | `curl` from another pod | Refused, exit code 56; through the Gateway `200` |
| Waypoint and layer 7 policy: only `GET` from the `curl` pod | `DELETE`, and `GET` from `reviews-v1` | `RBAC: access denied` for both; `GET` from `curl` works |
| Traffic split 90/10 with an HTTPRoute | 100 page loads | 96 to `reviews-v1`, 4 to `reviews-v2` |
