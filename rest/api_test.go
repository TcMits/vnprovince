package rest

import (
	"context"
	"testing"

	"github.com/TcMits/vnprovince/rest/proto"
)

func Test_vnProvinceService_ListDivisions(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.ListDivisionsRequest{PageSize: -10}
	if _, err := s.ListDivisions(ctx, req); err == nil {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, true)
	}

	req = &proto.ListDivisionsRequest{PageSize: 10}
	resp, err := s.ListDivisions(ctx, req)
	if err != nil {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if len(resp.Divisions) != 10 {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken == "" {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListDivisionsRequest{PageSize: 10, Filter: "province_code=1 AND district_code=1 AND ward_code=25"}
	resp, err = s.ListDivisions(ctx, req)
	if err != nil {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if len(resp.Divisions) != 1 {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken != "" {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if resp.Divisions[0].Id != 27 {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListDivisionsRequest{PageSize: 10, Skip: 10}
	resp, err = s.ListDivisions(ctx, req)
	if err != nil {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if len(resp.Divisions) != 10 {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if resp.Divisions[0].Id == 3 {
		t.Errorf("ListDivisions() error = %v, wantErr %v", err, false)
	}
}

func Test_vnProvinceService_GetDivision(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.GetDivisionRequest{Name: "divisions/1"}
	if _, err := s.GetDivision(ctx, req); err == nil {
		t.Errorf("GetDivision() error = %v, wantErr %v", err, true)
	}

	req = &proto.GetDivisionRequest{Name: "divisions/3"}
	resp, err := s.GetDivision(ctx, req)
	if err != nil {
		t.Errorf("GetDivision() error = %v, wantErr %v", err, false)
	}

	if resp.Id != 3 {
		t.Errorf("GetDivision() error = %v, wantErr %v", err, false)
	}
}
