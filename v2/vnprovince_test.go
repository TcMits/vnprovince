package vnprovince

import "testing"

func Test_Duplicate(t *testing.T) {
	existing := map[Division]struct{}{}
	EachDivision(func(d Division) bool {
		if _, ok := existing[d]; ok {
			t.Fatalf("duplicate division found: %s, %s", d.ProvinceName, d.WardName)
		}

		existing[d] = struct{}{}
		return true
	})
}
