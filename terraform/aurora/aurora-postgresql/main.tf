provider "aws" {
  region     = var.region
  access_key = var.AWS_ACCESS_KEY_ID
  secret_key = var.AWS_SECRET_ACCESS_KEY
}

module "wmware_aurora_postgresql" {
  source         = "../modules/aurora"
  engine         = "aurora-postgresql"
  engine_version = "14.3"

  instances = {
    1 = {
      instance_class      = "db.r5.large"
      publicly_accessible = true
    }
    2 = {
      identifier     = "postgresql-test-1"
      instance_class = "db.r5.large"
    }
    3 = {
      identifier     = "postgresql-test-2"
      instance_class = "db.r5.large"
    }
  }
}