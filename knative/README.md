# Knative

[Knative Serving](https://knative.dev/) on k3d (exercises 5.6 and 5.7): a serverless layer on top of Kubernetes. A Knative `Service` needs only an image: Knative creates the Deployment, the Services and the route, keeps every version as a revision, and runs no pods at all when there is no traffic.

## Install (exercise 5.6)

Knative needs a cluster without Traefik, so it has a cluster of its own. The material uses Kubernetes 1.34, but Knative 1.24 works with the current version too. It needs at least 1.35. Two clusters do not fit in the memory of the Colima VM, so the other one is stopped, which keeps its state.

```bash
k3d cluster stop k3s-default
k3d cluster create knative --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2 \
  --k3s-arg "--disable=traefik@server:0" --image rancher/k3s:v1.37.1-k3s1

V=knative-v1.24.0
kubectl apply --server-side -f https://github.com/knative/serving/releases/download/$V/serving-crds.yaml
kubectl apply --server-side -f https://github.com/knative/serving/releases/download/$V/serving-core.yaml
kubectl apply --server-side -f https://github.com/knative-extensions/net-kourier/releases/download/$V/kourier.yaml
kubectl patch configmap/config-network -n knative-serving --type merge \
  -p '{"data":{"ingress-class":"kourier.ingress.networking.knative.dev"}}'
kubectl apply --server-side -f https://github.com/knative/serving/releases/download/$V/serving-default-domain.yaml
kubectl get pods -n knative-serving
```

The crash loop that the exercise warns about did not appear. The inotify limit of the VM had already been raised for [Istio](../service-mesh#install-exercise-52), which is the likely reason, though that was not tested separately.

## The examples

Reach a service through Kourier with the host name from `kubectl get ksvc`:

```bash
kubectl apply -f hello.yaml
curl -H "Host: hello.default.172.19.0.3.sslip.io" http://localhost:8081    # Hello World!
kubectl apply -f hello-split.yaml                                           # a new revision, 50/50
```

| Example | Result |
|---|---|
| [`hello.yaml`](hello.yaml): Deploying a Knative Service | `Hello World!` |
| Autoscaling: no traffic | No pods running after 61 seconds |
| Autoscaling: the next request | `200` in 6.6 seconds, the following one in 0.008 seconds |
| [`hello-split.yaml`](hello-split.yaml): Traffic splitting | Of 200 requests, 108 to the new revision and 92 to the old one |
