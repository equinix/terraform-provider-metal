package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMetalUserAPIKey() *schema.Resource {
	userKeySchema := schemaMetalAPIKey()
	userKeySchema["user_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "UUID of user owning this key",
	}
	return &schema.Resource{
		DeprecationMessage: deprecatedProviderMsg,
		Create:             resourceMetalAPIKeyCreate,
		Read:               resourceMetalAPIKeyRead,
		Delete:             resourceMetalAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: userKeySchema,
	}
}
