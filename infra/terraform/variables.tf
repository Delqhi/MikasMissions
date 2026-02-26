variable "project_name" {
  description = "Project name"
  type        = string
  default     = "mikasmissions"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "kubeconfig_path" {
  description = "Path to kubeconfig"
  type        = string
  default     = "~/.kube/config"
}
