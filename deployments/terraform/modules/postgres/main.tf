resource "aws_db_parameter_group" "this" {
  name        = "evidra-pg16-${var.environment}"
  family      = var.parameter_group_family
  description = "Evidra PostgreSQL 16 parameter group with pgvector"

  parameter {
    name  = "shared_preload_libraries"
    value = "vector"
  }

  parameter {
    name  = "vector.enable_stats"
    value = "1"
  }

  parameter {
    name         = "rds.force_ssl"
    value        = "1"
    apply_method = "pending-reboot"
  }

  tags = var.tags
}

resource "aws_db_subnet_group" "this" {
  name        = "evidra-${var.environment}"
  subnet_ids  = var.private_subnet_ids
  description = "Evidra ${var.environment} DB subnet group"

  tags = var.tags
}

resource "aws_security_group" "rds" {
  name        = "evidra-${var.environment}-rds"
  description = "Security group for Evidra RDS PostgreSQL"
  vpc_id      = var.vpc_id

  ingress {
    description = "PostgreSQL from Kubernetes"
    from_port   = var.db_port
    to_port     = var.db_port
    protocol    = "tcp"
    cidr_blocks = var.k8s_cidr_blocks
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.tags
}

resource "aws_db_instance" "this" {
  identifier     = "evidra-${var.environment}"
  engine         = "postgres"
  engine_version = var.engine_version
  instance_class = var.instance_class

  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage
  storage_type          = "gp3"
  storage_encrypted     = true

  db_name                        = var.db_name
  username                       = var.db_master_username
  password                       = var.db_master_password
  port                           = var.db_port
  manage_master_user_password    = false

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  parameter_group_name = aws_db_parameter_group.this.name

  backup_retention_period = var.backup_retention_days
  backup_window           = "03:00-04:00"
  maintenance_window      = "sun:04:00-sun:05:00"

  copy_tags_to_snapshot     = true
  skip_final_snapshot       = var.skip_final_snapshot
  final_snapshot_identifier = var.skip_final_snapshot ? null : "evidra-${var.environment}-final-${formatdate("YYYYMMDDhhmmss", timestamp())}"

  deletion_protection = var.deletion_protection

  enabled_cloudwatch_logs_exports = ["postgresql"]

  tags = var.tags
}

resource "aws_db_instance" "replica" {
  count = var.create_read_replica ? 1 : 0

  identifier = "evidra-${var.environment}-replica"
  replicate_source_db = aws_db_instance.this.identifier

  instance_class = var.instance_class
  storage_type   = "gp3"

  vpc_security_group_ids = [aws_security_group.rds.id]

  backup_retention_period = var.backup_retention_days
  backup_window           = "04:00-05:00"
  maintenance_window      = "sun:05:00-sun:06:00"

  copy_tags_to_snapshot = true

  deletion_protection = var.deletion_protection

  tags = var.tags
}

resource "random_password" "db_password" {
  count = var.db_master_password != "" ? 0 : 1

  length  = 24
  special = false
}

resource "aws_secretsmanager_secret" "database_url" {
  name = "evidra/${var.environment}/database-url"
  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "database_url" {
  secret_id = aws_secretsmanager_secret.database_url.id
  secret_string = jsonencode({
    url      = "postgres://${aws_db_instance.this.username}:${aws_db_instance.this.password}@${aws_db_instance.this.endpoint}:${var.db_port}/${var.db_name}?sslmode=require"
    endpoint = aws_db_instance.this.endpoint
    port     = var.db_port
    database = var.db_name
    username = aws_db_instance.this.username
  })
}
