# This is an example of a main.tf file that creates a Leostream center, gateway, and pool.
# The center is created with the name "aws-center-us-east-1" and the type "amazon".
# The gateway is created with the name "gateway_us_east_1_1" and the address "
# The pool is depending on the center and is created with the name "AWS desktop pool 2".

terraform {
  required_providers {
    leostream = {
      source = "registry.terraform.io/hocmodo/leostream"
    }
    time = {
      source = "hashicorp/time"
      version = "0.12.0"
    }
  }
}

# This is the Leostream provider needed for the Leostream resources
provider "leostream" {
  host     = "https://192.168.178.79"
  username = "api"
  password = var.leostream_api_password
}

# This is a time provider needed for the time_sleep resource
provider "time" {
  # Configuration options
}

# This is the Leostream center resource
resource "leostream_center" "awscenter" {
  center_definition = {
    name = var.center_name
    type = "amazon"
    allow_rogue_policy_id = 0
    vc_datacenter = "us-east-1"
    vc_name = var.center_aws_key
    vc_password = var.center_password
    vc_auth_method = "access_key"
    }
}

# This is the Leostream gateway resource
resource "leostream_gateway" "gw" {
  name = "gateway_us_east_1_1"
  address = "192.168.178.105"
  address_private = ""
  notes = "This is a gateway in EU-WEST-1"
}

# This is a time_sleep resource that waits for the center to be scanned
resource "time_sleep" "wait_for_center_scan" {
  create_duration = "20s"
  depends_on = [leostream_center.awscenter]
}

# This is the Leostream center data source for looking up images and sizes
data "leostream_center_ds" "center_ds" {
   id = leostream_center.awscenter.id
   depends_on = [time_sleep.wait_for_center_scan]
}

# This is the Leostream pool resource
resource "leostream_pool" "awspool" {
  name         = "AWS desktop pool 2"
  display_name = "AWS desktop pool 2"

  pool_definition = {
    restrict_by = "A"
    parent_pool_id = 1
    server_ids = []
     attributes = [
       {
         vm_table_field     = "server_id"
         text_to_match      = leostream_center.awscenter.id
         condition_type     = "eq"
       }
     ]
  }

  provision = {
    provision_server_id = leostream_center.awscenter.id
    provision_vm_name   = "desktop-{SEQUENCE}"
    center = {
      name = var.center_name
      type = "amazon"
      id   = leostream_center.awscenter.id
      aws_size = element([for type in data.leostream_center_ds.center_ds.center_info.aws_sizes: type if "${type}" == var.desktop_type], 0)
      aws_sec_group = var.aws_sec_group_id
      aws_vpc_id = var.aws_vpc_id
      aws_iam_name = var.aws_iam_name
      aws_sub_net = var.aws_sub_net
    }
    provision_on_off = 0
    provision_max = 0
    provision_threshold = 0
    provision_vm_display_name = "Image name:'emr 5.23.0-ami-roller-7 hvm ebs' SBOM:http://blog.joustie.nl"
    provision_vm_name_next_value = 8
    provision_vm_id = element([for image in data.leostream_center_ds.center_ds.images: image if "${image.name}" == var.pool_image_name], 0).id
    mark_deletable = 1
  }

}
