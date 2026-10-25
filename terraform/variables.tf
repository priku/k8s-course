variable "resource_group_name" {
  description = "Existing resource group for the course resources (created with the budget, see README.md)"
  type        = string
  default     = "rg-dwk-aks-sc"
}

variable "node_resource_group_name" {
  description = "Resource group that AKS creates for the nodes; the course budget covers it by name"
  type        = string
  default     = "rg-dwk-aks-nodes-sc"
}

variable "cluster_name" {
  description = "Name of the AKS cluster"
  type        = string
  default     = "aks-dwk-sc"
}

variable "kubernetes_version" {
  description = "Kubernetes minor version; AKS picks the latest patch"
  type        = string
  default     = "1.37"
}

variable "node_count" {
  description = "Number of nodes in the default node pool"
  type        = number
  default     = 2
}

variable "vm_size" {
  description = "Node size: the latest generation (v7) with 2 vCPU and 8 GiB, enough for the Istio gateway and two Postgres databases"
  type        = string
  default     = "Standard_D2as_v7"
}

variable "acr_name" {
  description = "Globally unique name of the container registry; the AKS manifests refer to <name>.azurecr.io"
  type        = string
  default     = "acrdwkakssc"
}

variable "stop_time_zone" {
  description = "Windows time zone of the nightly stop schedule"
  type        = string
  default     = "FLE Standard Time" # Helsinki, follows daylight saving time
}

variable "tags" {
  description = "Tags for all resources"
  type        = map(string)
  default = {
    project = "dwk-k8s-course"
    purpose = "course"
    owner   = "priku"
  }
}
