variable "resource_group_name" {
  description = "Name of the resource group"
  type        = string
  default     = "dwk-aks-rg"
}

variable "location" {
  description = "Azure region for resources"
  type        = string
  default     = "swedencentral"
}

variable "cluster_name" {
  description = "Name of the AKS cluster"
  type        = string
  default     = "dwk-aks-cluster"
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.30"
}

variable "node_count" {
  description = "Number of nodes in the default node pool"
  type        = number
  default     = 2
}

variable "vm_size" {
  description = "VM size for the nodes"
  type        = string
  default     = "Standard_B2s_v2" # Cost-effective for learning (v2 available in swedencentral)
}

variable "tags" {
  description = "Tags for resources"
  type        = map(string)
  default = {
    Environment = "Development"
    Course      = "DevOps-with-Kubernetes"
  }
}
