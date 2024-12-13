# Copyright (c) HashiCorp, Inc.

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
  password = "System@123"
}
