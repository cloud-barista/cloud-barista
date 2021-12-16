package mcar

import (
	"context"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
	pb "github.com/cloud-barista/cb-mcks/src/grpc-api/protobuf/cbmcks"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// Healthy - 상태확인
func (s *MCARService) Healthy(ctx context.Context, req *pb.Empty) (*pb.MessageResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.Healthy()")

	resp := &pb.MessageResponse{Message: "cb-barista cb-mcks"}
	return resp, nil
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
