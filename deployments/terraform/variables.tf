variable "environment" {
  description = "Deployment environment (dev, staging, prod)"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "tags" {
  description = "Common resource tags"
  type        = map(string)
  default     = {}
}

variable "vpc_id" {
  description = "VPC ID for infrastructure resources"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for RDS"
  type        = list(string)
}

variable "k8s_cidr_blocks" {
  description = "Kubernetes cluster CIDR blocks allowed to access infrastructure"
  type        = list(string)
}

variable "kubernetes_host" {
  description = "Kubernetes API server endpoint"
  type        = string
}

variable "kubernetes_cluster_ca_certificate" {
  description = "Kubernetes cluster CA certificate (base64)"
  type        = string
}

variable "kubernetes_token" {
  description = "Kubernetes service account token"
  type        = string
  sensitive   = true
}

# --- Database ---

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
}

variable "db_allocated_storage" {
  description = "Allocated storage for RDS (GB)"
  type        = number
  default     = 20
}

variable "db_name" {
  description = "Initial database name"
  type        = string
  default     = "evidra"
}

variable "db_master_username" {
  description = "RDS master username"
  type        = string
  sensitive   = true
}

variable "db_master_password" {
  description = "RDS master password"
  type        = string
  sensitive   = true
}

variable "db_port" {
  description = "Database port"
  type        = number
  default     = 5432
}

variable "db_engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "16.3"
}

variable "db_parameter_group_family" {
  description = "RDS parameter group family"
  type        = string
  default     = "postgres16"
}

# --- NATS ---

variable "nats_namespace" {
  description = "Kubernetes namespace for NATS"
  type        = string
  default     = "nats"
}

variable "nats_replicas" {
  description = "Number of NATS cluster replicas"
  type        = number
  default     = 3
}

variable "nats_resources" {
  description = "Resource requests/limits for NATS pods"
  type = object({
    requests = object({
      cpu    = string
      memory = string
    })
    limits = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    requests = {
      cpu    = "250m"
      memory = "256Mi"
    }
    limits = {
      cpu    = "500m"
      memory = "512Mi"
    }
  }
}

variable "nats_storage_size" {
  description = "Persistent volume size for NATS JetStream"
  type        = string
  default     = "10Gi"
}

# --- Storage ---

variable "storage_buckets" {
  description = "List of S3 bucket names to create"
  type        = list(string)
}

variable "storage_force_destroy" {
  description = "Allow force destroy of non-empty buckets"
  type        = bool
  default     = false
}

variable "storage_enable_versioning" {
  description = "Enable S3 versioning on buckets"
  type        = bool
  default     = true
}

variable "storage_enable_encryption" {
  description = "Enable default SSE-S3 encryption on buckets"
  type        = bool
  default     = true
}

variable "storage_kms_key_arn" {
  description = "KMS key ARN for S3 SSE-KMS encryption (empty uses SSE-S3)"
  type        = string
  default     = ""
}
