package vnprovince

import (
	"strconv"
	"unsafe"
)

func stob(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Division is a division of Vietnam.
type Division struct {
	ProvinceCode int64  `json:"provinceCode"`
	DistrictCode int64  `json:"districtCode"`
	WardCode     int64  `json:"wardCode"`
	ProvinceName string `json:"provinceName"`
	DistrictName string `json:"districtName"`
	WardName     string `json:"wardName"`
}

// EachDivision calls fn for each division in the data directory.
func EachDivision(fn func(d Division) bool) {
	if fn == nil {
		return
	}

	commaSep := [7]int{}
	startIdx := 0
	for i := range stob(dataDirFS) {
		switch dataDirFS[i] {
		case '\n':
			row := dataDirFS[startIdx:i]
			commaIdx := 0
			for j, c := range row {
				if c == ',' {
					commaSep[commaIdx] = j
					commaIdx++
				}
			}

			var wardCode int64 = 0
			if wcs := row[commaSep[4]+1 : commaSep[5]]; wcs != "" {
				wardCode = must(strconv.ParseInt(wcs, 10, 64))
			}

			if !fn(Division{
				ProvinceName: row[:commaSep[0]],
				ProvinceCode: must(strconv.ParseInt(row[commaSep[0]+1:commaSep[1]], 10, 64)),
				DistrictName: row[commaSep[1]+1 : commaSep[2]],
				DistrictCode: must(strconv.ParseInt(row[commaSep[2]+1:commaSep[3]], 10, 64)),
				WardName:     row[commaSep[3]+1 : commaSep[4]],
				WardCode:     wardCode,
			}) {
				return
			}

			startIdx = i + 1 // skip newline character
		}
	}
}
