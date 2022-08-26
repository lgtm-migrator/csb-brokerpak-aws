# =====================================================================================================================
# ==================================== aws_rds_cluster ================================================================
# =====================================================================================================================

locals {
  cluster_identifier  = "ga-service-test-${replace(basename(path.cwd), "_", "-")}-${random_string.cluster_name.result}"
  port                = coalesce(var.port, (var.engine == "aurora-postgresql" ? 5432 : 3306))
  has_serverless_conf = length(keys(var.serverlessv2_scaling_configuration)) != 0

  tags = {
    Owner       = "VMware"
    Environment = "development"
  }
}

// SET TF_VAR_aws_access_key_id environment variable
variable "AWS_ACCESS_KEY_ID" {
  type    = string
  default = null
}

// SET TF_VAR_aws_secret_access_key environment variable
variable "AWS_SECRET_ACCESS_KEY" {
  type    = string
  default = null
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
}

variable "engine_mode" {
  description = "The database engine mode. Valid values: `global` (only valid for Aurora MySQL 1.21 and earlier), `multimaster`, `parallelquery`, `provisioned`, `serverless` (only valid for serverless v1)."
  type        = string
  default     = "provisioned"
}

variable "copy_tags_to_snapshot" {
  description = "Copy all Cluster `tags` to snapshots"
  type        = bool
  default     = false
}

variable "serverlessv2_scaling_configuration" {
  description = "Map of nested attributes with serverless v2 scaling properties. Only valid when `engine_mode` is set to `provisioned`"
  type        = map(string)
  default     = {}
}

variable "is_primary_cluster" {
  description = "Determines whether cluster is primary cluster with writer instance. Note: set to `false` replica elements in global cluster"
  type        = bool
  default     = true
}

variable "source_region" {
  description = "The source region for a replica DB cluster"
  type        = string
  default     = null
}

variable "database_name" {
  type    = string
  default = "mydb"
}

variable "global_cluster_identifier" {
  description = "The global cluster identifier specified on `aws_rds_global_cluster`"
  type        = string
  default     = null
}


# =====================================================================================================================
# ==================================== aws_rds_cluster_instances ======================================================
# =====================================================================================================================
variable "instances" {
  description = "Map of cluster instances"
  type        = any
  default     = {}
}


variable "instance_class" {
  description = "Instance type to use in instances. It can be overridden using the instance map"
  type        = string
  default     = ""
}

variable "publicly_accessible" {
  description = "Determines whether instances are publicly accessible. Default false"
  type        = bool
  default     = false
}

variable "apply_immediately" {
  description = "Specifies whether any cluster modifications are applied immediately, or during the next maintenance window. Default is `false`"
  type        = bool
  default     = false
}

# =====================================================================================================================
# =========================================== db_subnet_group  ========================================================
# =====================================================================================================================
variable "db_subnet_group_name" {
  description = "The name of the subnet group name"
  type        = string
  default     = ""
}

variable "db_subnet_group_subnets" {
  description = "List of subnet IDs used by database subnet group created"
  type        = list(string)
  default     = []
}
