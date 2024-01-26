package rest

import (
	"context"
	"errors"
	"strconv"

	"github.com/TcMits/aipstr"
	"github.com/TcMits/vnprovince"
	"github.com/TcMits/vnprovince/rest/proto"
	"github.com/alecthomas/participle/v2"
	"go.einride.tech/aip/pagination"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var stopIteration = errors.New("stop iteration")

type VNProvinceService struct {
	proto.UnimplementedVNProvinceServiceServer

	divisionDeclaration *aipstr.Declaration[selector]
	parser              *participle.Parser[aipstr.Filter]
}

func NewVNProvinceService() *VNProvinceService {
	return &VNProvinceService{
		divisionDeclaration: getFilterDeclaration(),
		parser:              aipstr.NewFilterParser(),
	}
}

func apiDivisionFromDivision(dst *proto.Division, src *vnprovince.Division) {
	dst.Name = proto.DivisionResourceName{DivisionId: strconv.FormatInt(src.ID(), 10)}.String()
	dst.Id = int32(src.ID())
	dst.ProvinceCode = int32(src.ProvinceCode)
	dst.DistrictCode = int32(src.DistrictCode)
	dst.WardCode = int32(src.WardCode)
	dst.ProvinceName = src.ProvinceName
	dst.DistrictName = src.DistrictName
	dst.WardName = src.WardName
}

func (s *VNProvinceService) ListDivisions(ctx context.Context, req *proto.ListDivisionsRequest) (*proto.ListDivisionsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	filter := trueSelector
	if filterStr := req.GetFilter(); filterStr != "" {
		filterStruct, err := s.parser.ParseString("", filterStr)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		filter, err = s.divisionDeclaration.WhereClause(filterStruct)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	pageToken, _ := pagination.ParsePageToken(req)
	pageSize := req.GetPageSize()
	resp := proto.ListDivisionsResponse{Divisions: make([]*proto.Division, 0, pageSize)}
	index := 0
	toOffset := int(int32(pageToken.Offset) + pageSize)
	offset := int(pageToken.Offset)

	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if !filter(&d) {
			return nil
		}

		if index >= offset && index < toOffset {
			division := proto.Division{}
			apiDivisionFromDivision(&division, &d)
			resp.Divisions = append(resp.Divisions, &division)
		}

		if index >= toOffset {
			resp.NextPageToken = pageToken.Next(req).String()
			return stopIteration
		}

		index += 1
		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &resp, nil
}

func (s *VNProvinceService) GetDivision(ctx context.Context, req *proto.GetDivisionRequest) (*proto.Division, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rn := proto.DivisionResourceName{}
	if err := errors.Join(rn.UnmarshalString(req.GetName()), rn.Validate()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	expectedId, err := strconv.ParseInt(rn.DivisionId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp := proto.Division{}
	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if d.ID() == expectedId {
			apiDivisionFromDivision(&resp, &d)
			return stopIteration
		}

		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.Id == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &resp, nil
}
