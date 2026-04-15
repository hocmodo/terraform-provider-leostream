# Copyright (c) HashiCorp, Inc.

resource "leostream_aws_pool" "pool_1" {

  name         = "AWS desktop pool 1"
  display_name = "Test"

  pool_definition = {
    restrict_by    = "A"
    parent_pool_id = 1
    server_ids     = []
    attributes = [
      {
        vm_table_field = "server_id"
        text_to_match  = "51"
        condition_type = "eq"
      }
    ]
  }

  provision = {
    provision_server_id = 51
    provision_vm_name   = "desktop-{SEQUENCE}"
    center = {
      name = "aws-center-us-east-1"
      type = "amazon"
      id   = 51
    }
    provision_on_off             = 0
    provision_max                = 0
    provision_threshold          = 0
    provision_vm_display_name    = "aws-desktop-1"
    provision_vm_name_next_value = 8
    provision_vm_id              = 15
    mark_deletable               = 1
  }

  # Optional: Configure logging thresholds and history retention
  log = {
    log_information_threshold = 100
    log_warning_threshold     = 50
    log_error_threshold       = 10
    retain_history = {
      pool_history_age      = 90   # days
      pool_history_interval = 1440 # minutes
    }
  }

  # Note: pool_stats is a read-only computed field that provides
  # real-time statistics about the pool, including:
  # - total_vm, total_vm_running, total_vm_stopped, total_vm_suspended
  # - total_agent_running, total_logged_in, total_connected
  # - assigned_vm, available_vm, unavailable_vm
  # - counts_updated (timestamp)
}
