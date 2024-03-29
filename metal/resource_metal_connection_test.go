package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

const (
	metalDedicatedConnIDEnvVar = "TF_ACC_METAL_DEDICATED_CONNECTION_ID"
)

func TestSpeedConversion(t *testing.T) {
	speedUint, err := speedStrToUint("50Mbps")
	if err != nil {
		t.Errorf("Error converting speed string to uint64: %s", err)
	}
	if speedUint != 50*mega {
		t.Errorf("Speed string conversion failed. Expected: %d, got: %d", 50*mega, speedUint)
	}

	speedStr, err := speedUintToStr(50 * mega)
	if err != nil {
		t.Errorf("Error converting speed uint to string: %s", err)
	}
	if speedStr != "50Mbps" {
		t.Errorf("Speed uint conversion failed. Expected: %s, got: %s", "50Mbps", speedStr)
	}

	speedUint, err = speedStrToUint("100Gbps")
	if err == nil {
		t.Errorf("Expected error converting invalid speed string to uint, got: %d", speedUint)
	}

	speedStr, err = speedUintToStr(100 * giga)
	if err == nil {
		t.Errorf("Expected error converting invalid speed uint to string, got: %s", speedStr)
	}
}

func testAccCheckMetalConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_connection" {
			continue
		}
		if _, _, err := client.Connections.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}

func testAccMetalConnectionConfig_Shared(randstr string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
	name = "%[1]s-pro-conn-%[2]s"
}

resource "metal_connection" "test" {
	name               = "%[1]s-conn-%[2]s"
	project_id         = metal_project.test.id
	type               = "shared"
	redundancy         = "redundant"
	metro              = "sv"
	speed              = "50Mbps"
	service_token_type = "a_side"
}
`, tstResourcePrefix, randstr)
}

func TestAccMetalConnection_Shared(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_Shared(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr(
						"metal_connection.test", "service_tokens.0.type", "a_side"),
					resource.TestCheckResourceAttr(
						"metal_connection.test", "service_token_type", "a_side"),
					resource.TestCheckResourceAttr(
						"metal_connection.test", "service_tokens.0.max_allowed_speed", "50Mbps"),
				),
			},
			{
				ResourceName:      "metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalConnectionConfig_Shared(rs) + testDataSourceMetalConnectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr("data.metal_connection.test", "mode", "standard"),
					resource.TestCheckResourceAttr("data.metal_connection.test", "type", "shared"),
					resource.TestCheckResourceAttr("data.metal_connection.test", "redundancy", "redundant"),
					resource.TestCheckResourceAttr(
						"data.metal_connection.test", "service_tokens.0.type", "a_side"),
					resource.TestCheckResourceAttr(
						"data.metal_connection.test", "service_token_type", "a_side"),
					resource.TestCheckResourceAttr(
						"data.metal_connection.test", "service_tokens.0.max_allowed_speed", "50Mbps"),
				),
			},
		},
	})
}

func testAccMetalConnectionConfig_Dedicated(randstr string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
	name = "%[1]s-pro-conn-%[2]s"
}

// We use the project resource to get organization_id
resource "metal_connection" "test" {
	name            = "%[1]s-conn-%[2]s"
	metro           = "sv"
	organization_id = metal_project.test.organization_id
	type            = "dedicated"
	redundancy      = "redundant"
	tags            = ["%[1]s"]
	speed           = "50Mbps"
	mode            = "standard"
}
`, tstResourcePrefix, randstr)
}

func testDataSourceMetalConnectionConfig() string {
	return `
data "metal_connection" "test" {
	connection_id = metal_connection.test.id
}
`
}

func TestAccMetalConnection_Dedicated(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_Dedicated(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr("metal_connection.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("metal_connection.test", "mode", "standard"),
					resource.TestCheckResourceAttr("metal_connection.test", "type", "dedicated"),
					resource.TestCheckResourceAttr("metal_connection.test", "redundancy", "redundant"),
				),
			},
			{
				ResourceName:      "metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalConnectionConfig_Dedicated(rs) + testDataSourceMetalConnectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr("data.metal_connection.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("data.metal_connection.test", "mode", "standard"),
					resource.TestCheckResourceAttr("data.metal_connection.test", "type", "dedicated"),
					resource.TestCheckResourceAttr("data.metal_connection.test", "redundancy", "redundant"),
				),
			},
		},
	})
}

func testAccMetalConnectionConfig_Tunnel(randstr string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
	name = "%[1]s-pro-conn-%[2]s"
}

// We use the project resource to get organization_id internally
resource "metal_connection" "test" {
	name            = "%[1]s-conn-%[2]s"
	project_id      = metal_project.test.id
	metro           = "sv"
	redundancy      = "redundant"
	type            = "dedicated"
	mode            = "tunnel"
	speed           = "50Mbps"
}
`, tstResourcePrefix, randstr)
}

func TestAccMetalConnection_Tunnel(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_Tunnel(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_connection.test", "mode", "tunnel"),
				),
			},
			{
				ResourceName:            "metal_connection.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}
