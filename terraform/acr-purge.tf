# ==============================================================================
# Weekly purge of old production images. Every deployment of main pushes three images
# (`project-<sha>`, and `main-<sha>` from before exercise 3.7); this ACR task deletes the
# ones older than a week, keeping the five newest of each repository, plus untagged
# manifests. Branch images are deleted when the branch is (project-cleanup.yaml), and
# hand-built version tags such as v2.9 do not match the filter.
# ==============================================================================
resource "azurerm_container_registry_task" "purge_old_images" {
  name                  = "purge-old-project-images"
  container_registry_id = azurerm_container_registry.acr.id
  tags                  = var.tags

  platform {
    os = "Linux"
  }

  encoded_step {
    task_content = base64encode(<<-YAML
      version: v1.1.0
      steps:
        - cmd: >-
            acr purge
            --filter 'todo-project:^(project|main)-[0-9a-f]{40}$'
            --filter 'todo-backend:^(project|main)-[0-9a-f]{40}$'
            --filter 'wiki-todo-generator:^(project|main)-[0-9a-f]{40}$'
            --ago 7d --keep 5 --untagged
          disableWorkingDirectoryOverride: true
          timeout: 3600
    YAML
    )
  }

  timer_trigger {
    name     = "weekly"
    schedule = "0 3 * * 1" # Mondays at 03:00 UTC
  }
}
