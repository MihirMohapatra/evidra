terraform {
  required_version = ">= 1.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.31"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.15"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }

  backend "s3" {
    bucket = "evidra-tfstate"
    key    = "terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = var.tags
  }
}

provider "kubernetes" {
  host                   = var.kubernetes_host
  cluster_ca_certificate = base64decode(var.kubernetes_cluster_ca_certificate)
  token                  = var.kubernetes_token
}

provider "helm" {
  kubernetes {
    host                   = var.kubernetes_host
    cluster_ca_certificate = base64decode(var.kubernetes_cluster_ca_certificate)
    token                  = var.kubernetes_token
  }
}

module "postgres" {
  source = "./modules/postgres"

  environment          = var.environment
  vpc_id               = var.vpc_id
  private_subnet_ids   = var.private_subnet_ids
  k8s_cidr_blocks      = var.k8s_cidr_blocks
  instance_class       = var.db_instance_class
  allocated_storage    = var.db_allocated_storage
  db_name              = var.db_name
  db_master_username   = var.db_master_username
  db_master_password   = var.db_master_password
  db_port              = var.db_port
  engine_version       = var.db_engine_version
  parameter_group_family = var.db_parameter_group_family
  tags                 = var.tags
}

module "nats" {
  source = "./modules/nats"

  environment     = var.environment
  namespace       = var.nats_namespace
  replicas        = var.nats_replicas
  resources       = var.nats_resources
  storage_size    = var.nats_storage_size
  tags            = var.tags
}

module "storage" {
  source = "./modules/storage"

  environment        = var.environment
  buckets            = var.storage_buckets
  force_destroy      = var.storage_force_destroy
  enable_versioning  = var.storage_enable_versioning
  enable_encryption  = var.storage_enable_encryption
  kms_key_arn        = var.storage_kms_key_arn
  tags               = var.tags
}
