environment = "prod"
region      = "us-east-1"

tags = {
  Environment = "prod"
  Project     = "evidra"
  ManagedBy   = "terraform"
}

# VPC — must match your prod cluster VPC
vpc_id             = "vpc-0f1e2d3c4b5a"
private_subnet_ids = ["subnet-0f1e2d3c", "subnet-0b4a5f6e", "subnet-0d7c8b9a"]
k8s_cidr_blocks    = ["10.1.0.0/16"]

# Kubernetes connection
kubernetes_host                   = "https://<prod-cluster>.eks.amazonaws.com"
kubernetes_cluster_ca_certificate = "<base64-ca-cert>"
kubernetes_token                  = "<sa-token>"

# Database — HA with multi-AZ and read replica
db_instance_class       = "db.r6g.large"
db_allocated_storage    = 100
db_name                 = "evidra"
db_master_username      = "evidra"
db_master_password      = ""  # Use random_password in module
db_port                 = 5432
db_engine_version       = "16.3"
db_parameter_group_family = "postgres16"

# NATS — 3-node cluster for HA
nats_replicas     = 3
nats_storage_size = "50Gi"
nats_namespace    = "nats"
nats_resources = {
  requests = {
    cpu    = "1"
    memory = "2Gi"
  }
  limits = {
    cpu    = "2"
    memory = "4Gi"
  }
}

# Storage
storage_buckets           = ["questionnaires", "evidence"]
storage_force_destroy     = false
storage_enable_versioning = true
storage_enable_encryption = true
storage_kms_key_arn       = "arn:aws:kms:us-east-1:123456789012:key/evidra-prod-s3-key"
