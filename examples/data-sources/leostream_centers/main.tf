# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    leostream = {
      source = "registry.terraform.io/hocmodo/leostream"
    }
  }
}

provider "leostream" {
  host     = "https://192.168.178.79"
  username = "api"
  password = var.leostream_api_password
}

data "leostream_center_ds" "center_ds" {
  id = 51
}

# Output the ID of the image with the name "emr 5.23.0-ami-roller-7 hvm ebs"
output "center_ds_image_id" {
  value = element([for image in data.leostream_center_ds.center_ds.images : image if "${image.name}" == "emr 5.23.0-ami-roller-7 hvm ebs"], 0).id
}
