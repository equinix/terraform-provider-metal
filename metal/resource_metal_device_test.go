package metal

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

// list of plans and metros used as filter criteria to find available hardware to run tests
var preferable_plans = []string{"x1.small.x86", "t1.small.x86", "c2.small.x86", "c2.medium.x86", "c3.small.x86", "c3.medium.x86", "m3.small.x86", "m3.large.x86"}
var preferable_metros = []string{"ch", "ny", "sv", "ty", "am"}

func init() {
	resource.AddTestSweepers("metal_device", &resource.Sweeper{
		Name: "metal_device",
		F:    testSweepDevices,
	})
}

func testSweepDevices(region string) error {
	log.Printf("[DEBUG] Sweeping devices")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping devices: %s", err)
	}
	client := meta.(*packngo.Client)

	ps, _, err := client.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("Error getting project list for sweepeing devices: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if strings.HasPrefix(p.Name, "tfacc-") {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := client.Devices.List(pid, nil)
		if err != nil {
			return fmt.Errorf("Error listing devices to sweep: %s", err)
		}
		for _, d := range ds {
			dids = append(dids, d.ID)
		}
	}

	for _, did := range dids {
		log.Printf("Removing device %s", did)
		_, err := client.Devices.Delete(did, true)
		if err != nil {
			return fmt.Errorf("Error deleting device %s", err)
		}
	}
	return nil
}

// Regexp vars for use with resource.ExpectError
var matchErrMustBeProvided = regexp.MustCompile(".* must be provided when .*")
var matchErrShouldNotBeAnIPXE = regexp.MustCompile(`.*"user_data" should not be an iPXE.*`)

// This function should be used to find available plans in all test where a metal_device resource is needed.
// To prevent unexpected plan/facilities changes (i.e. run out of a plan in a metro after first apply)
// during tests that have several config updates, resource metal_device should include a lifecycle
// like the one defined below.
//
// lifecycle {
//     ignore_changes = [
//       plan,
//       facilities,
//     ]
//   }
func confAccMetalDevice_base(plans, metros []string) string {
	return fmt.Sprintf(`
data "metal_plans" "test" {
    sort {
        attribute = "pricing_hour"
        direction = "asc"
    }
    filter {
        attribute = "name"
        values    = [%s]
    }
    filter {
        attribute = "available_in_metros"
        values    = [%s]
    }
}

locals {
    plan       = data.metal_plans.test.plans[0].slug
    metro      = tolist(data.metal_plans.test.plans[0].available_in_metros)[1]
    facilities = tolist(setsubtract(data.metal_plans.test.plans[0].available_in, ["sjc1", "ld7", "sy4"]))
}
`, fmt.Sprintf("\"%s\"", strings.Join(plans[:], `","`)), fmt.Sprintf("\"%s\"", strings.Join(metros[:], `","`)))
}

func TestAccMetalDevice_FacilityList(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDeviceConfig_facility_list(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &device),
				),
			},
		},
	})
}

func TestAccMetalDevice_NetworkPortsOrder(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDeviceConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &device),
					testAccCheckMetalDeviceNetworkOrder(r),
					testAccCheckMetalDevicePortsOrder(r),
				),
			},
		},
	})
}

func TestAccMetalDevice_Basic(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDeviceConfig_minimal(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &device),
					testAccCheckMetalDeviceNetwork(r),
					resource.TestCheckResourceAttrSet(
						r, "hostname"),
					resource.TestCheckResourceAttr(
						r, "billing_cycle", "hourly"),
					resource.TestCheckResourceAttr(
						r, "network_type", "layer3"),
					resource.TestCheckResourceAttr(
						r, "ipxe_script_url", ""),
					resource.TestCheckResourceAttr(
						r, "always_pxe", "false"),
					resource.TestCheckResourceAttrSet(
						r, "root_password"),
					resource.TestCheckResourceAttrPair(
						r, "deployed_facility", r, "facilities.0"),
				),
			},
			{
				Config: testAccCheckMetalDeviceConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &device),
					testAccCheckMetalDeviceNetwork(r),
					testAccCheckMetalDeviceAttributes(&device),
				),
			},
		},
	})
}

func TestAccMetalDevice_Metro(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDeviceConfig_metro(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &device),
					testAccCheckMetalDeviceNetwork(r),
					testAccCheckMetalDeviceAttributes(&device),
					resource.TestCheckResourceAttr(
						r, "metro", "da"),
				),
			},
		},
	})
}
func TestAccMetalDevice_Update(t *testing.T) {
	var d1, d2, d3, d4, d5 packngo.Device
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	r := "metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDeviceConfig_varname(rInt, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &d1),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-device-test-%d", rInt)),
				),
			},
			{
				Config: testAccCheckMetalDeviceConfig_varname(rInt+1, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &d2),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-device-test-%d", rInt+1)),
					testAccCheckMetalSameDevice(t, &d1, &d2),
				),
			},
			{
				Config: testAccCheckMetalDeviceConfig_varname(rInt+2, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &d3),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-device-test-%d", rInt+2)),
					resource.TestCheckResourceAttr(r, "description", fmt.Sprintf("test-desc-%d", rInt+2)),
					resource.TestCheckResourceAttr(r, "tags.0", fmt.Sprintf("%d", rInt+2)),
					testAccCheckMetalSameDevice(t, &d2, &d3),
				),
			},
			{
				Config: testAccCheckMetalDeviceConfig_no_description(rInt+3, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &d4),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-device-test-%d", rInt+3)),
					resource.TestCheckResourceAttr(r, "tags.0", fmt.Sprintf("%d", rInt+3)),
					testAccCheckMetalSameDevice(t, &d3, &d4),
				),
			},
			{
				Config: testAccCheckMetalDeviceConfig_reinstall(rInt+4, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &d5),
					testAccCheckMetalSameDevice(t, &d4, &d5),
				),
			},
		},
	})
}

func TestAccMetalDevice_IPXEScriptUrl(t *testing.T) {
	var device, d2 packngo.Device
	rs := acctest.RandString(10)
	r := "metal_device.test_ipxe_script_url"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDeviceConfig_ipxe_script_url(rs, "https://boot.netboot.xyz", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &device),
					testAccCheckMetalDeviceNetwork(r),
					resource.TestCheckResourceAttr(
						r, "ipxe_script_url", "https://boot.netboot.xyz"),
					resource.TestCheckResourceAttr(
						r, "always_pxe", "true"),
				),
			},
			{
				Config: testAccCheckMetalDeviceConfig_ipxe_script_url(rs, "https://new.netboot.xyz", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &d2),
					testAccCheckMetalDeviceNetwork(r),
					resource.TestCheckResourceAttr(
						r, "ipxe_script_url", "https://new.netboot.xyz"),
					resource.TestCheckResourceAttr(
						r, "always_pxe", "false"),
					testAccCheckMetalSameDevice(t, &device, &d2),
				),
			},
		},
	})
}

func TestAccMetalDevice_IPXEConflictingFields(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "metal_device.test_ipxe_conflict"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckMetalDeviceConfig_ipxe_conflict, confAccMetalDevice_base(preferable_plans, preferable_metros), rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &device),
				),
				ExpectError: matchErrShouldNotBeAnIPXE,
			},
		},
	})
}

func TestAccMetalDevice_IPXEConfigMissing(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "metal_device.test_ipxe_config_missing"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckMetalDeviceConfig_ipxe_missing, confAccMetalDevice_base(preferable_plans, preferable_metros), rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalDeviceExists(r, &device),
				),
				ExpectError: matchErrMustBeProvided,
			},
		},
	})
}

func testAccCheckMetalDeviceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_device" {
			continue
		}
		if _, _, err := client.Devices.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Device still exists")
		}
	}
	return nil
}

func testAccCheckMetalDeviceAttributes(device *packngo.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if device.Hostname != "tfacc-device-test" {
			return fmt.Errorf("Bad name: %s", device.Hostname)
		}
		if device.State != "active" {
			return fmt.Errorf("Device should be 'active', not '%s'", device.State)
		}

		return nil
	}
}

func testAccCheckMetalDeviceExists(n string, device *packngo.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundDevice, _, err := client.Devices.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundDevice.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundDevice)
		}

		*device = *foundDevice

		return nil
	}
}

func testAccCheckMetalSameDevice(t *testing.T, before, after *packngo.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID != after.ID {
			t.Fatalf("Expected device to be the same, but it was recreated: %s -> %s", before.ID, after.ID)
		}
		return nil
	}
}

func testAccCheckMetalDevicePortsOrder(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.Attributes["ports.0.name"] != "bond0" {
			return fmt.Errorf("first port should be bond0")
		}
		if rs.Primary.Attributes["ports.1.name"] != "eth0" {
			return fmt.Errorf("second port should be eth0")
		}
		if rs.Primary.Attributes["ports.2.name"] != "eth1" {
			return fmt.Errorf("third port should be eth1")
		}
		return nil
	}
}

func testAccCheckMetalDeviceNetworkOrder(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.Attributes["network.0.family"] != "4" {
			return fmt.Errorf("first netowrk should be public IPv4")
		}
		if rs.Primary.Attributes["network.0.public"] != "true" {
			return fmt.Errorf("first netowrk should be public IPv4")
		}
		if rs.Primary.Attributes["network.1.family"] != "6" {
			return fmt.Errorf("second netowrk should be public IPv6")
		}
		if rs.Primary.Attributes["network.2.family"] != "4" {
			return fmt.Errorf("third netowrk should be private IPv4")
		}
		if rs.Primary.Attributes["network.2.public"] == "true" {
			return fmt.Errorf("third netowrk should be private IPv4")
		}
		return nil
	}
}

func testAccCheckMetalDeviceNetwork(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var ip net.IP
		var k, v string
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		k = "access_public_ipv6"
		v = rs.Primary.Attributes[k]
		ip = net.ParseIP(v)
		if ip == nil {
			return fmt.Errorf("\"%s\" is not a valid IP address: %s",
				k, v)
		}

		k = "access_public_ipv4"
		v = rs.Primary.Attributes[k]
		ip = net.ParseIP(v)
		if ip == nil {
			return fmt.Errorf("\"%s\" is not a valid IP address: %s",
				k, v)
		}

		k = "access_private_ipv4"
		v = rs.Primary.Attributes[k]
		ip = net.ParseIP(v)
		if ip == nil {
			return fmt.Errorf("\"%s\" is not a valid IP address: %s",
				k, v)
		}

		return nil
	}
}

func TestAccMetalDevice_importBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDeviceConfig_basic(rs),
			},
			{
				ResourceName:      "metal_device.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMetalDeviceConfig_no_description(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-device-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test-%d"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
  tags             = ["%d"]

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix, rInt, rInt)
}

func testAccCheckMetalDeviceConfig_reinstall(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-device-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test-%d"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
  tags             = ["%d"]
  user_data = "#!/usr/bin/env sh\necho Reinstall\n"

  reinstall {
	  enabled = true
	  deprovision_fast = true
  }

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix, rInt, rInt)
}

func testAccCheckMetalDeviceConfig_varname(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-device-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test-%d"
  description      = "test-desc-%d"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
  tags             = ["%d"]

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix, rInt, rInt, rInt)
}

func testAccCheckMetalDeviceConfig_varname_pxe(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-device-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test-%d"
  description      = "test-desc-%d"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
  tags             = ["%d"]
  always_pxe       = true
  ipxe_script_url  = "http://matchbox.foo.wtf:8080/boot.ipxe"

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix, rInt, rInt, rInt)
}

func testAccCheckMetalDeviceConfig_metro(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-device-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"

  lifecycle {
    ignore_changes = [
      plan,
      metro,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix)
}

func testAccCheckMetalDeviceConfig_minimal(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-device-%s"
}

resource "metal_device" "test" {
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "ubuntu_16_04"
  project_id       = "${metal_project.test.id}"

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix)
}

func testAccCheckMetalDeviceConfig_basic(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
    name = "tfacc-pro-device-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-test"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix)
}

func testAccCheckMetalDeviceConfig_facility_list(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
  name = "tfacc-pro-device-%s"
}

resource "metal_device" "test"  {
  hostname         = "tfacc-device-test-ipxe-script-url"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix)
}

func testAccCheckMetalDeviceConfig_ipxe_script_url(projSuffix, url, pxe string) string {
	return fmt.Sprintf(`
%s

resource "metal_project" "test" {
  name = "tfacc-pro-device-%s"
}

resource "metal_device" "test_ipxe_script_url"  {
  hostname         = "tfacc-device-test-ipxe-script-url"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "custom_ipxe"
  user_data        = "#!/bin/sh\ntouch /tmp/test"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
  ipxe_script_url  = "%s"
  always_pxe       = "%s"

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros), projSuffix, url, pxe)
}

var testAccCheckMetalDeviceConfig_ipxe_conflict = `
%s

resource "metal_project" "test" {
  name = "tfacc-pro-device-%s"
}

resource "metal_device" "test_ipxe_conflict" {
  hostname         = "tfacc-device-test-ipxe-conflict"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "custom_ipxe"
  user_data        = "#!ipxe\nset conflict ipxe_script_url"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
  ipxe_script_url  = "https://boot.netboot.xyz"
  always_pxe       = true

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`

var testAccCheckMetalDeviceConfig_ipxe_missing = `
%s

resource "metal_project" "test" {
  name = "tfacc-pro-device-%s"
}

resource "metal_device" "test_ipxe_missing" {
  hostname         = "tfacc-device-test-ipxe-missing"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = "custom_ipxe"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
  always_pxe       = true

  lifecycle {
    ignore_changes = [
      plan,
      facilities,
    ]
  }
}`
