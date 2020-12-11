package mcir

import (
	"context"
	"fmt"

	gc "github.com/cloud-barista/cb-tumblebug/src/api/grpc/common"
	"github.com/cloud-barista/cb-tumblebug/src/api/grpc/logger"
	pb "github.com/cloud-barista/cb-tumblebug/src/api/grpc/protobuf/cbtumblebug"

	"github.com/cloud-barista/cb-tumblebug/src/core/mcir"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// CreateSpecWithInfo - Spec 생성
func (s *MCIRService) CreateSpecWithInfo(ctx context.Context, req *pb.TbSpecInfoRequest) (*pb.TbSpecInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCIRService.CreateSpecWithInfo()")

	// GRPC 메시지에서 MCIR 객체로 복사
	var mcirObj mcir.TbSpecInfo
	err := gc.CopySrcToDest(&req.Item, &mcirObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.CreateSpecWithInfo()")
	}

	content, err := mcir.RegisterSpecWithInfo(req.NsId, &mcirObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.CreateSpecWithInfo()")
	}

	// MCIR 객체에서 GRPC 메시지로 복사
	var grpcObj pb.TbSpecInfo
	err = gc.CopySrcToDest(&content, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.CreateSpecWithInfo()")
	}

	resp := &pb.TbSpecInfoResponse{Item: &grpcObj}
	return resp, nil
}

// CreateSpecWithSpecName - Spec 생성
func (s *MCIRService) CreateSpecWithSpecName(ctx context.Context, req *pb.TbSpecCreateRequest) (*pb.TbSpecInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCIRService.CreateSpecWithSpecName()")

	// GRPC 메시지에서 MCIR 객체로 복사
	var mcirObj mcir.TbSpecReq
	err := gc.CopySrcToDest(&req.Item, &mcirObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.CreateSpecWithSpecName()")
	}

	content, err := mcir.RegisterSpecWithCspSpecName(req.NsId, &mcirObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.CreateSpecWithSpecName()")
	}

	// MCIR 객체에서 GRPC 메시지로 복사
	var grpcObj pb.TbSpecInfo
	err = gc.CopySrcToDest(&content, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.CreateSpecWithSpecName()")
	}

	resp := &pb.TbSpecInfoResponse{Item: &grpcObj}
	return resp, nil
}

// ListSpec - Spec 목록
func (s *MCIRService) ListSpec(ctx context.Context, req *pb.ResourceAllQryRequest) (*pb.ListTbSpecInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCIRService.ListSpec()")

	resourceList, err := mcir.ListResource(req.NsId, req.ResourceType)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.ListSpec()")
	}

	// MCIR 객체에서 GRPC 메시지로 복사
	var grpcObj []*pb.TbSpecInfo
	err = gc.CopySrcToDest(&resourceList, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.ListSpec()")
	}

	resp := &pb.ListTbSpecInfoResponse{Items: grpcObj}
	return resp, nil
}

// GetSpec - Spec 조회
func (s *MCIRService) GetSpec(ctx context.Context, req *pb.ResourceQryRequest) (*pb.TbSpecInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCIRService.GetSpec()")

	res, err := mcir.GetResource(req.NsId, req.ResourceType, req.ResourceId)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.GetSpec()")
	}

	// MCIR 객체에서 GRPC 메시지로 복사
	var grpcObj pb.TbSpecInfo
	err = gc.CopySrcToDest(&res, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.GetSpec()")
	}

	resp := &pb.TbSpecInfoResponse{Item: &grpcObj}
	return resp, nil
}

// DeleteSpec - Spec 삭제
func (s *MCIRService) DeleteSpec(ctx context.Context, req *pb.ResourceQryRequest) (*pb.MessageResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCIRService.DeleteSpec()")

	err := mcir.DelResource(req.NsId, req.ResourceType, req.ResourceId, req.Force)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.DeleteSpec()")
	}

	resp := &pb.MessageResponse{Message: "The " + req.ResourceType + " " + req.ResourceId + " has been deleted"}
	return resp, nil
}

// DeleteAllSpec - Spec 전체 삭제
func (s *MCIRService) DeleteAllSpec(ctx context.Context, req *pb.ResourceAllQryRequest) (*pb.MessageResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCIRService.DeleteAllSpec()")

	err := mcir.DelAllResources(req.NsId, req.ResourceType, req.Force)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.DeleteAllSpec()")
	}

	resp := &pb.MessageResponse{Message: "All " + req.ResourceType + "s has been deleted"}
	return resp, nil
}

// FetchSpec - Spec 가져오기
func (s *MCIRService) FetchSpec(ctx context.Context, req *pb.FetchSpecQryRequest) (*pb.MessageResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCIRService.FetchSpec()")

	connConfigCount, specCount, err := mcir.FetchSpecs(req.NsId)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.FetchSpec()")
	}

	resp := &pb.MessageResponse{Message: "Fetched " + fmt.Sprint(specCount) + " specs (from " + fmt.Sprint(connConfigCount) + " connConfigs)"}
	return resp, nil
}

// FilterSpec
func (s *MCIRService) FilterSpec(ctx context.Context, req *pb.TbSpecInfoRequest) (*pb.ListTbSpecInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCIRService.FilterSpec()")

	// GRPC 메시지에서 MCIR 객체로 복사
	var mcirObj mcir.TbSpecInfo
	err := gc.CopySrcToDest(&req.Item, &mcirObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.FilterSpec()")
	}

	resourceList, err := mcir.FilterSpecs(req.NsId, mcirObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.FilterSpec()")
	}

	// MCIR 객체에서 GRPC 메시지로 복사
	var grpcObj []*pb.TbSpecInfo
	err = gc.CopySrcToDest(&resourceList, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCIRService.FilterSpec()")
	}

	resp := &pb.ListTbSpecInfoResponse{Items: grpcObj}
	return resp, nil
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
