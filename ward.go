package vnprovince

import "errors"

const (
	// WardsLength is the number of wards.
	WardsLength   = 10599
	wardsCapacity = int(WardsLength / DistrictsLength)
)

// Ward is a ward in Vietnam.
type Ward struct {
	Code int64  `json:"code"`
	Name string `json:"name"`
}

// GetWards returns all wards.
func GetWards() ([]*Ward, error) {
	out := make([]*Ward, 0, WardsLength)

	if err := EachWard(func(w Ward) error {
		out = append(out, &w)
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

// EachWard iterates over all wards.
func EachWard(fn func(w Ward) error) error {
	if fn == nil {
		return errors.New("fn is nil")
	}

	return EachDivision(func(d Division) error {
		if d.WardCode == 0 {
			return nil
		}

		if err := fn(Ward{
			Code: d.WardCode,
			Name: d.WardName,
		}); err != nil {
			return err
		}

		return nil
	})
}

// wardFromDivision sets the ward from the division.
func wardFromDivision(d *Division, w *Ward) {
	w.Code = d.WardCode
	w.Name = d.WardName
}
