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

data "leostream_centers" "centers" {}

output "centers_output" {
  value = data.leostream_centers.centers
}
