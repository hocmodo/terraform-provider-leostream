# Copyright (c) HashiCorp, Inc.

resource "leostream_aws_pool" "pool_1" {

  name         = "AWS desktop pool 1"
  display_name = "Test"

  pool_definition = {
    restrict_by    = "A"
    parent_pool_id = 1
    server_ids     = []
    attributes = [
      {
        vm_table_field = "server_id"
        text_to_match  = "51"
        condition_type = "eq"
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
    provision_on_off             = 0
    provision_max                = 0
    provision_threshold          = 0
    provision_vm_display_name    = "aws-desktop-1"
    provision_vm_name_next_value = 8
    provision_vm_id              = 15
    mark_deletable               = 1
  }
}
