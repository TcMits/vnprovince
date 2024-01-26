package api

import (
	"strings"

	"github.com/TcMits/aipstr"
	"github.com/TcMits/vnprovince"
)

type selector func(d *vnprovince.Division) bool

func andSelector(s ...selector) selector {
	return func(d *vnprovince.Division) bool {
		for _, selector := range s {
			if !selector(d) {
				return false
			}
		}

		return true
	}
}

func orSelector(s ...selector) selector {
	return func(d *vnprovince.Division) bool {
		for _, selector := range s {
			if selector(d) {
				return true
			}
		}

		return false
	}
}

func trueSelector(d *vnprovince.Division) bool {
	return true
}

func falseSelector(d *vnprovince.Division) bool {
	return false
}

func eqFieldWithInt(field string, value int64) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "id":
			return d.ID() == value
		case "province_code":
			return d.ProvinceCode == value
		case "district_code":
			return d.DistrictCode == value
		case "ward_code":
			return d.WardCode == value
		}

		return false
	}, nil
}

func ltFieldWithInt(field string, value int64) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "id":
			return d.ID() < value
		case "province_code":
			return d.ProvinceCode < value
		case "district_code":
			return d.DistrictCode < value
		case "ward_code":
			return d.WardCode < value
		}

		return false
	}, nil
}

func gtFieldWithInt(field string, value int64) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "id":
			return d.ID() > value
		case "province_code":
			return d.ProvinceCode > value
		case "district_code":
			return d.DistrictCode > value
		case "ward_code":
			return d.WardCode > value
		}

		return false
	}, nil
}

func lteFieldWithInt(field string, value int64) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "id":
			return d.ID() <= value
		case "province_code":
			return d.ProvinceCode <= value
		case "district_code":
			return d.DistrictCode <= value
		case "ward_code":
			return d.WardCode <= value
		}

		return false
	}, nil
}

func gteFieldWithInt(field string, value int64) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "id":
			return d.ID() >= value
		case "province_code":
			return d.ProvinceCode >= value
		case "district_code":
			return d.DistrictCode >= value
		case "ward_code":
			return d.WardCode >= value
		}

		return false
	}, nil
}

func neFieldWithInt(field string, value int64) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "id":
			return d.ID() != value
		case "province_code":
			return d.ProvinceCode != value
		case "district_code":
			return d.DistrictCode != value
		case "ward_code":
			return d.WardCode != value
		}

		return false
	}, nil
}

func eqFieldWithString(field string, value string) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "province_name":
			return d.ProvinceName == value
		case "district_name":
			return d.DistrictName == value
		case "ward_name":
			return d.WardName == value
		}

		return false
	}, nil
}

func ltFieldWithString(field string, value string) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "province_name":
			return d.ProvinceName < value
		case "district_name":
			return d.DistrictName < value
		case "ward_name":
			return d.WardName < value
		}

		return false
	}, nil
}

func gtFieldWithString(field string, value string) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "province_name":
			return d.ProvinceName > value
		case "district_name":
			return d.DistrictName > value
		case "ward_name":
			return d.WardName > value
		}

		return false
	}, nil
}

func lteFieldWithString(field string, value string) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "province_name":
			return d.ProvinceName <= value
		case "district_name":
			return d.DistrictName <= value
		case "ward_name":
			return d.WardName <= value
		}

		return false
	}, nil
}

func gteFieldWithString(field string, value string) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "province_name":
			return d.ProvinceName >= value
		case "district_name":
			return d.DistrictName >= value
		case "ward_name":
			return d.WardName >= value
		}

		return false
	}, nil
}

func neFieldWithString(field string, value string) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "province_name":
			return d.ProvinceName != value
		case "district_name":
			return d.DistrictName != value
		case "ward_name":
			return d.WardName != value
		}

		return false
	}, nil
}

func hasFieldWithString(field string, value string) (selector, error) {
	return func(d *vnprovince.Division) bool {
		switch field {
		case "province_name":
			return strings.Contains(
				strings.ToLower(d.ProvinceName),
				strings.ToLower(value),
			)
		case "district_name":
			return strings.Contains(
				strings.ToLower(d.DistrictName),
				strings.ToLower(value),
			)
		case "ward_name":
			return strings.Contains(
				strings.ToLower(d.WardName),
				strings.ToLower(value),
			)
		}

		return false
	}, nil
}

func getBasicOperator() []*aipstr.DeclarationOperatorFunc[selector] {
	return []*aipstr.DeclarationOperatorFunc[selector]{
		aipstr.NewOperatorFunc(
			aipstr.EqOp,
			aipstr.WithFieldWithValueInt(eqFieldWithInt),
			aipstr.WithFieldWithValueString(eqFieldWithString),
		),
		aipstr.NewOperatorFunc(
			aipstr.LtOp,
			aipstr.WithFieldWithValueInt(ltFieldWithInt),
			aipstr.WithFieldWithValueString(ltFieldWithString),
		),
		aipstr.NewOperatorFunc(
			aipstr.GtOp,
			aipstr.WithFieldWithValueInt(gtFieldWithInt),
			aipstr.WithFieldWithValueString(gtFieldWithString),
		),
		aipstr.NewOperatorFunc(
			aipstr.LeOp,
			aipstr.WithFieldWithValueInt(lteFieldWithInt),
			aipstr.WithFieldWithValueString(lteFieldWithString),
		),
		aipstr.NewOperatorFunc(
			aipstr.GeOp,
			aipstr.WithFieldWithValueInt(gteFieldWithInt),
			aipstr.WithFieldWithValueString(gteFieldWithString),
		),
		aipstr.NewOperatorFunc(
			aipstr.NeOp,
			aipstr.WithFieldWithValueInt(neFieldWithInt),
			aipstr.WithFieldWithValueString(neFieldWithString),
		),
		aipstr.NewOperatorFunc(
			aipstr.HasOp,
			aipstr.WithFieldWithValueString(hasFieldWithString),
		),
		aipstr.NewOperatorFunc(
			aipstr.AndOp,
			aipstr.WithCombineNoErr(andSelector),
		),
		aipstr.NewOperatorFunc(
			aipstr.OrOp,
			aipstr.WithCombineNoErr(orSelector),
		),
		aipstr.NewOperatorFunc(
			aipstr.TrueOp,
			aipstr.WithNoField(func() selector { return trueSelector }),
		),
		aipstr.NewOperatorFunc(
			aipstr.FalseOp,
			aipstr.WithNoField(func() selector { return falseSelector }),
		),
	}
}

func getFilterDeclaration() *aipstr.Declaration[selector] {
	return aipstr.NewDeclaration(
		aipstr.WithColumns(
			aipstr.NewColumn("id", aipstr.Filterable[selector]()),
			aipstr.NewColumn("province_code", aipstr.Filterable[selector]()),
			aipstr.NewColumn("province_name", aipstr.Filterable[selector]()),
			aipstr.NewColumn("district_code", aipstr.Filterable[selector]()),
			aipstr.NewColumn("district_name", aipstr.Filterable[selector]()),
			aipstr.NewColumn("ward_code", aipstr.Filterable[selector]()),
			aipstr.NewColumn("ward_name", aipstr.Filterable[selector]()),
		),
		aipstr.WithOperatorFuncs(getBasicOperator()...),
	)
}
