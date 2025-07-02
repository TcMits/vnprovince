//go:build ignore

package main

import (
	"bytes"
	"embed"
	"encoding/csv"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
	"unicode"

	v1 "github.com/TcMits/vnprovince"
	v2 "github.com/TcMits/vnprovince/v2"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//go:embed *.csv
var fs embed.FS

// based on https://vi.wikipedia.org/wiki/Danh_sách_đơn_vị_hành_chính_Việt_Nam_trong_đợt_cải_cách_thể_chế_2024–2025

const prefixTmpl = `
package vnprovince


func V1IndexToV2Index(idx int) (int, bool) {
	result, ok := v1MapV2[idx]
	return result, ok
}

var v1MapV2 = map[int]int{
`

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var fileMap = map[string]string{
	"Tỉnh An Giang":         "an_giang.csv",
	"Tỉnh Bắc Ninh":         "bac_ninh.csv",
	"Tỉnh Cà Mau":           "ca_mau.csv",
	"Thành phố Cần Thơ":     "can_tho.csv",
	"Tỉnh Cao Bằng":         "cao_bang.csv",
	"Thành phố Đà Nẵng":     "da_nang.csv",
	"Tỉnh Đắk Lắk":          "dak_lak.csv",
	"Tỉnh Điện Biên":        "dien_bien.csv",
	"Tỉnh Đồng Nai":         "dong_nai.csv",
	"Tỉnh Đồng Tháp":        "dong_thap.csv",
	"Tỉnh Gia Lai":          "gia_lai.csv",
	"Thành phố Hà Nội":      "ha_noi.csv",
	"Tỉnh Hà Tĩnh":          "ha_tinh.csv",
	"Thành phố Hải Phòng":   "hai_phong.csv",
	"Thành phố Hồ Chí Minh": "ho_chi_minh.csv",
	"Thành phố Huế":         "hue.csv",
	"Tỉnh Hưng Yên":         "hung_yen.csv",
	"Tỉnh Khánh Hòa":        "khanh_hoa.csv",
	"Tỉnh Lai Châu":         "lai_chau.csv",
	"Tỉnh Lâm Đồng":         "lam_dong.csv",
	"Tỉnh Lạng Sơn":         "lang_son.csv",
	"Tỉnh Lào Cai":          "lao_cai.csv",
	"Tỉnh Nghệ An":          "nghe_an.csv",
	"Tỉnh Ninh Bình":        "ninh_binh.csv",
	"Tỉnh Phú Thọ":          "phu_tho.csv",
	"Tỉnh Quảng Ngãi":       "quang_ngai.csv",
	"Tỉnh Quảng Ninh":       "quang_ninh.csv",
	"Tỉnh Quảng Trị":        "quang_tri.csv",
	"Tỉnh Sơn La":           "son_la.csv",
	"Tỉnh Tây Ninh":         "tay_ninh.csv",
	"Tỉnh Thái Nguyên":      "thai_nguyen.csv",
	"Tỉnh Thanh Hóa":        "thanh_hoa.csv",
	"Tỉnh Tuyên Quang":      "tuyen_quang.csv",
	"Tỉnh Vĩnh Long":        "vinh_long.csv",
}

var provinceMap = map[string]string{
	"Tỉnh Quảng Nam":         "Thành phố Đà Nẵng",
	"Tỉnh Bình Định":         "Tỉnh Gia Lai",
	"Tỉnh Long An":           "Tỉnh Tây Ninh",
	"Tỉnh Bắc Giang":         "Tỉnh Bắc Ninh",
	"Tỉnh Hải Dương":         "Thành phố Hải Phòng",
	"Tỉnh Yên Bái":           "Tỉnh Lào Cai",
	"Tỉnh Hà Giang":          "Tỉnh Tuyên Quang",
	"Tỉnh Bắc Kạn":           "Tỉnh Thái Nguyên",
	"Tỉnh Vĩnh Phúc":         "Tỉnh Phú Thọ",
	"Tỉnh Hoà Bình":          "Tỉnh Phú Thọ",
	"Tỉnh Thái Bình":         "Tỉnh Hưng Yên",
	"Tỉnh Hà Nam":            "Tỉnh Ninh Bình",
	"Tỉnh Nam Định":          "Tỉnh Ninh Bình",
	"Tỉnh Bến Tre":           "Tỉnh Vĩnh Long",
	"Tỉnh Trà Vinh":          "Tỉnh Vĩnh Long",
	"Tỉnh Bà Rịa - Vũng Tàu": "Thành phố Hồ Chí Minh",
	"Tỉnh Bình Dương":        "Thành phố Hồ Chí Minh",
	"Tỉnh Bình Phước":        "Tỉnh Đồng Nai",
	"Tỉnh Hậu Giang":         "Thành phố Cần Thơ",
	"Tỉnh Sóc Trăng":         "Thành phố Cần Thơ",
	"Tỉnh Kiên Giang":        "Tỉnh An Giang",
	"Tỉnh Tiền Giang":        "Tỉnh Đồng Tháp",
	"Tỉnh Bạc Liêu":          "Tỉnh Cà Mau",
	"Tỉnh Ninh Thuận":        "Tỉnh Khánh Hòa",
	"Tỉnh Phú Yên":           "Tỉnh Đắk Lắk",
	"Tỉnh Bình Thuận":        "Tỉnh Lâm Đồng",
	"Tỉnh Đắk Nông":          "Tỉnh Lâm Đồng",
	"Tỉnh Kon Tum":           "Tỉnh Quảng Ngãi",
	"Tỉnh Quảng Bình":        "Tỉnh Quảng Trị",
	"Tỉnh Thừa Thiên Huế":    "Thành phố Huế",
}

func main() {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	buf := bytes.Buffer{}
	buf.WriteString(prefixTmpl)
	v1Idx := 0
	v1.EachDivision(func(d v1.Division) bool {
		newName := provinceMap[d.ProvinceName]
		if newName == "" {
			newName = d.ProvinceName
		}

		fileName := fileMap[newName]
		if fileName == "" {
			panic("file not found for " + newName)
		}

		file := must(fs.Open(fileName))
		defer file.Close()

		newWardName := ""

		r := csv.NewReader(file)
		for {
			row, err := r.Read()
			if err == io.EOF {
				break
			}

			col0, _, err := transform.String(t, strings.TrimSpace(row[0]))
			if err != nil {
				panic(fmt.Sprintf("Error transforming string %s: %v", row[0], err))
			}

			col1, _, err := transform.String(t, strings.TrimSpace(row[1]))
			if err != nil {
				panic(fmt.Sprintf("Error transforming string %s: %v", row[1], err))
			}

			col2, _, err := transform.String(t, strings.TrimSpace(row[2]))
			if err != nil {
				panic(fmt.Sprintf("Error transforming string %s: %v", row[2], err))
			}

			dn, _, err := transform.String(t, d.DistrictName)
			if err != nil {
				panic(fmt.Sprintf("Error transforming string %s: %v", d.DistrictName, err))
			}

			wn, _, err := transform.String(t, d.WardName)
			if err != nil {
				panic(fmt.Sprintf("Error transforming string %s: %v", d.WardName, err))
			}

			if col2 == "" {
				continue
			}

			if col2 == wn {
				newWardName = col2
				break
			}

			if !strings.Contains(strings.ToUpper(col0), strings.ToUpper(dn)) && col0 != "" {
				continue
			}

			compareWardName := strings.TrimPrefix(d.WardName, "Xã ")
			compareWardName = strings.TrimPrefix(compareWardName, "Thị trấn ")
			compareWardName = strings.TrimPrefix(compareWardName, "Phường ")
			wn, _, err = transform.String(t, compareWardName)
			if err != nil {
				panic(fmt.Sprintf("Error transforming string %s: %v", d.WardName, err))
			}

			if !strings.Contains(strings.ToUpper(col1), strings.ToUpper(wn)) {
				continue
			}

			newWardName = col2
			break
		}

		if newWardName == "" {
			fmt.Printf("No new ward name found for %s, %s, %s\n", d.ProvinceName, d.DistrictName, d.WardName)
			v1Idx++
			return true
		}

		pn, _, err := transform.String(t, newName)
		if err != nil {
			panic(fmt.Sprintf("Error transforming string %s: %v", newName, err))
		}

		v2Idx := 0
		finded := false
		v2.EachDivision(func(d2 v2.Division) bool {
			compareProvinceName, _, err := transform.String(t, d2.ProvinceName)
			if err != nil {
				panic(fmt.Sprintf("Error transforming string %s: %v", d2.ProvinceName, err))
			}

			compareWardName, _, err := transform.String(t, d2.WardName)
			if err != nil {
				panic(fmt.Sprintf("Error transforming string %s: %v", d2.WardName, err))
			}

			if compareProvinceName == pn && compareWardName == newWardName {
				finded = true
				return false
			}

			v2Idx++
			return true
		})

		if !finded {
			fmt.Printf("No v2 division found for %s, %s\n", newName, newWardName)
			v1Idx++
			return true
		}

		buf.WriteString(fmt.Sprintf("\t%d: %d,\n", v1Idx, v2Idx))
		v1Idx++
		return true
	})
	buf.WriteString("}\n")

	if err := os.WriteFile("v1_map.go", must(format.Source(buf.Bytes())), os.ModePerm); err != nil {
		panic(err)
	}
}
