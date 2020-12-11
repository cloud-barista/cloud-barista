package server

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
)

func StartGRPCServer() {
	grpcPort := config.GetInstance().GetGrpcConfig().Port
	tcpConn, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if err != nil {
		logrus.Error("failed to listen server address: ", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMONServer(grpcServer, MonitoringService{})
	err = grpcServer.Serve(tcpConn)
	if err != nil {
		logrus.Error("failed to run grpc server: ", err)
		return
	}
}
