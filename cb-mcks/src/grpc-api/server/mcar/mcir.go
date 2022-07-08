package mcar

import (
	"context"
	"strconv"

	gc "github.com/cloud-barista/cb-mcks/src/grpc-api/common"
	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
	pb "github.com/cloud-barista/cb-mcks/src/grpc-api/protobuf/cbmcks"

	"github.com/cloud-barista/cb-mcks/src/core/service"
)

func (s *MCARService) ListSpec(ctx context.Context, req *pb.SpecQryRequest) (*pb.ListSpecInfoResponse, error) {
	logger := logger.NewLogger()

	logger.Debug("calling MCARService.ListSpec()")

	control_plane := req.ControlPlane
	if len(req.ControlPlane) == 0 {
		control_plane = "N"
	}

	cpumin, _ := strconv.Atoi(req.CpuMin)
	cpumax, _ := strconv.Atoi(req.CpuMax)
	memorymin, _ := strconv.Atoi(req.MemoryMin)
	memorymax, _ := strconv.Atoi(req.MemoryMax)

	if req.ControlPlane == "Y" {
		if len(req.CpuMin) == 0 {
			cpumin = 2
		}
		if len(req.MemoryMin) == 0 {
			memorymin = 2
		}
	} else {
		if len(req.CpuMin) == 0 {
			cpumin = 1
		}
		if len(req.MemoryMin) == 0 {
			memorymin = 1
		}
	}

	if len(req.CpuMax) == 0 {
		cpumax = 99999
	}
	if len(req.MemoryMax) == 0 {
		memorymax = 99999
	}

	specList, err := service.VerifySpecList(req.Connectionname, control_plane, cpumin, cpumax, memorymin, memorymax)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.ListSpec()")
	}

	// MCKS 객체에서 GRPC 메시지로 복사
	var grpcObj pb.ListSpecInfoResponse
	err = gc.CopySrcToDest(&specList, &grpcObj)
	if err != nil {
		return nil, gc.ConvGrpcStatusErr(err, "", "MCARService.ListSpec()")
	}

	return &grpcObj, nil
}
