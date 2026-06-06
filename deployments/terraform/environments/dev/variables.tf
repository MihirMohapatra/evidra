variable "environment" {
  type = string
}

variable "region" {
  type    = string
  default = "us-east-1"
}

variable "tags" {
  type    = map(string)
  default = {}
}

variable "vpc_id" {
  type = string
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "k8s_cidr_blocks" {
  type = list(string)
}

variable "kubernetes_host" {
  type = string
}

variable "kubernetes_cluster_ca_certificate" {
  type = string
}

variable "kubernetes_token" {
  type      = string
  sensitive = true
}

variable "db_instance_class" {
  type = string
}

variable "db_allocated_storage" {
  type    = number
  default = 20
}

variable "db_name" {
  type    = string
  default = "evidra"
}

variable "db_master_username" {
  type      = string
  sensitive = true
}

variable "db_master_password" {
  type      = string
  sensitive = true
}

variable "db_port" {
  type    = number
  default = 5432
}

variable "db_engine_version" {
  type    = string
  default = "16.3"
}

variable "db_parameter_group_family" {
  type    = string
  default = "postgres16"
}

variable "nats_namespace" {
  type    = string
  default = "nats"
}

variable "nats_replicas" {
  type    = number
  default = 1
}

variable "nats_resources" {
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
  type    = string
  default = "5Gi"
}

variable "storage_buckets" {
  type = list(string)
}

variable "storage_force_destroy" {
  type    = bool
  default = true
}

variable "storage_enable_versioning" {
  type    = bool
  default = true
}

variable "storage_enable_encryption" {
  type    = bool
  default = true
}

variable "storage_kms_key_arn" {
  type    = string
  default = ""
}
