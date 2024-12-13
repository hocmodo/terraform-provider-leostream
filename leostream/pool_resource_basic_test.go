// Copyright (c) HashiCorp, Inc.

package leostream

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPoolBasicResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read testing
            {
                Config: providerConfig + `
resource "leostream_basic_pool" "test" {
  name         = "Basic desktop pool 1"
  display_name = "Test 1"

  pool_definition = {
    restrict_by    = "A"
    server_ids     = []
    attributes = [
      {
        vm_table_field = "server_id"
        text_to_match  = "1"
        condition_type = "eq"
      }
    ]
  }
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify number of attributes
                    //resource.TestCheckResourceAttr("leostream_basic_pool.test", "attributes.#", "1"),
                    // Verify attributes
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "name", "Basic desktop pool 1"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "display_name", "Test 1"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "notes", ""),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "pool_definition.restrict_by", "A"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "pool_definition.attributes.0.vm_table_field", "server_id"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "pool_definition.attributes.0.text_to_match", "1"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "pool_definition.attributes.0.condition_type", "eq"),
                    // Verify if item has Computed attributes filled.
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "running_desktops_threshold", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "pool_definition.never_rogue", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "pool_definition.parent_pool_id", "1"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "pool_definition.pool_attribute_join", "A"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "pool_definition.use_vmotion", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.mark_deletable", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_limits_enforce", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_max", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_on_off", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_server_id", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_tenant_id", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_threshold", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_url", ""),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_vm_display_name", ""),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_vm_id", "0"),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_vm_name", ""),
                    resource.TestCheckResourceAttr("leostream_basic_pool.test", "provision.provision_vm_name_next_value", "0"),
                    // Verify dynamic values have any value set in the state.
                    resource.TestCheckResourceAttrSet("leostream_basic_pool.test", "id"),
                ),
            },
            // // ImportState testing
            // {
            //     ResourceName:      "hashicups_order.test",
            //     ImportState:       true,
            //     ImportStateVerify: true,
            //     // The last_updated attribute does not exist in the HashiCups
            //     // API, therefore there is no value for it during import.
            //     ImportStateVerifyIgnore: []string{"last_updated"},
            // },
//             // Update and Read testing
//             {
//                 Config: providerConfig + `
// resource "hashicups_order" "test" {
//   attributes = [
//     {
//       coffee = {
//         id = 2
//       }
//       quantity = 2
//     },
//   ]
// }
// `,
//                 Check: resource.ComposeAggregateTestCheckFunc(
//                     // Verify first order item updated
//                     resource.TestCheckResourceAttr("hashicups_order.test", "attributes.0.quantity", "2"),
//                     resource.TestCheckResourceAttr("hashicups_order.test", "attributes.0.coffee.id", "2"),
//                     // Verify first coffee item has Computed attributes updated.
//                     resource.TestCheckResourceAttr("hashicups_order.test", "attributes.0.coffee.description", ""),
//                     resource.TestCheckResourceAttr("hashicups_order.test", "attributes.0.coffee.image", "/packer.png"),
//                     resource.TestCheckResourceAttr("hashicups_order.test", "attributes.0.coffee.name", "Packer Spiced Latte"),
//                     resource.TestCheckResourceAttr("hashicups_order.test", "attributes.0.coffee.price", "350"),
//                     resource.TestCheckResourceAttr("hashicups_order.test", "attributes.0.coffee.teaser", "Packed with goodness to spice up your images"),
//                 ),
//             },
//             // Delete testing automatically occurs in TestCase
         },
    })
}
