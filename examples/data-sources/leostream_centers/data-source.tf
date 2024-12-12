# Copyright (c) HashiCorp, Inc.

data "leostream_centers" "center_list" {}

output "center_list_output" {
  value = one(data.leostream_centers.center_list.centers[*].name)
}
