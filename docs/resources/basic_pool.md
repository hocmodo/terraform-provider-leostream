---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "leostream_basic_pool Resource - leostream"
subcategory: ""
description: |-
  The basic pool resource allows you to manage Leostream pools. Basic pools are used to group desktops together for management and provisioning.
---

# leostream_basic_pool (Resource)

The basic pool resource allows you to manage Leostream pools. Basic pools are used to group desktops together for management and provisioning.

## Example Usage

```terraform
# Copyright (c) HashiCorp, Inc.

resource "leostream_basic_pool" "pool_1" {

  name         = "Basic desktop pool 1"
  display_name = "Test"

  pool_definition = {
    restrict_by    = "A"
    parent_pool_id = 1
    server_ids     = []
    attributes = [
      {
        vm_table_field = "server_id"
        text_to_match  = "machine1"
        condition_type = "eq"
      }
    ]
  }

  provision = {
    provision_server_id          = 51
    provision_vm_name            = "desktop-{SEQUENCE}"
    provision_on_off             = 0
    provision_max                = 0
    provision_threshold          = 0
    provision_vm_display_name    = "vm 1"
    provision_vm_name_next_value = 0
    provision_vm_id              = 15
    mark_deletable               = 1
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) Display name of the pool.
- `id` (String) Unique identifier for the pool.
- `name` (String) Name of the pool.
- `notes` (String) Notes for the pool.
- `pool_definition` (Attributes) Pool definition (see [below for nested schema](#nestedatt--pool_definition))
- `provision` (Attributes) Container for parameters related to Provisioning.
				Provisioning parameters depends on what Centers are defined in the Connection Broker and which sets of values in every Center type (e.g. Azure, AWS, etc.) are defined. (see [below for nested schema](#nestedatt--provision))
- `running_desktops_threshold` (Number) Running and available desktops in the pool.

<a id="nestedatt--pool_definition"></a>
### Nested Schema for `pool_definition`

Optional:

- `attributes` (Attributes List) Array container for Pool attributes (restrict_by is 'A') or for LDAP attributes (restrict_by is 'Z', requires Active Directory Centers). (see [below for nested schema](#nestedatt--pool_definition--attributes))
- `never_rogue` (Number) 0 or 1: A boolean field indicating if desktops in this pool treat any user as the assigned user
- `parent_pool_id` (Number) ID of the parent pool
- `pool_attribute_join` (String) A or O: How do the pool attributes get joined:
						A = And
						O = Or
- `restrict_by` (String) Restrict by:
						A = by attribute (default)
- `server_ids` (List of Number) List of tag IDs defining this pool
- `use_vmotion` (Number) 0 or 1: A boolean field indicating whether VMs of this pool will vMotion to new host

<a id="nestedatt--pool_definition--attributes"></a>
### Nested Schema for `pool_definition.attributes`

Optional:

- `ad_attribute_field` (String) Desktop attribute, mandatory for LDAP attributes,
									see possible values for an AD Center in centers.get response, field ldap_attributes.
									annot exist if vm_table_field or vm_gpu_field is populated.
- `condition_type` (String) The search conditional:
									ip - "matches (CIDR notation)";
									np - "does not match (CIDR)";
									eq - "is equal to";
									ne - "is not equal to";
									gt - "is greater than";
									lt - "is less than";
									ct - "contains";
									nc - "does not contain";
									bw - "begins with";
									ew - "ends with".
- `text_to_match` (String) The free form text attribute
- `vm_gpu_field` (String) The GPU field to search; must be a column in the vm_gpu table. Cannot exist if vm_table_field or ad_attribute_field is populated.
- `vm_table_field` (String) The machine's attribute to search; must be a column in the vm table. Cannot exist if ad_attribute_field or vm_gpu_field is populated.
									name - Name;
									display_name - Display name;
									windows_name - Machine name;
									ip - Hostname or IP address;
									partition_names - Disk partition name;
									partition_mount_points - Partition mount point;
									guest_os - Operating system;
									os_version - Operating system version;
									installed_protocols - Installed protocols;
									vc_memory_mb - Memory (in MB);
									vc_num_cpu - Number of CPUs;
									vc_num_ethernet_cards - Number of NICs;
									num_disks - Number of disks;
									computer_model - Computer model;
									bios_serial_number - BIOS serial number;
									max_clock_speed - CPU speed (GHz);
									notes - Notes;
									vc_annotation - Center "Notes";
									tag_filter - Tags;
									server_id - Servers.



<a id="nestedatt--provision"></a>
### Nested Schema for `provision`

Optional:

- `mark_deletable` (Number) 0 or 1: Specifies whether to initialize newly-provisioned desktops as 'deletable'.
- `provision_limits_enforce` (Number) 0 or 1: A boolean field indicating if Broker creates and deletes virtual machines to meet the start and max threshold.
- `provision_max` (Number) The maximum number of new machines that will be provisioned when the threshold is reached.
- `provision_on_off` (Number) A boolean field indicating if state of provisioning for this pool is:
						Running - provision according to thresholds
						Stopped - disabled by user or the Broker by error
- `provision_server_id` (Number) The ID of the server which will do the provisioning, or 0 if URL notification only
- `provision_tenant_id` (Number) The tenant to provision into
- `provision_threshold` (Number) Minimum number of available VMs before triggering provisioning.
- `provision_url` (String) The URL to notify when a new machine is provisioned.
- `provision_vm_display_name` (String) The display name of the VM to be provisioned.
- `provision_vm_id` (Number) The ID of the server which will do the provisioning, or 0 if URL notification only
- `provision_vm_name` (String) The name of the VM to be provisioned.
- `provision_vm_name_next_value` (Number) The next value for sequential VM names

## Import

Import is supported using the following syntax:

```shell
# Copyright (c) HashiCorp, Inc.

# Order can be imported by specifying the numeric identifier.

terraform import leostream_basic_pool 123
```
