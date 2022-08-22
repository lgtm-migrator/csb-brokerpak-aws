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
output "status" {
  value = format(
    "created db %s (id: %s) on server %s URL: https://%s.console.aws.amazon.com/rds/home?region=%s#database:id=%s;is-cluster=false",
    aws_rds_cluster.default.database_name,
    aws_rds_cluster.default.id,
    aws_rds_cluster.default.endpoint,
    var.region,
    var.region,
    aws_rds_cluster.default.id,
  )
  sensitive = true
}

output "cluster_arn" {
  description = "Amazon Resource Name (ARN) of cluster"
  value       = try(aws_rds_cluster.default.arn, "")
}


output "cluster_master_password" {
  description = "The database master password"
  value       = try(aws_rds_cluster.default.master_password, "")
  sensitive   = true
}

output "cluster_master_username" {
  description = "The database master username"
  value       = try(aws_rds_cluster.default.master_username, "")
  sensitive   = true
}
