//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"net/http"
	"os"
	"strings"

	v1 "github.com/TcMits/vnprovince"
	v2 "github.com/TcMits/vnprovince/v2"
)

// based on https://vi.wikipedia.org/wiki/Danh_sách_đơn_vị_hành_chính_Việt_Nam_trong_đợt_cải_cách_thể_chế_2024–2025

const prefixTmpl = `
package vnprovince


func V1IndexToV2Index(idx int) (int, bool) {
	if idx < 0 || idx >= len(v1MapV2) {
		return 0, false
	}


	result := v1MapV2[idx]
	return result, result != -1
}

var v1MapV2 = [...]int{
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

	type districtData struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	type wardData struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	type convertData struct {
		Status int `json:"status"`
		Data   []struct {
			ProvinceNameNew string `json:"ProvinceNameNew"`
			WardNameNew     string `json:"wardNameNew"`
		} `json:"data"`
	}

	provincesResp := must(http.Get("https://diachi.vnpost.vn/api/address/option/provinces?type=1"))
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
			districtsResp := must(http.Get("https://diachi.vnpost.vn/api/address/option/districts?type=1&provinceCode=" + p.Code))
			defer districtsResp.Body.Close()
			districts := make([]districtData, 0)

			if err := json.NewDecoder(districtsResp.Body).Decode(&districts); err != nil {
				panic(err)
			}

			for _, d := range districts {
				switch {
				case strings.HasPrefix(d.Name, "TP. "):
					d.Name = "Thành phố " + strings.TrimPrefix(d.Name, "TP. ")
				case strings.HasPrefix(d.Name, "TP."):
					d.Name = "Thành phố " + strings.TrimPrefix(d.Name, "TP.")
				case strings.HasPrefix(d.Name, "TX. "):
					d.Name = "Thị xã " + strings.TrimPrefix(d.Name, "TX. ")
				case strings.HasPrefix(d.Name, "TX "):
					d.Name = "Thị xã " + strings.TrimPrefix(d.Name, "TX ")
				case strings.HasPrefix(d.Name, "H. "):
					d.Name = "Huyện " + strings.TrimPrefix(d.Name, "H. ")
				case strings.HasPrefix(d.Name, "H "):
					d.Name = "Huyện " + strings.TrimPrefix(d.Name, "H ")
				case strings.HasPrefix(d.Name, "Q. "):
					d.Name = "Quận " + strings.TrimPrefix(d.Name, "Q. ")
				case strings.HasPrefix(d.Name, "Huyện "):
				default:
					panic("district name does not start with TP., TX., H., or Q.: " + d.Name)
				}

				func() {
					wardsResp := must(http.Get("https://diachi.vnpost.vn/api/address/option/wards?type=1&districtCode=" + d.Code + "&provinceCode=" + p.Code))
					defer wardsResp.Body.Close()
					wards := make([]wardData, 0)

					if err := json.NewDecoder(wardsResp.Body).Decode(&wards); err != nil {
						panic(err)
					}

					mapwards := make(map[string]struct{})
					for _, w := range wards {
						switch {
						case strings.HasPrefix(w.Name, "P. "):
							w.Name = "Phường " + strings.TrimPrefix(w.Name, "P. ")
						case strings.HasPrefix(w.Name, "X. "):
							w.Name = "Xã " + strings.TrimPrefix(w.Name, "X. ")
						case strings.HasPrefix(w.Name, "TT. "):
							w.Name = "Thị trấn " + strings.TrimPrefix(w.Name, "TT. ")
						case strings.HasPrefix(w.Name, "TT."):
							w.Name = "Thị trấn " + strings.TrimPrefix(w.Name, "TT.")
						case strings.HasPrefix(w.Name, "Đ. "):
							w.Name = "Đảo " + strings.TrimPrefix(w.Name, "Đ. ")
						case strings.HasPrefix(w.Name, "Đặc khu "):
						case w.Name == "":
							continue
						case w.Name == "Ea Chà Rang":
							w.Name = "Xã Ea Chà Rang"
						default:
							panic("ward name does not start with P. or X. or TT.: " + w.Name)
						}

						if _, ok := mapwards[w.Name]; ok {
							continue
						}

						mapwards[w.Name] = struct{}{}
						v1Idx := 0
						foundv1 := false
						v1.EachDivision(func(di v1.Division) bool {
							if di.ProvinceName == p.Name && di.DistrictName == d.Name && di.WardName == w.Name {
								foundv1 = true
								return false
							}

							v1Idx++
							return true
						})

						if !foundv1 {
							panic("could not find division in v1: " + p.Name + ", " + d.Name + ", " + w.Name)
						}

						func() {
							v2Idx := 0
							defer func() {
								buf.WriteString(fmt.Sprintf("\t%d,\n", v2Idx))
							}()
							reader := strings.NewReader(
								fmt.Sprintf(`{"direction":1,"provinceCode":"` + p.Code + `","districtCode":"` + d.Code + `","wardCode":"` + w.Code + `","pageIndex":0,"pageSize":10,"provinceName":"","districtName":"","wardName":"","name":true,"bcCode":false,"bcqgCode":false,"inputText":"","territory":1,"captchaCode":"2165"}`),
							)

							newWardResp := must(http.Post("https://diachi.vnpost.vn/api/address/convert", "application/json", reader))
							defer newWardResp.Body.Close()

							newWard := convertData{}
							if err := json.NewDecoder(newWardResp.Body).Decode(&newWard); err != nil {
								panic(err)
							}

							if newWard.Status != 200 || len(newWard.Data) == 0 {
								fmt.Println("can't not convert ward name for", p.Name, d.Name, w.Name, "status:", newWard.Status)
								v2Idx = -1
								return
							}

							switch {
							case strings.HasPrefix(newWard.Data[0].WardNameNew, "P. "):
								newWard.Data[0].WardNameNew = "Phường " + strings.TrimPrefix(newWard.Data[0].WardNameNew, "P. ")
							case strings.HasPrefix(newWard.Data[0].WardNameNew, "X. "):
								newWard.Data[0].WardNameNew = "Xã " + strings.TrimPrefix(newWard.Data[0].WardNameNew, "X. ")
							case strings.HasPrefix(newWard.Data[0].WardNameNew, "TT. "):
								newWard.Data[0].WardNameNew = "Thị trấn " + strings.TrimPrefix(newWard.Data[0].WardNameNew, "TT. ")
							case strings.HasPrefix(newWard.Data[0].WardNameNew, "Đặc khu "):
							default:
								panic("ward name does not start with P. or X. or TT.: " + newWard.Data[0].WardNameNew)
							}

							switch {
							case strings.HasPrefix(newWard.Data[0].ProvinceNameNew, "TP. "):
								newWard.Data[0].ProvinceNameNew = "Thành phố " + strings.TrimPrefix(newWard.Data[0].ProvinceNameNew, "TP. ")
							case strings.HasPrefix(newWard.Data[0].ProvinceNameNew, "Tỉnh "):
							default:
								panic("province name does not start with TP. or Tỉnh: " + newWard.Data[0].ProvinceNameNew)
							}

							findedv2 := false
							v2.EachDivision(func(di v2.Division) bool {
								if di.ProvinceName == newWard.Data[0].ProvinceNameNew && di.WardName == newWard.Data[0].WardNameNew {
									findedv2 = true
									return false
								}

								v2Idx++
								return true
							})

							if !findedv2 {
								println("could not find division in v2: " + newWard.Data[0].ProvinceNameNew + ", " + newWard.Data[0].WardNameNew)
								v2Idx = -1
							}
						}()
					}
				}()
			}
		}()
	}

	buf.WriteString("}\n")

	if err := os.WriteFile("v1_map.go", must(format.Source(buf.Bytes())), os.ModePerm); err != nil {
		panic(err)
	}
}
