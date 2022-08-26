provider "aws" {
  region     = var.region
  access_key = var.AWS_ACCESS_KEY_ID
  secret_key = var.AWS_SECRET_ACCESS_KEY
}

module "wmware_aurora_postgresql_serverlessv2" {
  source         = "../modules/aurora"
  engine         = "aurora-postgresql"
  engine_version = "14.3"
  engine_mode    = "provisioned" # It is the default value but I want to be explicit
  instance_class = "db.serverless"

#  two instances with default values
  instances = {for i in range(2) : i => {}}

  serverlessv2_scaling_configuration = {
    min_capacity = 2
    max_capacity = 6
  }

}