provider "aws" {
  region     = local.primary_region
  access_key = var.AWS_ACCESS_KEY_ID
  secret_key = var.AWS_SECRET_ACCESS_KEY
}

provider "aws" {
  alias      = "secondary"
  region     = local.secondary_region
  access_key = var.AWS_ACCESS_KEY_ID
  secret_key = var.AWS_SECRET_ACCESS_KEY
}


resource "aws_rds_global_cluster" "global_cluster" {
  global_cluster_identifier = local.global_cluster_identifier
  engine                    = "aurora-postgresql"
  engine_version            = "14.3"
  database_name             = "my_db"
}

module "wmware_aurora_primary" {
  source                    = "../modules/aurora"
  database_name             = aws_rds_global_cluster.global_cluster.database_name
  engine                    = aws_rds_global_cluster.global_cluster.engine
  engine_version            = aws_rds_global_cluster.global_cluster.engine_version
  global_cluster_identifier = aws_rds_global_cluster.global_cluster.id
  instance_class            = "db.r5.large"

  #  two instances with default values
  instances = {for i in range(2) : i => {}}
}

module "wmware_aurora_secondary" {
  source                    = "../modules/aurora"
  #  No database name
  engine                    = aws_rds_global_cluster.global_cluster.engine
  engine_version            = aws_rds_global_cluster.global_cluster.engine_version
  global_cluster_identifier = aws_rds_global_cluster.global_cluster.id
  instance_class            = "db.r5.large"
  is_primary_cluster        = false
  providers                 = {
    aws = aws.secondary
  }
  source_region = local.primary_region

  instances = {for i in range(2) : i => {}}

  depends_on = [
    module.wmware_aurora_primary
  ]
}

resource "random_string" "global_cluster_name" {
  length  = 10
  special = false
  upper   = false
}
