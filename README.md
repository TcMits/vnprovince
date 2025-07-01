### vnprovince

vnprovince provides a list of Vietnam administrative divisions

### Install

```sh
go get github.com/TcMits/vnprovince/v2@main
```

### Updates

- use `github.com/TcMits/vnprovince/v2@main` for 01/07/2025 divisions


### Examples

go

```go
package main

import (
  "fmt"
  "github.com/TcMits/vnprovince"
)

func main() {
	vnprovince.EachDivision(func(d Division) bool {
		return true
	})
}
```

### Data

- https://danhmuchanhchinh.gso.gov.vn/
- https://easyinvoice.vn/easyinvoice-cap-nhat-danh-sach-xa-phuong-moi-2025-sau-sap-nhap/
