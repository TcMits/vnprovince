package vnprovince

import (
	"errors"
	"testing"
)

func TestGetDistricts(t *testing.T) {
	districts, err := GetDistricts()
	if err != nil {
		t.Fatal(err)
	}

	if len(districts) != DistrictsLength {
		t.Fatalf("len(districts) = %d, want %d", len(districts), DistrictsLength)
	}

	wardsLength := 0
	for _, d := range districts {
		wardsLength += len(d.Wards)
	}

	if wardsLength != WardsLength {
		t.Fatalf("wardsLength = %d, want %d", wardsLength, WardsLength)
	}
}

func TestEachDistrict(t *testing.T) {
	type args struct {
		fn func(d District) error
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
				fn: func(d District) error { return errors.New("") },
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EachDistrict(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("EachDistrict() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
