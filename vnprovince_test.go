package vnprovince

import (
	"fmt"
	"strings"
	"testing"
)

func stringDivision(d Division) string {
	wardStr := ","
	if d.WardCode != 0 {
		wardStr = fmt.Sprintf("%s,%05d", d.WardName, d.WardCode)
	}

	return fmt.Sprintf(
		"%s,%02d,%s,%03d,",
		d.ProvinceName,
		d.ProvinceCode,
		d.DistrictName,
		d.DistrictCode,
	) + wardStr
}

func TestEachDivision(t *testing.T) {
	divisions := []Division{}
	EachDivision(func(d Division) bool {
		divisions = append(divisions, d)
		return true
	})

	lineIdx := 0
	startIdx := 0
	for i := range stob(dataDirFS) {
		switch dataDirFS[i] {
		case '\n':
			row := dataDirFS[startIdx:i]
			if !strings.HasPrefix(row, stringDivision(divisions[lineIdx])) {
				t.Fatalf("line %d: expected %q, got %q", lineIdx+1, stringDivision(divisions[lineIdx]), row)
			}

			startIdx = i + 1
			lineIdx++
		}
	}

	if lineIdx != len(divisions) {
		t.Fatalf("expected %d lines, got %d", lineIdx, len(divisions))
	}

	if startIdx != len(dataDirFS) {
		t.Fatalf("expected end of data at %d, got %d", len(dataDirFS), startIdx)
	}
}

func Benchmark_EachDivision(b *testing.B) {
	f := func(d Division) bool { return true }
	for i := 0; i < b.N; i++ {
		EachDivision(f)
	}
}
