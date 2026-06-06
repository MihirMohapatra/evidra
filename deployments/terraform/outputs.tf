output "postgres_endpoint" {
  description = "RDS PostgreSQL endpoint (write)"
  value       = module.postgres.endpoint
}

output "postgres_reader_endpoint" {
  description = "RDS PostgreSQL reader endpoint (read replicas)"
  value       = module.postgres.reader_endpoint
}

output "postgres_port" {
  description = "RDS PostgreSQL port"
  value       = module.postgres.port
}

output "postgres_security_group_id" {
  description = "RDS security group ID"
  value       = module.postgres.security_group_id
}

output "nats_url" {
  description = "NATS cluster URL for clients"
  value       = module.nats.url
}

output "nats_metrics_url" {
  description = "NATS monitoring endpoint"
  value       = module.nats.metrics_url
}

output "nats_namespace" {
  description = "NATS Kubernetes namespace"
  value       = module.nats.namespace
}

output "storage_endpoint" {
  description = "S3-compatible storage endpoint"
  value       = module.storage.endpoint
}

output "storage_buckets" {
  description = "Created S3 bucket names"
  value       = module.storage.bucket_names
}

output "storage_bucket_arns" {
  description = "Created S3 bucket ARNs"
  value       = module.storage.bucket_arns
}

output "database_url_secret_name" {
  description = "AWS Secrets Manager secret name for the database connection string"
  value       = module.postgres.secret_name
}
