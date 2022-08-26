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
