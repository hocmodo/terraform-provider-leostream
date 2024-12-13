// Copyright (c) HashiCorp, Inc.

package leostream

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCentersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "leostream_centers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.leostream_centers.test", "centers.#", "1"),
					// Verify the first center to ensure all attributes are set
					resource.TestCheckResourceAttr("data.leostream_centers.test", "centers.0.name", "Test AWS center"),
					resource.TestCheckResourceAttr("data.leostream_centers.test", "centers.0.flavor", "Z"),
					resource.TestCheckResourceAttr("data.leostream_centers.test", "centers.0.online", "1"),
					resource.TestCheckResourceAttr("data.leostream_centers.test", "centers.0.status_label", "Online"),
					resource.TestCheckResourceAttr("data.leostream_centers.test", "centers.0.center_type", "amazon"),
					resource.TestCheckResourceAttr("data.leostream_centers.test", "centers.0.type_label", "Amazon Web Services"),
				),
			},
		},
	})
}
