variable "environment" {
  type = string
}

variable "namespace" {
  type    = string
  default = "nats"
}

variable "replicas" {
  type    = number
  default = 3
}

variable "nats_image_tag" {
  type    = string
  default = "2.10.18-alpine"
}

variable "chart_version" {
  type    = string
  default = "1.2.3"
}

variable "jetstream_mem_size" {
  type    = string
  default = "2Gi"
}

variable "jetstream_file_size" {
  type    = string
  default = "10Gi"
}

variable "storage_size" {
  type    = string
  default = "10Gi"
}

variable "storage_class" {
  type    = string
  default = ""
}

variable "resources" {
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

variable "tags" {
  type    = map(string)
  default = {}
}
