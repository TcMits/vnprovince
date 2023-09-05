package vnprovince

import (
	"errors"
)

const (
	// districtsCapacity is the number of districts.
	DistrictsLength   = 705
	districtsCapacity = int(DistrictsLength / ProvincesLength)
)

// District is a district in Vietnam.
type District struct {
	Code  int64  `json:"code"`
	Name  string `json:"name"`
	Wards []Ward `json:"wards"`
}

// GetDistricts returns all districts and wards.
func GetDistricts() ([]*District, error) {
	out := make([]*District, 0, DistrictsLength)

	if err := EachDistrict(func(d District) error {
		out = append(out, &d)
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

// EachDistrict iterates over all districts and wards.
func EachDistrict(fn func(d District) error) error {
	if fn == nil {
		return errors.New("fn is nil")
	}

	var previousCode int64 = 1
	currentDistrict := &District{
		Wards: make([]Ward, 0, wardsCapacity),
	}

	if err := EachDivision(func(d Division) error {
		if previousCode != d.DistrictCode {
			if err := fn(*currentDistrict); err != nil {
				return err
			}

			// update previousCode
			previousCode = d.DistrictCode
			currentDistrict.Wards = make([]Ward, 0, wardsCapacity)
		}

		districtFromDivision(&d, currentDistrict)
		return nil
	}); err != nil {
		return err
	}

	// handle the last district
	if err := fn(*currentDistrict); err != nil {
		return err
	}

	return nil
}

// districtFromDivision updates the district from the division.
func districtFromDivision(d *Division, dist *District) {
	dist.Code = d.DistrictCode
	dist.Name = d.DistrictName

	if d.WardCode == 0 {
		return
	}

	// ward is the smallest unit of division
	dist.Wards = append(dist.Wards, Ward{})
	wardFromDivision(d, &dist.Wards[len(dist.Wards)-1])
}
