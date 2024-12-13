# Copyright (c) HashiCorp, Inc.

data "leostream_policies" "policy_list" {}

output "policy_list_output" {
  value = data.leostream_policies.policy_list
}
