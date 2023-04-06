package vnprovince

import "errors"

// WardsLength is the number of wards.
const WardsLength = 10599

// Ward is a ward in Vietnam.
type Ward struct {
	Name string `json:"name"`
	Code int64  `json:"code"`
}

// GetWards returns all wards.
func GetWards() ([]*Ward, error) {
	out := make([]*Ward, 0, WardsLength)

	err := EachWard(func(w Ward) error {
		out = append(out, &w)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, err
}

// EachWard iterates over all wards.
func EachWard(fn func(w Ward) error) error {
	if fn == nil {
		return errors.New("fn is nil")
	}

	currentWard := new(Ward)

	return EachDivision(func(d Division) error {
		if d.WardCode == 0 {
			return nil
		}

		wardFromDivision(&d, currentWard)
		if err := fn(*currentWard); err != nil {
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
