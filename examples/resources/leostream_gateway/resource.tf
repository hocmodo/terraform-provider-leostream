# Copyright (c) HashiCorp, Inc.

resource "leostream_gateway" "gw" {
  name            = "gateway_us_east_1"
  address         = "gateway.private.address"
  address_private = ""
  notes           = "This is a gateway in the us-east-1 region"
}

output "leostream_gateway" {
  value = leostream_gateway.gw
}
