# Wikipedia with init and sidecar

An app that serves Wikipedia pages, made of three containers in one Pod that share an `emptyDir` volume (exercise 5.4). There is no code of our own: everything is in [`manifests/deployment.yaml`](manifests/deployment.yaml).

| Container | Kind | Does |
|---|---|---|
| `fetch-first-page` | Init container | Fetches [the Kubernetes page](https://en.wikipedia.org/wiki/Kubernetes) and exits. The other containers start only after this has succeeded. |
| `fetch-random-page` | Sidecar | Waits for a random time between 5 and 15 minutes, fetches [a random page](https://en.wikipedia.org/wiki/Special:Random) and starts over |
| `nginx` | Main container | Serves whatever is in the volume |

The sidecar is a native Kubernetes one: an init container with `restartPolicy: Always`. It starts before nginx and keeps running beside it. It writes a new page to a temporary file and moves it into place, so nginx never serves a half-written page, and a failed fetch leaves the previous page in place.

```bash
kubectl apply -f manifests/namespace.yaml -f manifests/
curl -H 'Host: wikipedia.localhost' http://localhost:8081/   # or http://wikipedia.localhost:8081 in a browser
kubectl logs -n wikipedia deployment/wikipedia-dep -c fetch-random-page
```

The page is one HTML file, so its styles and images are missing. The volume is an `emptyDir`, so a new Pod starts over from the Kubernetes page.

## Test results (k3d)

| Test | Result |
|---|---|
| The first page after the Pod starts | `Kubernetes - Wikipedia`, and the Pod is `2/2` |
| The random fetch, with the delay shortened to 5–10 seconds for the test | Five different random pages in a row, without restarts |
| The real delay of 300–900 seconds, 2 000 draws in the same image | Smallest 300, largest 900 |

The delays are the env variables `MIN_DELAY_SECONDS` and `MAX_DELAY_SECONDS` of the sidecar.
