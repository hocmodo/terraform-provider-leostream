# Copyright (c) HashiCorp, Inc.

data "leostream_gateways" "gateway_list" {}

output "gateway_list_output" {
  value = data.leostream_gateways.gateway_list
}
