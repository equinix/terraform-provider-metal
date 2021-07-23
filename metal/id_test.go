// Licenses under the MPL 2.0
//
// Source: https://github.com/hashicorp/terraform-plugin-sdk/pull/508
package metal

import (
	"regexp"
	"strings"
	"testing"
)

var all36 = regexp.MustCompile(`^[a-z0-9]+$`)

func TestPrefixedUniqueId(t *testing.T) {
	prefix := "test-"
	iterations := 10000
	ids := make(map[string]struct{})
	var id string
	for i := 0; i < iterations; i++ {
		id = PrefixedUniqueId(prefix)

		if _, ok := ids[id]; ok {
			t.Fatalf("Got duplicated id! %s", id)
		}

		if !strings.HasPrefix(id, prefix) {
			t.Fatalf("Unique ID didn't have terraform- prefix! %s", id)
		}

		suffix := strings.TrimPrefix(id, prefix)

		if !all36.MatchString(suffix) {
			t.Fatalf("Suffix isn't in base 36! %s", suffix)
		}

		if len(suffix) > UniqueIDSuffixLength {
			t.Fatalf("Suffix length (%d) is larger than %d", len(suffix), UniqueIDSuffixLength)
		}

		ids[id] = struct{}{}
	}
}
