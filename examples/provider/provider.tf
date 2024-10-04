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
  host     = "https://leostream-broker"
  username = "api"
  password = "System@123"
}
