package metal

import (
	"fmt"
	"testing"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckMetalVirtualCircuitDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_virtual_circuit" {
			continue
		}
		if _, _, err := client.VirtualCircuits.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("VirtualCircuit still exists")
		}
	}

	return nil
}

func testAccMetalVirtualCircuitConfig_Dedicated(randstr string, randint int) string {
	// Dedicated connection in DA metro
	testConnection := os.Getenv(metalDedicatedConnIDEnvVar)
	return fmt.Sprintf(`
        data "metal_connection" "test" {
			connection_id = "%s"
		}

		resource "metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        resource "metal_vlan" "test" {
            project_id = metal_project.test.id
            metro      = "da"
        }

        resource "metal_virtual_circuit" "test" {
            connection_id = data.metal_connection.test.id
            project_id = metal_project.test.id
            port_id = data.metal_connection.test.ports[0].id
            vlan_id = metal_vlan.test.id
            nni_vlan = %d
        }`,
		testConnection, randstr, randint)
}

func TestAccMetalVirtualCircuit_Dedicated(t *testing.T) {

	rs := acctest.RandString(10)
	ri := acctest.RandIntRange(1024, 1093)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVirtualCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVirtualCircuitConfig_Dedicated(rs, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_virtual_circuit.test", "vlan_id",
						"metal_vlan.test", "id",
					),
				),
			},
			{
				ResourceName:      "metal_virtual_circuit.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
