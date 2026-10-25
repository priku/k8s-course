# Wiki Todo Generator

A Kubernetes CronJob that adds a reading todo to the [todo project](../todo-project) every hour: `Read <URL>`, where the URL is a random Wikipedia article.

## Exercises

- **2.9:** CronJob that creates a `Read <URL>` todo with a random Wikipedia article every hour

## How it works

1. `GET https://en.wikipedia.org/wiki/Special:Random` without following the redirect. Wikipedia answers `302 Found` with the random article in the `Location` header, as a protocol-relative address (`//en.wikipedia.org/wiki/Furfjord`) that is resolved to `https://`.
2. The title is decoded for reading (`Äänekoski` instead of `%C3%84%C3%A4nekoski`), unless decoding would bring back `?`, `#`, `%` or a space, which would change the URL when pasted.
3. If `Read <URL>` is longer than `TODO_MAX_LENGTH`, the job tries another article (up to five), since todo-backend would refuse it.
4. `POST http://todo-backend-svc:2345/todos` with `{"text": "Read <URL>"}`. The job runs inside the cluster, so it calls the backend's Service directly, not through the Ingress.

The run ends when the todo is created; a failure exits with an error, and Kubernetes retries the run.

## CronJob

`manifests/cronjob.yaml` in the namespace `project`:

- `schedule: "0 * * * *"` - every hour at minute 0
- `concurrencyPolicy: Forbid` - a new run is skipped while the previous one is still going
- `backoffLimit: 3` with `restartPolicy: Never` - a failed run is retried up to three times in new pods
- The last three successful and three failed Jobs are kept, so their logs can still be read

```bash
docker build -t wiki-todo-generator:v2.9 .
k3d image import wiki-todo-generator:v2.9
kubectl apply -k ../project/k3d   # deployed with the rest of the project (see ../project)
kubectl get cronjob,jobs -n project

# Run it now instead of waiting for the next hour
kubectl create job -n project --from=cronjob/wiki-todo-generator wiki-manual
kubectl logs -n project job/wiki-manual   # Created todo: Read https://en.wikipedia.org/wiki/...
```

## Configuration

There are no defaults in the code; the job stops if a variable is missing.

- `RANDOM_ARTICLE_URL` - `https://en.wikipedia.org/wiki/Special:Random`, in the CronJob
- `TODOS_URL` - `http://todo-backend-svc:2345/todos`, in the CronJob
- `TODO_MAX_LENGTH` - from the ConfigMap `todo-backend-config`, the same limit todo-backend checks

## Files

- `main.go` - Gets a random article and creates the todo
- `Dockerfile` - Multi-stage build (Go 1.27, Alpine 3.24 with CA certificates)
- `manifests/cronjob.yaml` - Hourly CronJob
- `manifests/kustomization.yaml` - Kustomize base, part of [`../project`](../project)

## Technologies

- **Language:** Go 1.27 (standard library only)
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes CronJob
