package metal

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("metal_vlan", &resource.Sweeper{
		Name:         "metal_vlan",
		Dependencies: []string{"metal_virtual_circuit", "metal_vrf", "metal_device"},
		F:            testSweepVlans,
	})
}

func testSweepVlans(region string) error {
	log.Printf("[DEBUG] Sweeping vlans")
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping vlans: %s", err)
	}
	metal := config.(*packngo.Client)
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping vlans: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if isSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := metal.ProjectVirtualNetworks.List(pid, nil)
		if err != nil {
			log.Printf("Error listing vlans to sweep: %s", err)
			continue
		}
		for _, d := range ds.VirtualNetworks {
			if isSweepableTestResource(d.Description) {
				dids = append(dids, d.ID)
			}
		}
	}

	for _, did := range dids {
		log.Printf("Removing vlan %s", did)
		_, err := metal.ProjectVirtualNetworks.Delete(did)
		if err != nil {
			return fmt.Errorf("Error deleting vlan %s", err)
		}
	}
	return nil
}

func testAccCheckMetalVlanConfig_metro(projSuffix, metro string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
	name = "%[1]s-pro-%[2]s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    metro = "%[3]s"
    description = "%[1]s-vlan foovlan"
    vxlan = 5
}
`, tstResourcePrefix, projSuffix, metro)
}

func TestAccMetalVlan_Metro(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVlanConfig_metro(rs, metro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_vlan.foovlan", "metro", metro),
					resource.TestCheckResourceAttr(
						"metal_vlan.foovlan", "facility", ""),
				),
			},
		},
	})
}

func TestAccMetalVlan_Basic(t *testing.T) {
	var vlan packngo.VirtualNetwork
	rs := acctest.RandString(10)
	fac := "ny5"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVlanConfig_var(rs, fac),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalVlanExists("metal_vlan.foovlan", &vlan),
					resource.TestCheckResourceAttr(
						"metal_vlan.foovlan", "description", "tfacc-vlan foovlan"),
					resource.TestCheckResourceAttr(
						"metal_vlan.foovlan", "facility", fac),
				),
			},
		},
	})
}

func testAccCheckMetalVlanExists(n string, vlan *packngo.VirtualNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundVlan, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundVlan.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundVlan)
		}

		*vlan = *foundVlan

		return nil
	}
}

func testAccCheckMetalVlanDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_vlan" {
			continue
		}
		if _, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Vlan still exists")
		}
	}

	return nil
}

func testAccCheckMetalVlanConfig_var(projSuffix, fac string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
	name = "%[1]s-pro-%[2]s"
}

resource "metal_vlan" "foovlan" {
    project_id  = "${metal_project.foobar.id}"
    facility    = "%[3]s"
    description = "%[1]s-vlan foovlan"
}
`, tstResourcePrefix, projSuffix, fac)
}

func TestAccMetalVlan_importBasic(t *testing.T) {
	rs := acctest.RandString(10)
	fac := "ny5"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVlanConfig_var(rs, fac),
			},
			{
				ResourceName:      "metal_vlan.foovlan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
