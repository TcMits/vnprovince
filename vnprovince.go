package vnprovince

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

// DivisionsLength the number of divisions in the data directory.
const DivisionsLength = 10604

// Division is a division of Vietnam.
type Division struct {
	ProvinceCode int64  `json:"provinceCode"`
	DistrictCode int64  `json:"districtCode"`
	WardCode     int64  `json:"wardCode"`
	ProvinceName string `json:"provinceName"`
	DistrictName string `json:"districtName"`
	WardName     string `json:"wardName"`
}

func (d *Division) ID() int64 {
	if d == nil {
		return 0
	}

	return d.WardCode + d.DistrictCode + d.ProvinceCode
}

// GetDivisions returns all divisions in the data directory.
func GetDivisions() ([]*Division, error) {
	out := make([]*Division, 0, DivisionsLength)

	if err := EachDivision(func(d Division) error {
		out = append(out, &d)
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

// EachDivision calls fn for each division in the data directory.
func EachDivision(fn func(d Division) error) error {
	if fn == nil {
		return errors.New("fn is nil")
	}

	f, err := DataDirFS.Open(DivisionPath)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if err := fn(Division{
			ProvinceName: row[0],
			ProvinceCode: must(strconv.ParseInt(row[1], 10, 64)),
			DistrictName: row[2],
			DistrictCode: must(strconv.ParseInt(row[3], 10, 64)),
			WardName:     row[4],
			WardCode:     ignore(strconv.ParseInt(row[5], 10, 64)),
		}); err != nil {
			return err
		}
	}

	return nil
}
