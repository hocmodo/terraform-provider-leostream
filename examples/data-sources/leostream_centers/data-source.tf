# Copyright (c) HashiCorp, Inc.

data "leostream_center_ds" "center_ds" {
  id = 51
}

# Output the ID of the image with the name "image_name"
output "center_ds_image_id" {
  value = element([for image in data.leostream_center_ds.center_ds.images : image if "${image.name}" == "image_name"], 0).id
}
