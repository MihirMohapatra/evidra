module "evidra" {
  source = "../.."

  environment = var.environment
  region      = var.region
  tags        = var.tags

  vpc_id             = var.vpc_id
  private_subnet_ids = var.private_subnet_ids
  k8s_cidr_blocks    = var.k8s_cidr_blocks

  kubernetes_host                   = var.kubernetes_host
  kubernetes_cluster_ca_certificate = var.kubernetes_cluster_ca_certificate
  kubernetes_token                  = var.kubernetes_token

  db_instance_class       = var.db_instance_class
  db_allocated_storage    = var.db_allocated_storage
  db_name                 = var.db_name
  db_master_username      = var.db_master_username
  db_master_password      = var.db_master_password
  db_port                 = var.db_port
  db_engine_version       = var.db_engine_version
  db_parameter_group_family = var.db_parameter_group_family

  nats_replicas   = var.nats_replicas
  nats_resources  = var.nats_resources
  nats_storage_size = var.nats_storage_size
  nats_namespace  = var.nats_namespace

  storage_buckets          = var.storage_buckets
  storage_force_destroy    = var.storage_force_destroy
  storage_enable_versioning = var.storage_enable_versioning
  storage_enable_encryption = var.storage_enable_encryption
  storage_kms_key_arn      = var.storage_kms_key_arn
}

terraform {
  required_version = ">= 1.6"

  backend "s3" {
    bucket = "evidra-tfstate"
    key    = "prod/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = var.tags
  }
}
