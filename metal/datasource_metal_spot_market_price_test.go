package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalSpotPrice_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSpotMarketRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalSpotMarketPrice(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.metal_spot_market_price.metro", "price"),
					resource.TestCheckResourceAttrSet(
						"data.metal_spot_market_price.facility", "price"),
				),
			},
		},
	})
}

func testDataSourceMetalSpotMarketPrice() string {
	return fmt.Sprintf(`
data "metal_spot_market_price" "metro" {
	metro    = "da"
	plan     = "c3.medium.x86"
}

data "metal_spot_market_price" "facility" {
	facility = "da11"
	plan     = "c3.medium.x86"
}
`)
}
