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
  value = try(module.wmware_aurora_mysql.status, "")
  sensitive = false
}

output "cluster_arn" {
  description = "Amazon Resource Name (ARN) of cluster"
  value       = try(module.wmware_aurora_mysql.cluster_arn, "")
}


output "cluster_master_password" {
  description = "The database master password"
  value       = try(module.wmware_aurora_mysql.cluster_master_password, "")
  sensitive   = true
}

output "cluster_master_username" {
  description = "The database master username"
  value       = try(module.wmware_aurora_mysql.cluster_master_username, "")
  sensitive   = false
}

output "cluster_instances" {
  description = "A map of cluster instances"
  value       = try(module.wmware_aurora_mysql.cluster_instances, "")
  sensitive   = false
}
