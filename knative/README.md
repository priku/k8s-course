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

## Ping-pong as a serverless service (exercise 5.7)

[`../ping-pong/manifests/knative/service.yaml`](../ping-pong/manifests/knative/service.yaml) makes [ping-pong](../ping-pong) a Knative Service, and Log output calls it at `http://ping-pong.exercises.svc.cluster.local/count`. The counter is in Postgres, so nothing is lost when the pods go away.

```bash
docker tag ping-pong:v4.1 dev.local/ping-pong:v5.7
docker build -t log-reader:v5.7 ../log-output/reader
k3d image import -c knative dev.local/ping-pong:v5.7 log-reader:v5.7 log-writer:v1.11

kubectl apply -f ../manifests/exercises-namespace.yaml
sops --decrypt ../ping-pong/secrets/postgres-secret.enc.yaml | kubectl apply -f -
kubectl apply -f ../ping-pong/manifests/statefulset.yaml
kubectl apply -f ../ping-pong/manifests/knative/service.yaml
kubectl apply -f ../log-output/manifests/configmap.yaml -f ../log-output/manifests/service.yaml \
  -f ../log-output/manifests/knative/deployment.yaml

curl -H "Host: ping-pong.exercises.172.19.0.3.sslip.io" http://localhost:8081/     # pong 0
kubectl port-forward svc/log-output-svc -n exercises 8090:2345                      # http://localhost:8090
```

This cluster has no Traefik, so Log output is reached with a port-forward.

Three things had to change for this to work:

| Problem | Why | Fix |
|---|---|---|
| Knative cannot find a local image | It looks up the digest of a tag from the registry | The prefix `dev.local/`, whose tags Knative does not look up |
| Ping-pong would never scale to zero | The readiness probe of Log output (exercise 4.1) calls it every 5 seconds, which counts as traffic | [The Knative version](../log-output/manifests/knative/deployment.yaml) of Log output has no such probe |
| Log output gave `502` when ping-pong was starting | Its limit of 2 seconds for a request was too short | 45 seconds |

| Test, all through Log output | Result |
|---|---|
| No traffic | No ping-pong pods after about a minute |
| The request after all the pods are gone | `200` in 1.6 seconds, and the count continues from where it was |
| A request while Knative is still shutting down the previous pod | `200` in 29 seconds: it waits for the drain |

The last one is the cost of scaling to zero that shows in practice: for about half a minute after a scale-down, the next request is slow.
