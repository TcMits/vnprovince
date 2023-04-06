package vnprovince

import (
	"errors"
	"testing"
)

func TestGetDistricts(t *testing.T) {
	type args struct {
		out *[]*District
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				out: nil,
			},
			wantErr: true,
		},
		{
			name: "t2",
			args: args{
				out: new([]*District),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetDistricts(tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("GetDistricts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDistrictsLength(t *testing.T) {
	districts := make([]*District, 0, DistrictsLength)
	if err := GetDistricts(&districts); err != nil {
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
