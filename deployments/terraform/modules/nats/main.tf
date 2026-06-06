resource "kubernetes_namespace" "this" {
  metadata {
    name = var.namespace
    labels = {
      name        = var.namespace
      environment = var.environment
      app         = "nats"
    }
  }
}

resource "helm_release" "nats" {
  name       = "nats"
  repository = "https://nats-io.github.io/helm-charts"
  chart      = "nats"
  version    = var.chart_version
  namespace  = kubernetes_namespace.this.metadata[0].name

  values = [
    yamlencode({
      nats = {
        image = {
          repository = "nats"
          tag        = var.nats_image_tag
        }
        jetstream = {
          enabled = true
          memStorage = {
            enabled = true
            size    = var.jetstream_mem_size
          }
          fileStorage = {
            enabled        = true
            size           = var.jetstream_file_size
            storageDirectory = "/data/jetstream"
            pvc = {
              enabled      = true
              storageSize  = var.storage_size
              storageClass = var.storage_class
            }
          }
        }
        cluster = {
          enabled = var.replicas > 1
          replicas = var.replicas
        }
        resources = var.resources
      }
      configReload = {
        enabled = true
        image = {
          repository = "natsio/nats-server-config-reloader"
          tag        = "0.14.0"
        }
      }
      exporter = {
        enabled = true
        image = {
          repository = "natsio/prometheus-nats-exporter"
          tag        = "0.15.0"
        }
      }
    })
  ]

  depends_on = [kubernetes_namespace.this]
}
