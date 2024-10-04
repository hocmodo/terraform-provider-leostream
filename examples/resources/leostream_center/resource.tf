# Copyright (c) HashiCorp, Inc.

resource "leostream_center" "awscenter" {
  center_definition = {
    name                  = "aws-center-us-east-1"
    type                  = "amazon"
    allow_rogue_policy_id = 0
    vc_datacenter         = "us-east-1"
    vc_name               = "aws_access key"
    vc_password           = "aws_secret key"
    vc_auth_method        = "access_key"
    wait_inst_status      = 1
    wait_sys_status       = 1
  }
}
