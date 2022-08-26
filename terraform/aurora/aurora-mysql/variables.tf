#variable "engine" {
#  description = "The name of the database engine to be used for this DB cluster. Defaults to `aurora`. Valid Values: `aurora`, `aurora-mysql`, `aurora-postgresql`, `mysql`, `postgres`. (Note that `mysql` and `postgres` are Multi-AZ RDS clusters)."
#  type        = string
#}
#

// SET TF_VAR_aws_access_key_id environment variable
variable "AWS_ACCESS_KEY_ID" {
  type = string
  default = null
}

// SET TF_VAR_aws_secret_access_key environment variable
variable "AWS_SECRET_ACCESS_KEY" {
  type = string
  default = null
}

variable "region" {
  type    = string
  default = "us-west-2"
}
