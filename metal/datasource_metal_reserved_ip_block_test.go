package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceMetalReservedIPBlockConfig_Basic(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
	name = "tfacc-pro-reserved_ip_block-%s"
}

resource "metal_reserved_ip_block" "test" {
	project_id  = metal_project.foobar.id
	metro       = "sv"
	type        = "public_ipv4"
	quantity    = 2
}

data "metal_reserved_ip_block" "test" {
	project_id  = metal_project.foobar.id
    ip_address  = cidrhost(metal_reserved_ip_block.test.cidr_notation,1)
}

data "metal_reserved_ip_block" "test_id" {
	id  = metal_reserved_ip_block.test.id
}

`, name)
}

func TestAccDataSourceMetalReservedIPBlock_Basic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalReservedIPBlockConfig_Basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_reserved_ip_block.test", "id",
						"data.metal_reserved_ip_block.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"metal_reserved_ip_block.test", "cidr_notation",
						"data.metal_reserved_ip_block.test_id", "cidr_notation",
					),
				),
			},
		},
	})
}
