package main

import (
	"bytes"
	_ "embed"
	"go/format"
	"os"
	"unsafe"
)

//go:embed divisions_16_10_2024.csv
var data string

const prefixTmpl = `
package vnprovince

type Division struct {
	ProvinceName string
	DistrictName string
	WardName     string
}

func EachDivision(fn func(d Division) bool) {
	for _, d := range divisions {
		if !fn(d) {
			return
		}
	}
}


var divisions = [...]Division{
`

const prefixTestTmpl = `
package vnprovince

import (
	"strings"
	"testing"
	"unsafe"
)

func stob(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func TestEachDivision(t *testing.T) {
	divisions := []Division{}
	EachDivision(func(d Division) bool {
		divisions = append(divisions, d)
		return true
	})

	lineIdx := 0
	startIdx := 0
	for i := range stob(data) {
		switch data[i] {
		case '\n':
			row := data[startIdx:i]
			div := divisions[lineIdx]
			if !strings.Contains(row, div.ProvinceName) || !strings.Contains(row, div.DistrictName) || !strings.Contains(row, div.WardName) {
				t.Fatalf("line %d: expected %+v, got %q", lineIdx+1, divisions[lineIdx], row)
			}

			startIdx = i + 1
			lineIdx++
		}
	}

	if lineIdx != len(divisions) {
		t.Fatalf("expected %d lines, got %d", lineIdx, len(divisions))
	}

	if startIdx != len(data) {
		t.Fatalf("expected end of data at %d, got %d", len(data), startIdx)
	}
}

const data = `

func stob(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

type Division struct {
	ProvinceName string `json:"provinceName"`
	DistrictName string `json:"districtName"`
	WardName     string `json:"wardName"`
}

func eachDivision(fn func(d Division) bool) {
	if fn == nil {
		return
	}

	commaSep := [7]int{}
	startIdx := 0
	for i := range stob(data) {
		switch data[i] {
		case '\n':
			row := data[startIdx:i]
			commaIdx := 0
			for j, c := range stob(row) {
				if c == ',' {
					commaSep[commaIdx] = j
					commaIdx++
				}
			}

			if !fn(Division{
				ProvinceName: row[:commaSep[0]],
				DistrictName: row[commaSep[1]+1 : commaSep[2]],
				WardName:     row[commaSep[3]+1 : commaSep[4]],
			}) {
				return
			}

			startIdx = i + 1 // skip newline character
		}
	}
}

func main() {
	buf := bytes.Buffer{}
	buf.WriteString(prefixTmpl)
	eachDivision(func(d Division) bool {
		buf.WriteString("\t{")
		buf.WriteString(`"` + d.ProvinceName + `",`)
		buf.WriteString(`"` + d.DistrictName + `",`)
		buf.WriteString(`"` + d.WardName + `"},`)
		buf.WriteByte('\n')
		return true
	})

	buf.WriteString("}\n")

	if err := os.WriteFile("vnprovince.go", must(format.Source(buf.Bytes())), os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile("vnprovince_test.go", must(format.Source([]byte(prefixTestTmpl+"`"+data+"`"))), os.ModePerm); err != nil {
		panic(err)
	}
}
