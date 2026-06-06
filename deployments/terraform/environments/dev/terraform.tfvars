environment = "dev"
region      = "us-east-1"

tags = {
  Environment = "dev"
  Project     = "evidra"
  ManagedBy   = "terraform"
}

# VPC — must match your dev cluster VPC
vpc_id             = "vpc-0a1b2c3d4e5f"
private_subnet_ids = ["subnet-0a1b2c3d", "subnet-0e4f5a6b", "subnet-0c7d8e9f"]
k8s_cidr_blocks    = ["10.0.0.0/16"]

# Kubernetes connection
kubernetes_host                   = "https://<dev-cluster>.eks.amazonaws.com"
kubernetes_cluster_ca_certificate = "<base64-ca-cert>"
kubernetes_token                  = "<sa-token>"

# Database — small single-AZ instance
db_instance_class       = "db.t3.medium"
db_allocated_storage    = 20
db_name                 = "evidra"
db_master_username      = "evidra"
db_master_password      = "changeme-dev-password"
db_port                 = 5432
db_engine_version       = "16.3"
db_parameter_group_family = "postgres16"

# NATS — single replica for dev
nats_replicas     = 1
nats_storage_size = "5Gi"
nats_namespace    = "nats"
nats_resources = {
  requests = {
    cpu    = "250m"
    memory = "256Mi"
  }
  limits = {
    cpu    = "500m"
    memory = "512Mi"
  }
}

# Storage
storage_buckets           = ["questionnaires", "evidence"]
storage_force_destroy     = true
storage_enable_versioning = true
storage_enable_encryption = true
storage_kms_key_arn       = ""
