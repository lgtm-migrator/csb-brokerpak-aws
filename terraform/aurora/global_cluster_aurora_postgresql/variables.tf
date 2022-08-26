locals {
  primary_region   = "eu-west-1"
  secondary_region = "us-east-1"
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

locals {
  global_cluster_identifier  = "ga-service-test-${replace(basename(path.cwd), "_", "-")}-${random_string.global_cluster_name.result}"
}
