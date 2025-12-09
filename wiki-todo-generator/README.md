# Wiki Todo Generator

A scheduled job that fetches random Wikipedia articles and automatically creates todo items in the todo-backend service. Runs as a Kubernetes CronJob to add daily reading suggestions.

## Exercises

This application fulfills the requirements for:
- **Exercise 2.9:** Daily todos from Wikipedia API
- **Exercise 2.10:** CronJob for scheduled Wikipedia todo fetcher

## Features

- Fetches random Wikipedia articles using the Wikipedia API
- Creates todos with format: "Read https://en.wikipedia.org/wiki/[Article_Title]"
- Runs on a schedule (hourly by default)
- Proper User-Agent header for Wikipedia API compliance
- Integrates with todo-backend REST API
- Automatic error handling and retry logic (via Kubernetes)

## How It Works

1. Queries Wikipedia API for a random article
2. Extracts the article title
3. Constructs a Wikipedia URL
4. Posts a new todo to the todo-backend service
5. Exits (Kubernetes CronJob will schedule the next run)

## Running Locally

### Using Docker

Build the image:
```bash
docker build -t wiki-todo-generator:v1.0 .
```

Run manually (requires todo-backend to be running):
```bash
docker run \
  -e BACKEND_URL="http://localhost:3000/todos" \
  wiki-todo-generator:v1.0
```

Expected output:
```
Random Wikipedia article: https://en.wikipedia.org/wiki/Atlantic_Ocean
Successfully created todo: Read https://en.wikipedia.org/wiki/Atlantic_Ocean
```

### Using Go

```bash
# Set the backend URL
export BACKEND_URL="http://localhost:3000/todos"

# Run directly
go run main.go
```

## Deploying to Kubernetes

### Local Deployment (k3d)

#### Prerequisites
- k3d cluster running
- kubectl configured
- todo-backend service deployed and accessible

#### Import image to k3d
```bash
k3d image import wiki-todo-generator:v1.0 -c k3s-default
```

#### Deploy CronJob
```bash
kubectl apply -f manifests/cronjob.yaml
```

#### Check CronJob status
```bash
# View CronJob details
kubectl get cronjobs -n project

# List jobs created by the CronJob
kubectl get jobs -n project

# View logs from the most recent run
kubectl logs -n project -l job-name=wiki-todo-generator-<timestamp>
```

#### Manually trigger a job (for testing)
```bash
kubectl create job -n project --from=cronjob/wiki-todo-generator wiki-todo-test
kubectl logs -n project -l job-name=wiki-todo-test -f
```

#### Delete CronJob
```bash
kubectl delete -f manifests/cronjob.yaml
```

### Cloud Deployment (Azure AKS)

#### Prerequisites
- Azure CLI installed and authenticated (`az login`)
- AKS cluster created (see [terraform/README.md](../terraform/README.md))
- kubectl configured for AKS (`az aks get-credentials`)
- Azure Container Registry (ACR) access
- todo-backend service deployed in AKS

#### Build and Push to ACR
```bash
# Get ACR name from Terraform output
ACR_NAME=$(cd ../terraform/bootstrap && terraform output -raw acr_name)

# Login to ACR
az acr login --name $ACR_NAME

# Build and tag image
docker build -t wiki-todo-generator:v1.0 .
docker tag wiki-todo-generator:v1.0 $ACR_NAME.azurecr.io/wiki-todo-generator:v1.0

# Push to ACR
docker push $ACR_NAME.azurecr.io/wiki-todo-generator:v1.0
```

#### Update manifests for AKS
Update the image reference in `manifests/cronjob.yaml`:
```yaml
spec:
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: wiki-todo-generator
              image: <acr-name>.azurecr.io/wiki-todo-generator:v1.0
```

#### Deploy to AKS
```bash
kubectl apply -f manifests/cronjob.yaml
```

#### Check CronJob in AKS
```bash
# View CronJob details
kubectl get cronjobs -n project

# List jobs
kubectl get jobs -n project

# View logs from the most recent run
kubectl logs -n project -l job-name=wiki-todo-generator-<timestamp>
```

#### Manually trigger a job in AKS (for testing)
```bash
kubectl create job -n project --from=cronjob/wiki-todo-generator wiki-todo-test-aks
kubectl logs -n project -l job-name=wiki-todo-test-aks -f
```

#### Delete CronJob from AKS
```bash
kubectl delete -f manifests/cronjob.yaml
```

## Configuration

The application uses the following environment variables:

- `BACKEND_URL` - URL of the todo-backend service (default: http://todo-backend-svc:2345/todos)

### CronJob Schedule

The default schedule is `"0 * * * *"` which runs every hour at minute 0 (e.g., 1:00, 2:00, 3:00).

You can customize the schedule in the CronJob manifest:
```yaml
spec:
  schedule: "0 * * * *"  # Every hour
  # schedule: "0 9 * * *"  # Daily at 9:00 AM
  # schedule: "*/30 * * * *"  # Every 30 minutes
  # schedule: "0 0 * * 0"  # Weekly on Sunday at midnight
```

## Wikipedia API Integration

The application uses the Wikipedia REST API with the following endpoint:
```
https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnnamespace=0&rnlimit=1
```

**API Response Example:**
```json
{
  "query": {
    "random": [
      {
        "id": 12345,
        "title": "Atlantic Ocean"
      }
    ]
  }
}
```

**Important:** The application sets a proper User-Agent header as required by Wikipedia's API guidelines:
```
User-Agent: WikiTodoGenerator/1.0 (Kubernetes CronJob)
```

## Integration with Todo-Backend

The generator creates todos by sending POST requests to the todo-backend service:

**Request:**
```json
{
  "text": "Read https://en.wikipedia.org/wiki/Atlantic_Ocean"
}
```

The backend validates the text (must be ≤140 characters) and creates the todo.

## Files

- `main.go` - Main application code with Wikipedia API and todo-backend integration
- `go.mod` - Go module definition
- `Dockerfile` - Multi-stage Docker build
- `manifests/cronjob.yaml` - Kubernetes CronJob configuration

## Technologies

- **Language:** Go 1.21
- **Container:** Docker (multi-stage build with Alpine Linux)
- **Orchestration:** Kubernetes CronJob
- **External API:** Wikipedia REST API

## Error Handling

The application exits with a non-zero status code on errors:
- Wikipedia API unavailable
- Invalid API response
- Todo-backend unavailable
- Todo creation failed

Kubernetes will retry failed jobs according to the `restartPolicy: OnFailure` setting.

## Monitoring CronJob Execution

### View CronJob History
```bash
# See recent job executions
kubectl get jobs -n project -l app=wiki-todo-generator

# View logs from all job runs
kubectl logs -n project -l app=wiki-todo-generator --tail=50
```

### Check for Failed Jobs
```bash
# List failed jobs
kubectl get jobs -n project --field-selector status.successful=0

# Describe a specific job for error details
kubectl describe job -n project <job-name>
```

## Example Todo Items Created

The CronJob creates todos like:
- "Read https://en.wikipedia.org/wiki/Atlantic_Ocean"
- "Read https://en.wikipedia.org/wiki/Machine_Learning"
- "Read https://en.wikipedia.org/wiki/Solar_System"
- "Read https://en.wikipedia.org/wiki/Ancient_Rome"

Each todo is unique and educational, providing daily learning opportunities!

## Testing Tips

1. **Test the Wikipedia API integration:**
   ```bash
   curl "https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnnamespace=0&rnlimit=1"
   ```

2. **Verify todo-backend connectivity:**
   ```bash
   kubectl port-forward -n project svc/todo-backend-svc 2345:2345
   curl http://localhost:2345/todos
   ```

3. **Manually trigger and watch:**
   ```bash
   kubectl create job -n project --from=cronjob/wiki-todo-generator manual-test
   kubectl logs -n project -l job-name=manual-test -f
   ```

4. **Check generated todos:**
   ```bash
   curl http://localhost:8081/api/todos | jq '.[] | select(.text | startswith("Read https://en.wikipedia.org"))'
   ```
