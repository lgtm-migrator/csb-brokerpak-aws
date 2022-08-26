provider "aws" {
  region     = var.region
  access_key = var.AWS_ACCESS_KEY_ID
  secret_key = var.AWS_SECRET_ACCESS_KEY
}

module "wmware_aurora_mysql" {
  source         = "../modules/aurora"
  engine         = "aurora-mysql"
  engine_version = "5.7"
  instances      = {
    1 = {
      instance_class      = "db.r5.large"
      publicly_accessible = true
    }
    2 = {
      identifier     = "mysql-test-1"
      instance_class = "db.r5.large"
    }
    3 = {
      identifier     = "mysql-test-2"
      instance_class = "db.r5.large"
      promotion_tier = 15
    }
  }
}
