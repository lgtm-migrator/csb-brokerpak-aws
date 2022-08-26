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
  # Set preferred_backup_window here and not in the instances because will create an error
  cluster_identifier        = local.cluster_identifier
  engine                    = var.engine
  engine_version            = var.engine_version
  engine_mode               = var.engine_mode
  #  availability_zones      = ["us-west-2a", "us-west-2b", "us-west-2c"]
  database_name             = var.is_primary_cluster ? var.database_name : null
  master_username           = var.is_primary_cluster ? random_string.username.result : null
  master_password           = var.is_primary_cluster ? random_password.password.result : null
  backup_retention_period   = 5
  preferred_backup_window   = "07:00-09:00"
  port                      = local.port
  skip_final_snapshot       = true
  global_cluster_identifier = var.global_cluster_identifier


  dynamic "serverlessv2_scaling_configuration" {
    for_each = !local.has_serverless_conf ? [] : [var.serverlessv2_scaling_configuration]

    content {
      min_capacity = lookup(serverlessv2_scaling_configuration.value, "min_capacity", 1)
      max_capacity = lookup(serverlessv2_scaling_configuration.value, "max_capacity", 5)
    }
  }


  lifecycle {
    ignore_changes = [
      # Ignore changes to availability_zones
      # This only worked when availability zones were matching all available zones in the region.
      # If provided less or none, there was always a diff and terraform wanted to recreate the cluster.
      availability_zones,
      # See https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/rds_cluster#replication_source_identifier
      # Since this is used either in read-replica clusters or global clusters.
      replication_source_identifier,
      # See docs here https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/rds_global_cluster#new-global-cluster-from-existing-db-cluster
      global_cluster_identifier,
    ]
  }
}

resource "aws_rds_cluster_instance" "instance" {
  for_each = length(keys(var.instances)) != 0 ? var.instances : {}

  identifier            = lookup(each.value, "identifier", "${local.cluster_identifier}-${each.key}")
  cluster_identifier    = try(aws_rds_cluster.default.id, "")
  engine                = var.engine
  engine_version        = var.engine_version
  instance_class        = lookup(each.value, "instance_class", var.instance_class)
  publicly_accessible   = lookup(each.value, "publicly_accessible", var.publicly_accessible)
  apply_immediately     = lookup(each.value, "apply_immediately", var.apply_immediately)
  availability_zone     = lookup(each.value, "availability_zone", null)
  copy_tags_to_snapshot = lookup(each.value, "copy_tags_to_snapshot", var.copy_tags_to_snapshot)
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

resource "random_string" "cluster_name" {
  length  = 10
  special = false
  upper   = false
}
