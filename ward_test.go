package vnprovince

import (
	"errors"
	"testing"
)

func TestGetWardsLength(t *testing.T) {
	wards, err := GetWards()
	if err != nil {
		t.Fatal(err)
	}

	if len(wards) != WardsLength {
		t.Fatalf("len(wards) = %d, want %d", len(wards), WardsLength)
	}
}

func TestEachWard(t *testing.T) {
	type args struct {
		fn func(w Ward) error
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
				fn: func(d Ward) error { return errors.New("test") },
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EachWard(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("EachWard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
