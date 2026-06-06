output "endpoint" {
  value = aws_db_instance.this.endpoint
}

output "reader_endpoint" {
  value = try(aws_db_instance.replica[0].endpoint, aws_db_instance.this.endpoint)
}

output "port" {
  value = var.db_port
}

output "security_group_id" {
  value = aws_security_group.rds.id
}

output "secret_name" {
  value = aws_secretsmanager_secret.database_url.name
}

output "secret_arn" {
  value = aws_secretsmanager_secret.database_url.arn
}
