//go:build ignore

package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"go/format"
	"net/http"
	"os"
	"strings"
)

const prefixTmpl = `
package vnprovince

type Division struct {
	ProvinceName string
	WardName     string
}

func EachDivision(fn func(d Division) bool) {
	for _, d := range divisions {
		if !fn(d) {
			return
		}
	}
}

func AtIndex(idx int) (Division, bool){
	if idx < 0 || idx >= len(divisions) {
		return Division{}, false
	}

	return divisions[idx], true
}

func Len() int {
	return len(divisions)
}

var divisions = [...]Division{
`

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	buf := bytes.Buffer{}
	buf.WriteString(prefixTmpl)

	type provinceData struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	type wardData struct {
		Name string `json:"name"`
	}

	provincesResp := must(http.Get("https://diachi.vnpost.vn/api/address/option/provinces?type=2"))
	defer provincesResp.Body.Close()
	provinces := make([]provinceData, 0)
	if err := json.NewDecoder(provincesResp.Body).Decode(&provinces); err != nil {
		panic(err)
	}

	for _, p := range provinces {
		switch {
		case strings.HasPrefix(p.Name, "TP. "):
			p.Name = "Thành phố " + strings.TrimPrefix(p.Name, "TP. ")
		case strings.HasPrefix(p.Name, "Tỉnh "):
		default:
			panic("province name does not start with TP. or Tỉnh: " + p.Name)
		}

		func() {
			wardsResp := must(http.Get("https://diachi.vnpost.vn/api/address/option/wards?type=2&districtCode=null&provinceCode=" + p.Code))
			defer wardsResp.Body.Close()
			wards := make([]wardData, 0)

			if err := json.NewDecoder(wardsResp.Body).Decode(&wards); err != nil {
				panic(err)
			}

			for _, w := range wards {
				switch {
				case strings.HasPrefix(w.Name, "P. "):
					w.Name = "Phường " + strings.TrimPrefix(w.Name, "P. ")
				case strings.HasPrefix(w.Name, "X. "):
					w.Name = "Xã " + strings.TrimPrefix(w.Name, "X. ")
				case strings.HasPrefix(w.Name, "TT. "):
					w.Name = "Thị trấn " + strings.TrimPrefix(w.Name, "TT. ")
				case strings.HasPrefix(w.Name, "Đặc khu "):
				default:
					panic("ward name does not start with P. or X. or TT.: " + w.Name)
				}

				buf.WriteString("\t{")
				buf.WriteString(`"` + p.Name + `",`)
				buf.WriteString(`"` + w.Name + `"},`)
				buf.WriteByte('\n')
			}
		}()
	}

	buf.WriteString("}\n")

	if err := os.WriteFile("vnprovince.go", must(format.Source(buf.Bytes())), os.ModePerm); err != nil {
		panic(err)
	}
}
