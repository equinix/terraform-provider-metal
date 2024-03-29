package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetalDatasourcePreCreatedIPBlock_Basic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDatasourcePreCreatedIPBlockConfig_Basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.metal_precreated_ip_block.test", "cidr_notation"),
					resource.TestCheckResourceAttrPair(
						"metal_ip_attachment.test", "device_id",
						"metal_device.test", "id"),
				),
			},
			{
				ResourceName:      "metal_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDatasourcePreCreatedIPBlockConfig_Basic(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-recreated_ip_block-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test-ip-blockt"
  plan             = local.plan
  metro            = local.metro
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = metal_project.test.id

  lifecycle {
    ignore_changes = [
      plan,
      metro,
    ]
  }
}

data "metal_precreated_ip_block" "test" {
    facility         = metal_device.test.deployed_facility
    project_id       = metal_device.test.project_id
    address_family   = 6
    public           = true
}

resource "metal_ip_attachment" "test" {
    device_id = metal_device.test.id
    cidr_notation = cidrsubnet(data.metal_precreated_ip_block.test.cidr_notation,8,2)
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros), name)
}
