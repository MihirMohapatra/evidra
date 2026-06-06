variable "environment" {
  type = string
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

variable "instance_class" {
  type = string
}

variable "allocated_storage" {
  type = number
}

variable "max_allocated_storage" {
  type    = number
  default = 100
}

variable "db_name" {
  type = string
}

variable "db_master_username" {
  type      = string
  sensitive = true
}

variable "db_master_password" {
  type      = string
  sensitive = true
  default   = ""
}

variable "db_port" {
  type = number
}

variable "engine_version" {
  type = string
}

variable "parameter_group_family" {
  type = string
}

variable "backup_retention_days" {
  type    = number
  default = 7
}

variable "skip_final_snapshot" {
  type    = bool
  default = false
}

variable "deletion_protection" {
  type    = bool
  default = true
}

variable "create_read_replica" {
  type    = bool
  default = false
}

variable "tags" {
  type    = map(string)
  default = {}
}
