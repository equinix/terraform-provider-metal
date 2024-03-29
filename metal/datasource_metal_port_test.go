package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetalPort_ByName(t *testing.T) {

	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPortConfig_ByName(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_port.test", "bond_name", "bond0"),
				),
			},
		},
	})
}

func testPortConfig_ByName(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-port-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test-port"
  plan             = local.plan
  metro            = local.metro
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = metal_project.test.id

  lifecycle {
    ignore_changes = [
      plan,
      metro,
    ]
  }
}

data "metal_port" "test" {
    device_id = metal_device.test.id
    name      = "eth0"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), name)
}

func TestAccMetalPort_ById(t *testing.T) {

	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPortConfig_ById(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.metal_port.test", "name"),
				),
			},
		},
	})
}

func testPortConfig_ById(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-port-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test-port"
  plan             = local.plan
  metro            = local.metro
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = metal_project.test.id

  lifecycle {
    ignore_changes = [
      plan,
      metro,
    ]
  }
}

data "metal_port" "test" {
  port_id        = metal_device.test.ports[0].id
}

`, confAccMetalDevice_base(preferable_plans, preferable_metros), name)
}
