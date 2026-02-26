resource "kubernetes_namespace" "platform" {
  metadata {
    name = "${var.project_name}-${var.environment}"
    labels = {
      app = var.project_name
      env = var.environment
    }
  }
}
