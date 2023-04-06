# vnprovince

vnprovince is a Golang library that provides a list of Vietnam administrative divisions, including tỉnh thành (provinces), quận huyện (districts), and phường xã (communes/wards). The library includes both the name and code for each administrative division, as defined by the General Statistics Office of Vietnam (Tổng cục Thống kê).

## Installation

You can install the library using `go get`:

```sh
go get github.com/TcMits/vnprovince
```

## Usage

To use the vnprovince library in your Golang project, simply import it and use the provided functions to access the administrative divisions. Here's a simple example:

go

```go
package main

import (
  "fmt"
  "github.com/TcMits/vnprovince"
)

func main() {
  provinces, err := vnprovince.GetProvinces()
  if err != nil {
    fmt.Println(err)
  }

  for _, province := range provinces {
    fmt.Println(province.Name)
  }
}
```

## Contributing

Contributions to the vnprovince library are welcome! If you find a bug or have an idea for a new feature, please open an issue or submit a pull request.

## Data

- https://danhmuchanhchinh.gso.gov.vn/
