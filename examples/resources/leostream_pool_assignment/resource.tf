# Copyright (c) HashiCorp, Inc.

resource "leostream_pool_assignment" "poolassignment_1" {

  policy_id    = 2
  pool_id      = 148
  offer_filter = 1
  offer_filter_json = {
    filters = [
      {
        offer_filter_attribute = "isMemberOf",
        offer_filter_condition = "eq",
        offer_filter_value     = "test"
      }
    ],
    join = "o"
  }
  plan_protocol_id      = 2
  offer_quantity        = 1
  display_mode          = 0
  start_if_stopped      = 1
  on_assign_url         = "http://mycallback.com"
  on_assign_url_cb      = 1
  on_assign_url_timeout = 5
}
