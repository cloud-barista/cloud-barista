package mcar

import (
	"context"

	gc "github.com/cloud-barista/cb-mcks/src/grpc-api/common"
	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
	pb "github.com/cloud-barista/cb-mcks/src/grpc-api/protobuf/cbmcks"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/core/service"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// CreateCluster - Cluster 생성
func (s *MCARService) CreateCluster(ctx context.Context, req *pb.ClusterCreateRequest) (*pb.ClusterInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.CreateCluster()")

	// GRPC 메시지에서 MCKS 객체로 복사
	var mcarObj app.ClusterReq
	err := gc.CopySrcToDest(&req.Item, &mcarObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.CreateCluster()")
	}

	s.ClusterReqDef(&mcarObj)
	err = s.ClusterReqValidate(mcarObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.CreateCluster()")
	}

	cluster, err := service.CreateCluster(req.Namespace, req.Minorversion, req.Patchversion, &mcarObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.CreateCluster()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.ClusterInfo
	err = gc.CopySrcToDest(&cluster, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.CreateCluster()")
	}

	resp := &pb.ClusterInfoResponse{Item: &grpcObj}
	return resp, nil
}

// ListCluster - Cluster 목록
func (s *MCARService) ListCluster(ctx context.Context, req *pb.ClusterAllQryRequest) (*pb.ListClusterInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.ListCluster()")

	clusterList, err := service.ListCluster(req.Namespace)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.ListCluster()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.ListClusterInfoResponse
	err = gc.CopySrcToDest(&clusterList, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.ListCluster()")
	}

	return &grpcObj, nil
}

// GetCluster - Cluster 조회
func (s *MCARService) GetCluster(ctx context.Context, req *pb.ClusterQryRequest) (*pb.ClusterInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.GetCluster()")

	if err := s.Validate(map[string]string{"namespace": req.Namespace, "cluster": req.Cluster}); err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.GetCluster()")
	}

	cluster, err := service.GetCluster(req.Namespace, req.Cluster)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.GetCluster()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.ClusterInfo
	err = gc.CopySrcToDest(&cluster, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.GetCluster()")
	}

	resp := &pb.ClusterInfoResponse{Item: &grpcObj}
	return resp, nil
}

// DeleteCluster - Cluster 삭제
func (s *MCARService) DeleteCluster(ctx context.Context, req *pb.ClusterQryRequest) (*pb.StatusResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.DeleteCluster()")

	if err := s.Validate(map[string]string{"namespace": req.Namespace, "cluster": req.Cluster}); err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.DeleteCluster()")
	}

	status, err := service.DeleteCluster(req.Namespace, req.Cluster)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.DeleteCluster()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.StatusResponse
	err = gc.CopySrcToDest(&status, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.DeleteCluster()")
	}

	return &grpcObj, nil
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
