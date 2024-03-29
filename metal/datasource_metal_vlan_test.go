package metal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckMetalDatasourceVlanConfig_ByVxlanFacility(projSuffix, fac string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "%[1]s-pro-vlan-%[2]s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    facility = "%[3]s"
    description = "%[1]s-vlan foovlan"
}

data "metal_vlan" "dsvlan" {
    facility = metal_vlan.foovlan.facility
    project_id = metal_vlan.foovlan.project_id
    vxlan = metal_vlan.foovlan.vxlan
}
`, tstResourcePrefix, projSuffix, fac)
}

func TestAccMetalDatasourceVlan_ByVxlanFacility(t *testing.T) {
	rs := acctest.RandString(10)
	fac := "sv15"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDatasourceVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDatasourceVlanConfig_ByVxlanFacility(rs, fac),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "vxlan",
						"data.metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "id",
						"data.metal_vlan.dsvlan", "id",
					),
				),
			},
		},
	})
}

func testAccCheckMetalDatasourceVlanConfig_ByVxlanMetro(projSuffix, metro string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "%[1]s-pro-vlan-%[2]s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    metro = "%[3]s"
    description = "%[1]s-vlan foovlan"
    vxlan = 5
}

data "metal_vlan" "dsvlan" {
    metro = metal_vlan.foovlan.metro
    project_id = metal_vlan.foovlan.project_id
    vxlan = metal_vlan.foovlan.vxlan
}

resource "metal_vlan" "barvlan" {
    project_id = metal_project.foobar.id
    metro = metal_vlan.foovlan.metro
	description = "%[1]s-vlan barvlan"
    vxlan = 6
}

data "metal_vlan" "bardsvlan" {
    metro = metal_vlan.barvlan.metro
    project_id = metal_vlan.barvlan.project_id
    vxlan = metal_vlan.barvlan.vxlan
}
`, tstResourcePrefix, projSuffix, metro)
}

func TestAccMetalDatasourceVlan_ByVxlanMetro(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDatasourceVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDatasourceVlanConfig_ByVxlanMetro(rs, metro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "vxlan",
						"data.metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "id",
						"data.metal_vlan.dsvlan", "id",
					),
					resource.TestCheckResourceAttr(
						"metal_vlan.barvlan", "vxlan", "6",
					),
					resource.TestCheckResourceAttr(
						"data.metal_vlan.bardsvlan", "vxlan", "6",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.barvlan", "id",
						"data.metal_vlan.bardsvlan", "id",
					),
				),
			},
		},
	})
}

func testAccCheckMetalDatasourceVlanDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_vlan" {
			continue
		}
		if _, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("DatasourceVlan still exists")
		}
	}

	return nil
}

func testAccCheckMetalDatasourceVlanConfig_ByVlanId(projSuffix, metro string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "%[1]s-pro-vlan-%[2]s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    metro = "%[3]s"
    description = "%[1]s-vlan foovlan"
    vxlan = 5
}

data "metal_vlan" "dsvlan" {
    vlan_id = metal_vlan.foovlan.id
}
`, tstResourcePrefix, projSuffix, metro)
}

func TestAccMetalDatasourceVlan_ByVlanId(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDatasourceVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDatasourceVlanConfig_ByVlanId(rs, metro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "vxlan",
						"data.metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "project_id",
						"data.metal_vlan.dsvlan", "project_id",
					),
				),
			},
		},
	})
}

func testAccCheckMetalDatasourceVlanConfig_ByProjectId(projSuffix, metro string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "%[1]s-pro-vlan-%[2]s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    metro = "%[3]s"
    description = "%[1]s-vlan foovlan"
    vxlan = 5
}

data "metal_vlan" "dsvlan" {
    project_id = metal_vlan.foovlan.project_id
}
`, tstResourcePrefix, projSuffix, metro)
}

func TestAccMetalDatasourceVlan_ByProjectId(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDatasourceVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDatasourceVlanConfig_ByProjectId(rs, metro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "vxlan",
						"data.metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "project_id",
						"data.metal_vlan.dsvlan", "project_id",
					),
				),
			},
		},
	})
}

func Test_matchingVlan(t *testing.T) {
	type args struct {
		vlans     []packngo.VirtualNetwork
		vxlan     int
		projectID string
		facility  string
		metro     string
	}
	tests := []struct {
		name    string
		args    args
		want    *packngo.VirtualNetwork
		wantErr bool
	}{{
		name: "MatchingVLAN",
		args: args{
			vlans:     []packngo.VirtualNetwork{{VXLAN: 123}},
			vxlan:     123,
			projectID: "",
			facility:  "",
			metro:     "",
		},
		want:    &packngo.VirtualNetwork{VXLAN: 123},
		wantErr: false,
	},
		{
			name: "MatchingFac",
			args: args{
				vlans:    []packngo.VirtualNetwork{{FacilityCode: "fac"}},
				facility: "fac",
			},
			want:    &packngo.VirtualNetwork{FacilityCode: "fac"},
			wantErr: false,
		},
		{
			name: "MatchingMet",
			args: args{
				vlans: []packngo.VirtualNetwork{{MetroCode: "met"}},
				metro: "met",
			},
			want:    &packngo.VirtualNetwork{MetroCode: "met"},
			wantErr: false,
		},
		{
			name: "SecondMatch",
			args: args{
				vlans: []packngo.VirtualNetwork{{FacilityCode: "fac"}, {MetroCode: "met"}},
				metro: "met",
			},
			want:    &packngo.VirtualNetwork{MetroCode: "met"},
			wantErr: false,
		},
		{
			name: "TwoMatches",
			args: args{
				vlans: []packngo.VirtualNetwork{{MetroCode: "met"}, {MetroCode: "met"}},
				metro: "met",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ComplexMatch",
			args: args{
				vlans: []packngo.VirtualNetwork{{VXLAN: 987, FacilityCode: "fac", MetroCode: "skip"}, {VXLAN: 123, FacilityCode: "fac", MetroCode: "met"}, {VXLAN: 456, FacilityCode: "fac", MetroCode: "nope"}},
				metro: "met",
			},
			want:    &packngo.VirtualNetwork{VXLAN: 123, FacilityCode: "fac", MetroCode: "met"},
			wantErr: false,
		},
		{
			name: "NoMatch",
			args: args{
				vlans:     nil,
				vxlan:     123,
				projectID: "pid",
				facility:  "fac",
				metro:     "met",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchingVlan(tt.args.vlans, tt.args.vxlan, tt.args.projectID, tt.args.facility, tt.args.metro)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchingVlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchingVlan() = %v, want %v", got, tt.want)
			}
		})
	}
}
