package metal

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("metal_vrf", &resource.Sweeper{
		Name: "metal_vrf",
		Dependencies: []string{
			"metal_device",
			"metal_virtual_circuit",
			// TODO: add sweeper when offered
			// "metal_reserved_ip_block",
		},
		F: testSweepVRFs,
	})
}

func testSweepVRFs(region string) error {
	log.Printf("[DEBUG] Sweeping VRFs")
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping VRFs: %s", err)
	}
	metal := config.(*packngo.Client)
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping VRFs: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if isSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := metal.VRFs.List(pid, nil)
		if err != nil {
			log.Printf("Error listing VRFs to sweep: %s", err)
			continue
		}
		for _, d := range ds {
			if isSweepableTestResource(d.Name) {
				dids = append(dids, d.ID)
			}
		}
	}

	for _, did := range dids {
		log.Printf("Removing VRFs %s", did)
		_, err := metal.VRFs.Delete(did)
		if err != nil {
			return fmt.Errorf("Error deleting VRFs %s", err)
		}
	}
	return nil
}

func TestAccMetalVRF_basic(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrSet(
						"metal_vrf.test", "local_asn"),
				),
			},
			{
				ResourceName:      "metal_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalVRF_withIPRanges(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withIPRanges(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				ResourceName:      "metal_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
		},
	})
}

func TestAccMetalVRF_withIPReservations(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_withIPRanges(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withIPReservations(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrPair("metal_vrf.test", "id", "metal_reserved_ip_block.test", "vrf_id"),
				),
			},
			{
				ResourceName:      "metal_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:            "metal_reserved_ip_block.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_state"},
			},
		},
	})
}

func TestAccMetalVRF_withGateway(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_withIPReservations(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withGateway(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrPair("metal_vrf.test", "id", "metal_gateway.test", "vrf_id"),
				),
			},
			{
				ResourceName:      "metal_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "metal_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalVRFConfig_withConnection(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()
	nniVlan := acctest.RandIntRange(2024, 2093)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_withVC(rInt, nniVlan),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"metal_virtual_circuit.test", "name", fmt.Sprintf("tfacc-vc-%d", rInt)),
					resource.TestCheckResourceAttr(
						"metal_virtual_circuit.test",
						"nni_vlan", strconv.Itoa(nniVlan)),
					resource.TestCheckResourceAttrPair(
						"metal_virtual_circuit.test",
						"vrf_id", "metal_vrf.test", "id"),
					resource.TestCheckNoResourceAttr(
						"metal_virtual_circuit.test",
						"vlan_id"),
					resource.TestCheckResourceAttr(
						"metal_virtual_circuit.test",
						"peer_asn", "65530"),
					resource.TestCheckResourceAttr(
						"metal_virtual_circuit.test",
						"subnet", "192.168.100.16/31"),
					resource.TestCheckResourceAttr(
						"metal_virtual_circuit.test",
						"metal_ip", "192.168.100.16"),
					resource.TestCheckResourceAttr(
						"metal_virtual_circuit.test",
						"customer_ip", "192.168.100.17"),
				),
			},
			{
				ResourceName:      "metal_virtual_circuit.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalVRFConfig_withVCGateway(rInt, nniVlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"metal_virtual_circuit.test",
						"vlan_id"),
				),
			},
			{
				ResourceName:      "metal_reserved_ip_block.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "metal_vlan.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "metal_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalVRFCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_vrf" {
			continue
		}
		if _, _, err := client.VRFs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal VRF still exists")
		}
	}

	return nil
}

func testAccMetalVRFExists(n string, vrf *packngo.VRF) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundResource, _, err := client.VRFs.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundResource.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundResource)
		}

		*vrf = *foundResource

		return nil
	}
}

func testAccMetalVRFConfig_basic(r int) string {
	testMetro := "da"

	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-pro-vrf-%d"
}

resource "metal_vrf" "test" {
	name = "tfacc-vrf-%d"
	metro = "%s"
	project_id = "${metal_project.test.id}"
}`, r, r, testMetro)
}

func testAccMetalVRFConfig_withIPRanges(r int) string {
	testMetro := "da"

	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-pro-vrf-%d"
}

resource "metal_vrf" "test" {
	name = "tfacc-vrf-%d"
	metro = "%s"
	description = "tfacc-vrf-%d"
	local_asn = "65000"
	ip_ranges = ["192.168.100.0/25"]
	project_id = metal_project.test.id
}`, r, r, testMetro, r)
}

func testAccMetalVRFConfig_withIPReservations(r int) string {
	testMetro := "da"

	return testAccMetalVRFConfig_withIPRanges(r) + fmt.Sprintf(`

resource "metal_reserved_ip_block" "test" {
	vrf_id = metal_vrf.test.id
	cidr = 29
	description = "tfacc-reserved-ip-block-%d"
	network = "192.168.100.0"
	type = "vrf"
	metro = "%s"
	project_id = metal_project.test.id
}
`, r, testMetro)
}

func testAccMetalVRFConfig_withGateway(r int) string {
	testMetro := "da"

	return testAccMetalVRFConfig_withIPReservations(r) + fmt.Sprintf(`

resource "metal_vlan" "test" {
	description = "%s-vlan vrf"
	metro       = "%s"
	project_id  = metal_project.test.id
}

resource "metal_gateway" "test" {
    project_id        = metal_project.test.id
    vlan_id           = metal_vlan.test.id
    ip_reservation_id = metal_reserved_ip_block.test.id
}
`, tstResourcePrefix, testMetro)
}

func testAccMetalVRFConfig_withVC(r, nniVlan int) string {
	// Dedicated connection in DA metro
	testConnection := os.Getenv(metalDedicatedConnIDEnvVar)
	return testAccMetalVRFConfig_withIPRanges(r) + fmt.Sprintf(`

	data "metal_connection" "test" {
		connection_id = "%[1]s"
	}

	resource "metal_virtual_circuit" "test" {
		name = "%[4]s-vc-%[2]d"
		description = "%[4]s-vc-%[2]d"
		connection_id = data.metal_connection.test.id
		project_id = metal_project.test.id
		port_id = data.metal_connection.test.ports[0].id
		nni_vlan = %[3]d
		vrf_id = metal_vrf.test.id
		peer_asn = 65530
		subnet = "192.168.100.16/31"
		metal_ip = "192.168.100.16"
		customer_ip = "192.168.100.17"
	}
	`, testConnection, r, nniVlan, tstResourcePrefix)
}

func testAccMetalVRFConfig_withVCGateway(r, nniVlan int) string {
	// Dedicated connection in DA metro
	testConnection := os.Getenv(metalDedicatedConnIDEnvVar)
	return testAccMetalVRFConfig_withGateway(r) + fmt.Sprintf(`
	data "metal_connection" "test" {
		connection_id = "%[1]s"
	}

	resource "metal_virtual_circuit" "test" {
		name = "%[4]s-vc-%[2]d"
		description = "%[4]s-vc-%[2]d"
		connection_id = data.metal_connection.test.id
		project_id = metal_project.test.id
		port_id = data.metal_connection.test.ports[0].id
		nni_vlan = %[3]d
		vrf_id = metal_vrf.test.id
		peer_asn = 65530
		subnet = "192.168.100.16/31"
		metal_ip = "192.168.100.16"
		customer_ip = "192.168.100.17"
	}`, testConnection, r, nniVlan, tstResourcePrefix)
}

func testAccMetalVRFConfig_withConnection(r int) string {
	testMetro := "da"
	return testAccMetalVRFConfig_basic(r) + fmt.Sprintf(`

	resource "metal_connection" "test" {
		name            = "%s-conn-%d"
		organization_id = metal_project.foobar.organization_id
		project_id      = metal_project.foobar.id
		metro           = "%s"
		redundancy      = "redundant"
		type            = "shared"
	}`, tstResourcePrefix, r, testMetro)
}
