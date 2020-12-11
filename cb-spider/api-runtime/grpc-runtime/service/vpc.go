// gRPC Runtime of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2020.09.

package service

import (
	"context"
	"strings"

	gc "github.com/cloud-barista/cb-spider/api-runtime/grpc-runtime/common"
	"github.com/cloud-barista/cb-spider/api-runtime/grpc-runtime/logger"
	pb "github.com/cloud-barista/cb-spider/api-runtime/grpc-runtime/stub/cbspider"

	cmrt "github.com/cloud-barista/cb-spider/api-runtime/common-runtime"
	cres "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// CreateVPC - VPC 생성
func (s *CCMService) CreateVPC(ctx context.Context, req *pb.VPCCreateRequest) (*pb.VPCInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling CCMService.CreateVPC()")

	// check the input Name to include the SUBNET: Prefix
	if strings.HasPrefix(req.Item.Name, rsSubnetPrefix) {
		return nil, gc.NewGrpcStatusErr(rsSubnetPrefix+" cannot be used for VPC name prefix!!", "", "CCMService.CreateVPC()")
	}
	// check the input Name to include the SecurityGroup Delimiter
	if strings.HasPrefix(req.Item.Name, sgDELIMITER) {
		return nil, gc.NewGrpcStatusErr(sgDELIMITER+" cannot be used in VPC name!!", "", "CCMService.CreateVPC()")
	}

	// Grpc RegInfo => Driver ReqInfo
	// (1) create SubnetInfo List
	subnetInfoList := []cres.SubnetInfo{}
	for _, info := range req.Item.SubnetInfoList {
		subnetInfo := cres.SubnetInfo{IId: cres.IID{NameId: info.Name, SystemId: ""}, IPv4_CIDR: info.Ipv4Cidr}
		subnetInfoList = append(subnetInfoList, subnetInfo)
	}
	// (2) create VPCReqInfo with SubnetInfo List
	reqInfo := cres.VPCReqInfo{
		IId:            cres.IID{NameId: req.Item.Name, SystemId: ""},
		IPv4_CIDR:      req.Item.Ipv4Cidr,
		SubnetInfoList: subnetInfoList,
	}

	// Call common-runtime API
	result, err := cmrt.CreateVPC(req.ConnectionName, rsVPC, reqInfo)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.CreateVPC()")
	}

	// CCM 객체에서 GRPC 메시지로 복사
	var grpcObj pb.VPCInfo
	err = gc.CopySrcToDest(result, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.CreateVPC()")
	}

	resp := &pb.VPCInfoResponse{Item: &grpcObj}
	return resp, nil
}

// ListVPC - VPC 목록
func (s *CCMService) ListVPC(ctx context.Context, req *pb.VPCAllQryRequest) (*pb.ListVPCInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling CCMService.ListVPC()")

	// Call common-runtime API
	result, err := cmrt.ListVPC(req.ConnectionName, rsVPC)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.ListVPC()")
	}

	// CCM 객체에서 GRPC 메시지로 복사
	var grpcObj []*pb.VPCInfo
	err = gc.CopySrcToDest(&result, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.ListVPC()")
	}

	resp := &pb.ListVPCInfoResponse{Items: grpcObj}
	return resp, nil
}

// GetVPC - VPC 조회
func (s *CCMService) GetVPC(ctx context.Context, req *pb.VPCQryRequest) (*pb.VPCInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling CCMService.GetVPC()")

	// Call common-runtime API
	result, err := cmrt.GetVPC(req.ConnectionName, rsVPC, req.Name)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.GetVPC()")
	}

	// CCM 객체에서 GRPC 메시지로 복사
	var grpcObj pb.VPCInfo
	err = gc.CopySrcToDest(result, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.GetVPC()")
	}

	resp := &pb.VPCInfoResponse{Item: &grpcObj}
	return resp, nil
}

// DeleteVPC - VPC 삭제
func (s *CCMService) DeleteVPC(ctx context.Context, req *pb.VPCQryRequest) (*pb.BooleanResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling CCMService.DeleteVPC()")

	// Call common-runtime API
	result, _, err := cmrt.DeleteResource(req.ConnectionName, rsVPC, req.Name, req.Force)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.DeleteVPC()")
	}

	resp := &pb.BooleanResponse{Result: result}
	return resp, nil
}

// ListAllVPC - 관리 VPC 목록
func (s *CCMService) ListAllVPC(ctx context.Context, req *pb.VPCAllQryRequest) (*pb.AllResourceInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling CCMService.ListAllVPC()")

	// Call common-runtime API
	allResourceList, err := cmrt.ListAllResource(req.ConnectionName, rsVPC)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.ListAllVPC()")
	}

	// CCM 객체에서 GRPC 메시지로 복사
	var grpcObj pb.AllResourceInfoResponse
	err = gc.CopySrcToDest(&allResourceList, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.ListAllVPC()")
	}

	return &grpcObj, nil
}

// DeleteCSPVPC - CSP VPC 삭제
func (s *CCMService) DeleteCSPVPC(ctx context.Context, req *pb.CSPVPCQryRequest) (*pb.BooleanResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling CCMService.DeleteCSPVPC()")

	// Call common-runtime API
	result, _, err := cmrt.DeleteCSPResource(req.ConnectionName, rsVPC, req.Id)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "CCMService.DeleteCSPVPC()")
	}

	resp := &pb.BooleanResponse{Result: result}
	return resp, nil
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
