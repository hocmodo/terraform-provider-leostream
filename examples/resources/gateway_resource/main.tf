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

resource "leostream_gateway" "gw" {
  name            = "gateway_us_east_1_1_1"
  address         = "192.168.178.105"
  address_private = ""
  notes           = "This is a gateway in EU-WEST-1"
}

output "leostream_gateway" {
  value = leostream_gateway.gw
}
