# Copyright 2020 Pivotal Software, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

resource "aws_rds_cluster" "default" {
  cluster_identifier      = local.cluster_identifier
  engine                  = local.engine
  engine_version          = local.is_serverless ? null : local.engine_version
  engine_mode             = var.engine_mode
  #  availability_zones      = ["us-west-2a", "us-west-2b", "us-west-2c"]
  database_name           = "mydb"
  master_username         = random_string.username.result
  master_password         = random_password.password.result
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  port                    = local.port

  lifecycle {
    ignore_changes = [
      # Ignore changes to availability_zones
      # This only worked when availability zones were matching all available zones in the region.
      # If provided less or none, there was always a diff and terraform wanted to recreate the cluster.
      availability_zones,
    ]
  }
}

resource "random_password" "password" {
  length           = 32
  special          = false
  override_special = "~_-."
}

resource "random_string" "username" {
  length           = 10
  special          = false
  override_special = "~_-."
}