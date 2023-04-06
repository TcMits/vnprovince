package vnprovince

import (
	"errors"
	"testing"
)

func TestGetProvinces(t *testing.T) {
	provinces, err := GetProvinces()
	if err != nil {
		t.Fatal(err)
	}

	if len(provinces) != ProvincesLength {
		t.Fatalf("len(provinces) = %d, want %d", len(provinces), ProvincesLength)
	}

	districtsLength := 0
	wardsLength := 0
	for _, p := range provinces {
		districtsLength += len(p.Districts)
		for _, d := range p.Districts {
			wardsLength += len(d.Wards)
		}
	}

	if districtsLength != DistrictsLength {
		t.Fatalf("districtsLength = %d, want %d", districtsLength, DistrictsLength)
	}

	if wardsLength != WardsLength {
		t.Fatalf("wardsLength = %d, want %d", wardsLength, WardsLength)
	}
}

func TestEachProvince(t *testing.T) {
	type args struct {
		fn func(p Province) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				fn: nil,
			},
			wantErr: true,
		},
		{
			name: "t2",
			args: args{
				fn: func(p Province) error { return errors.New("") },
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EachProvince(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("EachProvince() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
