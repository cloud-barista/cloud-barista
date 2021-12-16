package server

import (
	"fmt"
	"net"

	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"google.golang.org/grpc"

	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
)

func StartGRPCServer() {
	grpcConfig := config.GetGrpcInstance()
	tcpConn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcConfig.GrpcServer.Ip, grpcConfig.GrpcServer.Port))
	if err != nil {
		util.GetLogger().Error("failed to listen server address: ", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMONServer(grpcServer, MonitoringService{})
	err = grpcServer.Serve(tcpConn)
	if err != nil {
		util.GetLogger().Error("failed to run grpc server: ", err)
		return
	}
}
