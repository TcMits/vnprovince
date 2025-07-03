package vnprovince

import (
	"testing"

	v1 "github.com/TcMits/vnprovince"
)

func Test_V1IndexToV2Index(t *testing.T) {
	v1Idx := 0
	v1.EachDivision(func(d v1.Division) bool {
		if d.ProvinceName == "Thành phố Hồ Chí Minh" && d.DistrictName == "Quận 3" && d.WardName == "Phường 14" {
			return false
		}

		v1Idx++
		return true
	})

	v2Idx, ok := V1IndexToV2Index(v1Idx)
	if !ok {
		t.Fatalf("Expected to find index for v1 index %d, but got none", v1Idx)
	}

	d, ok := AtIndex(v2Idx)
	if !ok {
		t.Fatalf("Expected to find division at index %d, but got none", v2Idx)
	}

	if d.ProvinceName != "Thành phố Hồ Chí Minh" {
		t.Fatalf("Expected province name 'Thành phố Hồ Chí Minh', got '%s'", d.ProvinceName)
	}

	if d.WardName != "Phường Nhiêu Lộc" {
		t.Fatalf("Expected ward name 'Phường Nhiêu Lộc', got '%s'", d.WardName)
	}
}
