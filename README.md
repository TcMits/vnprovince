### vnprovince

vnprovince provides a list of Vietnam administrative divisions

### Install

```sh
go get github.com/TcMits/vnprovince@main
```


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
