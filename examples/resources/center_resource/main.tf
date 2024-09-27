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

resource "leostream_center" "awscenter" {
  center_definition = {
    name                  = "aws-center-us-east-1"
    type                  = "amazon"
    allow_rogue_policy_id = 0
    vc_datacenter         = "us-east-1"
    vc_name               = var.center_aws_key
    vc_password           = var.center_password
    vc_auth_method        = "access_key"
    wait_inst_status      = 1
    wait_sys_status       = 1
  }
}
