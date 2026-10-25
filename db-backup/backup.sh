#!/bin/sh
# shellcheck shell=busybox
# Dumps the todo database and uploads it to Azure Blob Storage (exercise 3.10).
#
# Database connection: the standard Postgres variables (PGHOST, PGDATABASE, PGUSER, PGPASSWORD, ...).
# Azure access: workload identity. The AKS webhook adds AZURE_CLIENT_ID, AZURE_TENANT_ID,
# AZURE_AUTHORITY_HOST and AZURE_FEDERATED_TOKEN_FILE (a Kubernetes ServiceAccount token),
# which is exchanged for an Entra ID token for Azure Storage. No keys are involved.
set -eu
set -o pipefail

: "${STORAGE_ACCOUNT:?STORAGE_ACCOUNT is not set}"
: "${CONTAINER:?CONTAINER is not set}"
: "${PGDATABASE:?PGDATABASE is not set}"
: "${AZURE_FEDERATED_TOKEN_FILE:?Workload identity is not set up for this pod}"

name="$PGDATABASE-$(date -u +%Y-%m-%dT%H%M%SZ).sql.gz"
file="/tmp/$name"

echo "Dumping database $PGDATABASE from $PGHOST"
pg_dump --no-owner --no-privileges | gzip > "$file"
echo "Dump is $(wc -c < "$file") bytes"

# Exchange the ServiceAccount token for an Azure Storage access token
token=$(curl -sS --fail-with-body -X POST \
  "${AZURE_AUTHORITY_HOST%/}/$AZURE_TENANT_ID/oauth2/v2.0/token" \
  --data-urlencode "client_id=$AZURE_CLIENT_ID" \
  --data-urlencode "grant_type=client_credentials" \
  --data-urlencode "scope=https://storage.azure.com/.default" \
  --data-urlencode "client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer" \
  --data-urlencode "client_assertion@$AZURE_FEDERATED_TOKEN_FILE" | jq -r .access_token)

# Put Blob: one request is enough for dumps up to 5000 MiB
url="https://$STORAGE_ACCOUNT.blob.core.windows.net/$CONTAINER/$name"
curl -sS --fail-with-body -X PUT "$url" \
  -H "Authorization: Bearer $token" \
  -H "x-ms-version: 2026-10-06" \
  -H "x-ms-blob-type: BlockBlob" \
  -H "Content-Type: application/gzip" \
  --data-binary "@$file"

echo "Uploaded $url"
