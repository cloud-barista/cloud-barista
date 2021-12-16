package mcar

import (
	"context"

	gc "github.com/cloud-barista/cb-mcks/src/grpc-api/common"
	pb "github.com/cloud-barista/cb-mcks/src/grpc-api/protobuf/cbmcks"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// Healthy - 상태확인
func (r *MCARRequest) Healthy() (string, error) {
	// 서버에 요청
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	resp, err := r.Client.Healthy(ctx, &pb.Empty{})

	if err != nil {
		return "", err
	}

	// 결과값 마샬링
	return gc.ConvertToOutput(r.OutType, &resp)
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
