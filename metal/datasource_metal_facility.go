package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalFacility() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalFacilityRead,
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the Facility to match",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of this Facility.",
				Optional:    true,
				Computed:    true,
			},
			"features": {
				Type:        schema.TypeList,
				Description: "The features of this Facility.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceMetalFacilityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	code := d.Get("code").(string)

	if code == "" {
		return fmt.Errorf("Error Facility code is required")
	}

	facilities, _, err := client.Facilities.List(nil)
	if err != nil {
		return fmt.Errorf("Error listing Facilities: %s", err)
	}

	for _, f := range facilities {
		if f.Code == code {
			d.SetId(f.ID)
			return setMap(d, map[string]interface{}{
				"code":     f.Code,
				"name":     f.Name,
				"features": f.Features,
			})
		}
	}

	return fmt.Errorf("Facility %s was not found", code)
}
