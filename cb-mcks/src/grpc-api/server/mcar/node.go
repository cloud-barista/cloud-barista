package mcar

import (
	"context"

	gc "github.com/cloud-barista/cb-mcks/src/grpc-api/common"
	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
	pb "github.com/cloud-barista/cb-mcks/src/grpc-api/protobuf/cbmcks"

	"github.com/cloud-barista/cb-mcks/src/core/model"
	"github.com/cloud-barista/cb-mcks/src/core/service"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// AddNode - Node 추가
func (s *MCARService) AddNode(ctx context.Context, req *pb.NodeCreateRequest) (*pb.ListNodeInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.AddNode()")

	if err := s.Validate(map[string]string{"cluster": req.Cluster}); err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.AddNode()")
	}

	// GRPC 메시지에서 MCKS 객체로 복사
	var mcarObj model.NodeReq
	err := gc.CopySrcToDest(&req.Item, &mcarObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.AddNode()")
	}

	err = s.NodeReqValidate(mcarObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.AddNode()")
	}

	node, err := service.AddNode(req.Namespace, req.Cluster, &mcarObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.AddNode()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.ListNodeInfoResponse
	err = gc.CopySrcToDest(&node, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.ListCluster()")
	}

	return &grpcObj, nil
}

// ListNode - Node 목록
func (s *MCARService) ListNode(ctx context.Context, req *pb.NodeAllQryRequest) (*pb.ListNodeInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.ListNode()")

	if err := s.Validate(map[string]string{"cluster": req.Cluster}); err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.ListNode()")
	}

	nodeList, err := service.ListNode(req.Namespace, req.Cluster)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.ListNode()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.ListNodeInfoResponse
	err = gc.CopySrcToDest(&nodeList, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.ListNode()")
	}

	return &grpcObj, nil
}

// GetNode - Node 조회
func (s *MCARService) GetNode(ctx context.Context, req *pb.NodeQryRequest) (*pb.NodeInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.GetNode()")

	if err := s.Validate(map[string]string{"cluster": req.Cluster, "node": req.Node}); err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.GetNode()")
	}

	node, err := service.GetNode(req.Namespace, req.Cluster, req.Node)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.GetNode()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.NodeInfo
	err = gc.CopySrcToDest(&node, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.GetNode()")
	}

	resp := &pb.NodeInfoResponse{Item: &grpcObj}
	return resp, nil
}

// RemoveNode - Node 삭제
func (s *MCARService) RemoveNode(ctx context.Context, req *pb.NodeQryRequest) (*pb.StatusResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.RemoveNode()")

	if err := s.Validate(map[string]string{"cluster": req.Cluster, "node": req.Node}); err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.GetNode()")
	}

	status, err := service.RemoveNode(req.Namespace, req.Cluster, req.Node)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.RemoveNode()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.StatusResponse
	err = gc.CopySrcToDest(&status, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.RemoveNode()")
	}

	return &grpcObj, nil
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
