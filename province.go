package vnprovince

import "errors"

// ProvincesLength is the number of provinces.
const ProvincesLength = 63

// Province represents a province.
type Province struct {
	Name      string     `json:"name"`
	Code      int64      `json:"code"`
	Districts []District `json:"districts"`
}

// GetProvinces returns all provinces and districts.
func GetProvinces() ([]*Province, error) {
	out := make([]*Province, 0, ProvincesLength)

	err := EachProvince(func(p Province) error {
		out = append(out, &p)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, err
}

// EachProvince iterates over all provinces and districts.
func EachProvince(fn func(p Province) error) error {
	if fn == nil {
		return errors.New("fn is nil")
	}

	var previousCode int64 = 1
	currentProvince := &Province{
		Districts: make([]District, 0, 2),
	}

	err := EachDivision(func(d Division) error {
		if previousCode != d.ProvinceCode {
			if err := fn(*currentProvince); err != nil {
				return err
			}

			// update previousCode
			previousCode = d.ProvinceCode
			currentProvince.Districts = make([]District, 0, 2)
		}

		provinceFromDivision(&d, currentProvince)
		return nil
	})
	if err != nil {
		return err
	}

	// handle the last province
	if err := fn(*currentProvince); err != nil {
		return err
	}

	return nil
}

// provinceFromDivision converts a division to a province.
func provinceFromDivision(d *Division, p *Province) {
	p.Code = d.ProvinceCode
	p.Name = d.ProvinceName

	var currentDistrict *District
	for i := range p.Districts {
		district := &p.Districts[i]
		if district.Code == d.DistrictCode {
			currentDistrict = district
		}
	}

	if currentDistrict == nil {
		p.Districts = append(p.Districts, District{})
		currentDistrict = &p.Districts[len(p.Districts)-1]
	}

	districtFromDivision(d, currentDistrict)
}
