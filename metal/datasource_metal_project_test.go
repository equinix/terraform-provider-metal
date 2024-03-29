package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/packethost/packngo"
)

func testAccCheckMetalDataSourceProject_Basic(r string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
	name = "tfacc-pro-%s"
	bgp_config {
		deployment_type = "local"
		md5 = "2SFsdfsg43"
		asn = 65000
	}
}

data metal_project "test" {
	project_id = metal_project.foobar.id
}

`, r)
}

func TestAccMetalDataSourceProject_Basic(t *testing.T) {
	var project packngo.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDataSourceProject_Basic(rn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "name", fmt.Sprintf("tfacc-pro-%s", rn)),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "bgp_config.0.md5",
						"2SFsdfsg43"),
					resource.TestCheckResourceAttrPair(
						"metal_project.foobar", "id",
						"data.metal_project.test", "id"),
				),
			},
		},
	})
}
