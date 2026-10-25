# Database backup

A CronJob that dumps the todo database of the [project](../project) once a day and uploads the dump to Azure Blob Storage (exercise 3.10). The course material uses Google Cloud Storage; this project runs on Azure.

## Exercises

- **3.10:** Daily `pg_dump` of the todo database to Azure Blob Storage, authenticated with workload identity

## How it works

`backup.sh` runs in an image based on the same `postgres:18.6-alpine3.24` as the database, so `pg_dump` matches the server version:

1. `pg_dump --no-owner --no-privileges | gzip` to `/tmp/todos-<UTC time>.sql.gz`, connecting with the standard Postgres variables (`PGHOST`, `PGDATABASE`, `PGUSER`, `PGPASSWORD`, …) from the ConfigMap `todo-backend-config` and the Secret `postgres-secret`.
2. **Workload identity**: the pod runs as the Kubernetes ServiceAccount `db-backup` with the label `azure.workload.identity/use: "true"`, so AKS mounts a token for that ServiceAccount and sets `AZURE_CLIENT_ID`, `AZURE_TENANT_ID` and `AZURE_FEDERATED_TOKEN_FILE`. The script exchanges the token at Entra ID for an access token of the managed identity `id-dwk-db-backup`, which trusts exactly this ServiceAccount through a federated credential.
3. `Put Blob` to `https://stdwkbackupssc.blob.core.windows.net/todo-backups/<name>` with that token.

There are no keys, connection strings or Secrets for Azure anywhere: the storage account does not even accept shared keys. The managed identity can only write and read blobs in the container `todo-backups` (`Storage Blob Data Contributor` on the container).

The material describes two ways to give the job access: a service account key stored as a Kubernetes Secret, or workload identity, "the recommended and safest approach". This uses workload identity.

| Course material (GCP) | Here (Azure) |
|---|---|
| Google Cloud Storage bucket | Storage account `stdwkbackupssc`, container `todo-backups` (Cool tier) |
| GCP service account `my-storage-sa` | Managed identity `id-dwk-db-backup` |
| `roles/storage.objectCreator`, `objectViewer` | `Storage Blob Data Contributor` on the container |
| Workload Identity for GKE | Microsoft Entra Workload ID for AKS: a federated credential for `system:serviceaccount:project:db-backup` |

The Azure resources are in [`../terraform/backup.tf`](../terraform/backup.tf); the CronJob and the ServiceAccount are the Kustomize component [`../project/components/db-backup`](../project/components/db-backup).

## Schedule

Every day at 03:00 Finnish time (`timeZone: Europe/Helsinki`). The cluster is stopped every night, so the run is usually missed. `startingDeadlineSeconds: 86400` makes Kubernetes start a missed run when the cluster is back, as long as it is less than a day late. `backoffLimit: 2` retries a failed run, which helps right after the cluster starts, when the database may not be ready yet. Blobs older than 30 days are deleted by a lifecycle rule of the storage account.

Only production backs up its database: the deployment workflow adds the component for `main` only, and the federated credential trusts the ServiceAccount of the namespace `project` only.

## Run, list and restore

```bash
# Back up now
kubectl create job -n project --from=cronjob/db-backup db-backup-manual
kubectl logs -n project job/db-backup-manual   # Uploaded https://stdwkbackupssc.blob.core.windows.net/todo-backups/todos-...

# List and download (Storage Blob Data Reader on the container is enough)
az storage blob list --account-name stdwkbackupssc --container-name todo-backups --auth-mode login -o table
az storage blob download --account-name stdwkbackupssc --container-name todo-backups \
  --auth-mode login --name todos-<time>.sql.gz --file backup.sql.gz

# Restore into an empty database
gunzip -c backup.sql.gz | psql -v ON_ERROR_STOP=1 -d todos
```

## Files

- `backup.sh` - Dump, token exchange and upload
- `Dockerfile` - `postgres:18.6-alpine3.24` with `curl` and `jq`, runs as the `postgres` user

## Technologies

- PostgreSQL 18 `pg_dump`, Azure Blob Storage REST API (`x-ms-version: 2026-10-06`)
- Microsoft Entra Workload ID (Kubernetes ServiceAccount token → managed identity token)
- Kubernetes CronJob
