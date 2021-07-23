// Licenses under the MPL 2.0
//
// Source: https://github.com/hashicorp/terraform-plugin-sdk/pull/508
package metal

import (
	"fmt"

	"github.com/rs/xid"
)

// UniqueIDSuffixLength is the string length of the suffix generated by
// PrefixedUniqueId. This can be used by length validation functions to
// ensure prefixes are the correct length for the target field.
const UniqueIDSuffixLength = 20

// Helper for a resource to generate a unique identifier w/ given prefix
func PrefixedUniqueId(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, xid.New().String())
}
