// This is an example of how to create a Leostream pool using Terraform.


// This block is required to use the Leostream provider.
terraform {
  required_providers {
    leostream = {
      source = "registry.terraform.io/hocmodo/leostream"

    }
  }
}

// This block configures the Leostream provider.
provider "leostream" {
  host     = "https://192.168.178.79"
  username = "api"
  password = var.leostream_api_password
}

// This block defines a Leostream pool resource.
#resource "leostream_pool" "pl" {
  # name         = "tst-us-rhel8-cloud-g4dn.2xl-pool"
  # display_name = "TST | US | RHEL8-CLOUD | 8CPU | 32GB"

  # pool_definition = {
  #   restrict_by = "A"
  #   parent_pool_id = 4
  #   server_ids = [2]
  #   attributes = [
  #     {
  #       vm_table_field     = "name"
  #       ad_attribute_field = ""
  #       vm_gpu_field       = ""
  #       text_to_match      = "tst-useast1-vdi-leo-rhel8-cloud-g4dn.2xl-"
  #       condition_type     = "ct"
  #     }
  #   ]
  # }

#   provision = {
#     provision_server_id = 13
#     provision_vm_name   = "tst-useast1-vdi-leo-rhel8-cloud-g4dn.2xl-{SEQUENCE}"
#     center = {
#       name = "aws-center-us-east-1"
#       type = "amazon"
#       id   = 13
#       aws_size = "g4dn.2xlarge"
#     }
#     provision_on_off = 0
#     provision_max = 2
#     provision_threshold = 0
#     provision_vm_display_name = "rhel8-cloud-g4dn.2xl-{SEQUENCE}"
#     mark_deletable = 1
#   }

# }

# output "new_leostream_pool" {
#   value = leostream_pool.pl
# }

// This block defines a Leostream pool resource.
resource "leostream_pool" "awspool" {

  name         = "AWS desktop pool 1"
  display_name = "Test"

  pool_definition = {
    restrict_by = "A"
    parent_pool_id = 1
    server_ids = []
    attributes = [
      {
        vm_table_field     = "server_id"
        text_to_match      = "51"
        condition_type     = "eq"
      }
    ]
  }

  provision = {
    provision_server_id = 51
    provision_vm_name   = "desktop-{SEQUENCE}"
    center = {
      name = "aws-center-us-east-1"
      type = "amazon"
      id   = 51
    }
    provision_on_off = 0
    provision_max = 0
    provision_threshold = 0
    provision_vm_display_name = "Image name:'emr 5.23.0-ami-roller-7 hvm ebs'"
    provision_vm_name_next_value = 8
    provision_vm_id = 15
    mark_deletable = 1
  }
}
