package vnprovince

import (
	"errors"
	"testing"
)

func TestGetDivisions(t *testing.T) {
	type args struct {
		out *[]*Division
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
				out: new([]*Division),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetDivisions(tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("GetDivisions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDivisionsLength(t *testing.T) {
	divisions := make([]*Division, 0, DivisionsLength)
	if err := GetDivisions(&divisions); err != nil {
		t.Fatal(err)
	}

	if len(divisions) != DivisionsLength {
		t.Fatalf("len(divisions) = %d, want %d", len(divisions), DivisionsLength)
	}
}

func TestEachDivision(t *testing.T) {
	type args struct {
		fn func(d Division) error
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
				fn: func(d Division) error { return errors.New("test") },
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EachDivision(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("EachDivision() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
