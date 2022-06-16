package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccOrgDataSource_Basic(t *testing.T) {
	var org packngo.Organization
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalOrgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalOrgDataSourceConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalOrgExists("metal_organization.test", &org),
					resource.TestCheckResourceAttr(
						"metal_organization.test", "name",
						fmt.Sprintf("tfacc-org-datasource-%d", rInt)),
					resource.TestCheckResourceAttr(
						"metal_organization.test", "description", "quux"),
					resource.TestCheckResourceAttr(
						"data.metal_organization.test", "name",
						fmt.Sprintf("tfacc-org-datasource-%d", rInt)),
					resource.TestCheckResourceAttrPair(
						"metal_organization.test", "address.0.address",
						"data.metal_organization.test", "address.0.address",
					),
					resource.TestCheckResourceAttrPair(
						"metal_organization.test", "address.0.city",
						"data.metal_organization.test", "address.0.city",
					),
					resource.TestCheckResourceAttrPair(
						"metal_organization.test", "address.0.country",
						"data.metal_organization.test", "address.0.country",
					),
					resource.TestCheckResourceAttrPair(
						"metal_organization.test", "address.0.zip_code",
						"data.metal_organization.test", "address.0.zip_code",
					),
				),
			},
		},
	})
}

func testAccCheckMetalOrgDataSourceConfigBasic(r int) string {
	return fmt.Sprintf(`
resource "metal_organization" "test" {
  name = "tfacc-org-datasource-%d"
  description = "quux"
  address {
	address = "tfacc org street"
	city = "london"
	zip_code = "12345"
	country = "GB"
  }
}

data "metal_organization" "test" {
  organization_id = metal_organization.test.id
}

`, r)
}
