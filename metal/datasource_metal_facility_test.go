package metal

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	matchErrMissingFeature = regexp.MustCompile(`.*doesn't have feature.*`)
	matchErrNoCapacity     = regexp.MustCompile(`Not enough capacity.*`)
)

func TestAccDataSourceFacility_Basic(t *testing.T) {
	testFac := "da11"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFacilityConfigBasic(testFac),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_facility.test", "code", testFac),
				),
			},
			{
				Config: testAccDataSourceFacilityConfigCapacityReasonable(testFac),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_facility.test", "code", testFac),
				),
			},
			{
				Config:      testAccDataSourceFacilityConfigCapacityUnreasonable(testFac),
				ExpectError: matchErrNoCapacity,
			},
			{
				Config:      testAccDataSourceFacilityConfigCapacityUnreasonableMultiple(testFac),
				ExpectError: matchErrNoCapacity,
			},
		},
	})
}

func TestAccDataSourceFacility_Features(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceFacilityConfigFeatures(),
				ExpectError: matchErrMissingFeature,
			},
		},
	})
}

func testAccDataSourceFacilityConfigFeatures() string {
	return `
data "metal_facility" "test" {
    code = "da11"
    features_required = ["baremetal", "ibx", "missingFeature"]
}
`
}

func testAccDataSourceFacilityConfigBasic(facCode string) string {
	return fmt.Sprintf(`
data "metal_facility" "test" {
    code = "%s"
}
`, facCode)
}

func testAccDataSourceFacilityConfigCapacityUnreasonable(facCode string) string {
	return fmt.Sprintf(`
data "metal_facility" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1000
    }
}
`, facCode)
}

func testAccDataSourceFacilityConfigCapacityReasonable(facCode string) string {
	return fmt.Sprintf(`
data "metal_facility" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1
    }
    capacity {
        plan = "c3.medium.x86"
        quantity = 1
    }
}
`, facCode)
}

func testAccDataSourceFacilityConfigCapacityUnreasonableMultiple(facCode string) string {
	return fmt.Sprintf(`
data "metal_facility" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1
    }
    capacity {
        plan = "c3.medium.x86"
        quantity = 1000
    }
}
`, facCode)
}
