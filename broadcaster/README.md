# broadcaster

Sends a status message of the [todo project](../project) to an external chat service whenever a todo is saved or updated (exercise 4.6).

```
todo-backend → NATS, subject "todos" → broadcaster (6 replicas, queue group) → webhook
```

## How it works

- [todo-backend](../todo-backend) publishes `{"action": "created", "todo": {...}}` or `"updated"` to the subject `todos`. Todos from [wiki-todo-generator](../wiki-todo-generator) go through the backend, so they are included.
- broadcaster subscribes with the queue group `broadcasters`. NATS gives each message to only one member of the group, so the number of replicas does not change how many times a message is sent.
- A message is forwarded as a `POST` with the exercise's "Generic" payload to the URL in `WEBHOOK_URL`:

```json
{"user": "bot", "message": "A todo was created: Buy milk"}
```

- A failed send is logged and not retried. The exercise accepts a missing message when something is broken, but not a duplicate.

The exercise's external service here is [`mendhak/http-https-echo`](https://github.com/mendhak/docker-http-https-echo) in [`manifests/webhook-receiver.yaml`](manifests/webhook-receiver.yaml), which logs what it gets. It needs no accounts, and the project's messages stay inside the cluster. To send them elsewhere, change `WEBHOOK_URL` in [`manifests/configmap.yaml`](manifests/configmap.yaml) to a URL that takes this payload.

NATS itself is in [`../project/base/nats.yaml`](../project/base/nats.yaml), a plain Deployment instead of the course's Helm chart. Every environment gets its own that way, so its messages cannot reach the broadcaster of another environment.

## Trying it

```bash
kubectl get pods -n project -l app=broadcaster     # 6 replicas
curl -X POST -H 'Content-Type: application/json' -d '{"text":"Buy milk"}' http://localhost:8081/api/todos
kubectl logs -n project deployment/webhook-receiver-dep | grep '"body"'
```

## Test results (6 replicas)

| Test | Result |
|---|---|
| k3d: 40 new todos and 40 updates | 80 messages sent, each exactly once, by all six replicas (9–16 messages each) |
| AKS: 15 new todos and 15 updates | 30 messages sent, each exactly once, by all six replicas (1–7 messages each) |
| k3d: NATS stopped | Pods turn unready without restarts; todos are still saved |
| k3d: NATS back | All replicas ready without restarts; the todo saved during the outage was sent once, after the connection returned |

## Configuration

| Variable | Value | Where |
|---|---|---|
| `PORT` | `3000` | Deployment |
| `NATS_URL` | `nats://nats-svc:4222` | ConfigMap `broadcaster-config` |
| `BROADCAST_MODE` | `webhook` sends the messages, `log` only logs them (staging, exercise 4.9) | ConfigMap `broadcaster-config` |
| `WEBHOOK_URL` | `http://webhook-receiver-svc:2345/` (not needed in `log` mode) | ConfigMap `broadcaster-config` |

## Health checks

- `GET /healthz` (readiness): `200` when connected to NATS, otherwise `503`
- `GET /livez` (liveness): `200` while the process runs, so a NATS outage does not restart the pod

## Files

- `main.go` - NATS subscription, message format and the webhook call
- `main_test.go` - Tests for the message format and the webhook call
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24)
- `manifests/deployment.yaml` - Deployment with 6 replicas
- `manifests/configmap.yaml` - The NATS and webhook addresses
- `manifests/webhook-receiver.yaml` - The stand-in for the external chat service

## Technologies

Go 1.27, [nats.go](https://github.com/nats-io/nats.go) 1.55, NATS 2.15
