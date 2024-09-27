# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    leostream = {
      source = "registry.terraform.io/hocmodo/leostream"
    }
  }
}

provider "leostream" {}

data "leostream_centers" "example" {}
