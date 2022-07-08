package mcar

import (
	"context"
	"errors"

	gc "github.com/cloud-barista/cb-mcks/src/grpc-api/common"
	pb "github.com/cloud-barista/cb-mcks/src/grpc-api/protobuf/cbmcks"
)

// ListSpec - Spec 목록
func (r *MCARRequest) ListSpec() (string, error) {
	// 입력데이터 검사
	if r.InData == "" {
		return "", errors.New("input data required")
	}

	// 입력데이터 언마샬링
	var item pb.SpecQryRequest
	err := gc.ConvertToMessage(r.InType, r.InData, &item)
	if err != nil {
		return "", err
	}

	// 서버에 요청
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	resp, err := r.Client.ListSpec(ctx, &item)
	if err != nil {
		return "", err
	}

	// 결과값 마샬링
	return gc.ConvertToOutput(r.OutType, &resp)
}
