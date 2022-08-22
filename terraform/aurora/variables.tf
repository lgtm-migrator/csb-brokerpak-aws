locals {
  cluster_identifier = "ga-service-test-${replace(basename(path.cwd), "_", "-")}"
  engine             = var.engine == null ? "aurora-mysql" : var.engine
  engine_version     = var.engine_version == null ? "5.7" : var.engine_version
  is_serverless      = var.engine_mode == "serverless"
  port               = coalesce(var.port, (local.engine == "aurora-postgresql" ? 5432 : 3306))

  tags = {
    Owner       = "VMware"
    Environment = "development"
  }
}

// SET TF_VAR_aws_access_key_id environment variable
variable "AWS_ACCESS_KEY_ID" {
  type = string
}

// SET TF_VAR_aws_secret_access_key environment variable
variable "AWS_SECRET_ACCESS_KEY" {
  type = string
}

variable "region" {
  type    = string
  default = "us-west-2"
}

variable "port" {
  type    = number
  default = null
}

variable "engine" {
  description = "The name of the database engine to be used for this DB cluster. Defaults to `aurora`. Valid Values: `aurora`, `aurora-mysql`, `aurora-postgresql`, `mysql`, `postgres`. (Note that `mysql` and `postgres` are Multi-AZ RDS clusters)."
  type        = string
  default     = null
}

variable "engine_version" {
  description = "The database engine version."
  type        = string
  default     = null
}

variable "engine_mode" {
  description = "The database engine mode. Valid values: `global` (only valid for Aurora MySQL 1.21 and earlier), `multimaster`, `parallelquery`, `provisioned`, `serverless`. Defaults to: provisioned"
  type        = string
  default     = null
}

variable "copy_tags_to_snapshot" {
  description = "Copy all Cluster `tags` to snapshots"
  type        = bool
  default     = false
}