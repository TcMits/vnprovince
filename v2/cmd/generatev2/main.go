//go:build ignore

package main

import (
	"bytes"
	_ "embed"
	"go/format"
	"os"
	"unsafe"
)

//go:embed data.csv
var data string

const prefixTmpl = `
package vnprovince

type Division struct {
	ProvinceName string
	OldDistrictName string
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
			if !strings.Contains(row, div.ProvinceName) || !strings.Contains(row, div.OldDistrictName) || !strings.Contains(row, div.WardName) {
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
	ProvinceName    string `json:"provinceName"`
	OldDistrictName string `json:"oldDistrictName"`
	WardName        string `json:"wardName"`
}

func eachDivision(fn func(d Division) bool) {
	if fn == nil {
		return
	}

	commaSep := [11]int{}
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
				ProvinceName:    row[commaSep[2]+1 : commaSep[3]],
				OldDistrictName: row[commaSep[5]+1 : commaSep[6]],
				WardName:        row[commaSep[8]+1 : commaSep[9]],
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
		buf.WriteString(`"` + d.OldDistrictName + `",`)
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
