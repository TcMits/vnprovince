### vnprovince

vnprovince provides a list of Vietnam administrative divisions

### Features

- get all divisions
- convert old divisions to 01/07/2025 divisions

### Install

```sh
go get github.com/TcMits/vnprovince/v2@main
```

### Updates

- use `github.com/TcMits/vnprovince/v2@main` for 01/07/2025 divisions


### Examples

Loop

```go
package main

import (
  "github.com/TcMits/vnprovince"
)

func main() {
	vnprovince.EachDivision(func(d Division) bool {
		return true
	})
}
```

Convert v1 to v2

```go
package main

import (
  "log"
  v1 "github.com/TcMits/vnprovince"
  v2 "github.com/TcMits/vnprovince/v2"
)

func main() {
	v1Idx := 0
	v1.EachDivision(func(d v1.Division) bool {
		if d.ProvinceName == "Thành phố Hồ Chí Minh" && d.DistrictName == "Quận 3" && d.WardName == "Phường 14" {
			return false
		}

		v1Idx++
		return true
	})

	v2Idx, ok := v2.V1IndexToV2Index(v1Idx)
	if !ok {
		log.Fatalf("Expected to find index for v1 index %d, but got none", v1Idx)
	}

	d, ok := v2.AtIndex(v2Idx)
	if !ok {
		log.Fatalf("Expected to find division at index %d, but got none", v2Idx)
	}

	if d.ProvinceName != "Thành phố Hồ Chí Minh" {
		log.Fatalf("Expected province name 'Thành phố Hồ Chí Minh', got '%s'", d.ProvinceName)
	}

	if d.WardName != "Phường Nhiêu Lộc" {
		log.Fatalf("Expected ward name 'Phường Nhiêu Lộc', got '%s'", d.WardName)
	}
}
```


### Data

- https://danhmuchanhchinh.gso.gov.vn/
- https://easyinvoice.vn/easyinvoice-cap-nhat-danh-sach-xa-phuong-moi-2025-sau-sap-nhap/
- https://vi.wikipedia.org/wiki/Danh_sách_đơn_vị_hành_chính_Việt_Nam_trong_đợt_cải_cách_thể_chế_2024–2025
- https://diachi.vnpost.vn
