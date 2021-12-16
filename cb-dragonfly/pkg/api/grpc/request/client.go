package request

import (
	"google.golang.org/grpc"

	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
)

func ConnectGRPC(grpcServerAddr string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	clientConn, err := grpc.Dial(grpcServerAddr, opts...)
	if err != nil {
		return nil, err
	}
	return clientConn, nil
}

func ConnectClient(clientConn *grpc.ClientConn) *pb.MONClient {
	monClient := pb.NewMONClient(clientConn)
	return &monClient
}
