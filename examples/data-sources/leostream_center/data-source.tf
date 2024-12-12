# Copyright (c) HashiCorp, Inc.

data "leostream_center_ds" "center_ds" {
  id = 51
}

# Output the ID of the image with the name "image_name"
output "center_ds_image_id" {
  value = element([for image in data.leostream_center_ds.center_ds.images : image if "${image.name}" == "emr 5.23.0-ami-roller-7 hvm ebs"], 0).id
}

# Output if there is an instance type t2.micro in the center
output "center_ds_aws_size_available" {
  value = contains(data.leostream_center_ds.center_ds.center_info.aws_sizes, "t2.micro")
}

# Output if there is an image called something in the center
output "center_ds_image_available" {
  value = contains(data.leostream_center_ds.center_ds.images[*].name, "emr 5.23.0-ami-roller-7 hvm ebs")
}

# Output AWS security groups in the center
output "center_sec_groups" {
  value = data.leostream_center_ds.center_ds.center_info.aws_sub_nets
}
