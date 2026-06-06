output "url" {
  value = "nats://nats.${var.namespace}.svc.cluster.local:4222"
}

output "metrics_url" {
  value = "http://nats.${var.namespace}.svc.cluster.local:8222"
}

output "namespace" {
  value = var.namespace
}
