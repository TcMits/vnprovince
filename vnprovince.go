package vnprovince

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

var StopIteration = io.EOF

// DivisionsLength the number of divisions in the data directory.
const DivisionsLength = 10604

// Division is a division of Vietnam.
type Division struct {
	ProvinceName string
	ProvinceCode int64
	DistrictName string
	DistrictCode int64
	WardName     string
	WardCode     int64
}

// GetDivisions returns all divisions in the data directory.
func GetDivisions(out *[]*Division) error {
	if out == nil {
		return errors.New("out is nil")
	}

	return EachDivision(func(d Division) error {
		*out = append(*out, &d)
		return nil
	})
}

// EachDivision calls fn for each division in the data directory.
func EachDivision(fn func(d Division) error) error {
	if fn == nil {
		return errors.New("fn is nil")
	}

	f, err := DataDirFS.Open(DivisionPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	d := new(Division)
	for {
		row, err := csvReader.Read()
		if err == StopIteration {
			break
		}

		if err != nil {
			return err
		}

		fromRow(row, d)

		if err := fn(*d); err != nil {
			return err
		}
	}

	return nil
}

// fromRow populates d from a row of the CSV file.
// it panics if the row is invalid.
func fromRow(row []string, d *Division) {
	d.ProvinceName = row[0]
	d.ProvinceCode = must(strconv.ParseInt(row[1], 10, 64))

	d.DistrictName = row[2]
	d.DistrictCode = must(strconv.ParseInt(row[3], 10, 64))

	if row[4] != "" {
		d.WardName = row[4]
		d.WardCode = must(strconv.ParseInt(row[5], 10, 64))
	} else {
		d.WardName = ""
		d.WardCode = 0
	}
}
