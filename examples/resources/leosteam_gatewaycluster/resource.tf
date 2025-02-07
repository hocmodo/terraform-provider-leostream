# Copyright (c) HashiCorp, Inc.

resource "leostream_gatewaycluster" "gwcluster" {
  name            = "gateway_us_east_1"
  address         = "gatewaycluster.public.address"
  notes           = "This is a gateway in the us-east-1 region"
  load_balance_via = "E"
}
