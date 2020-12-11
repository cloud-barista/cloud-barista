package request

import (
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
)

const GrpcPort = 9999

//var once sync.Once
//var grpcMonClient pb.MONClient

func ConnectGRPC() (*grpc.ClientConn, error) {
	//grpcPort := config.GetInstance().GetGrpcConfig().Port
	opts := []grpc.DialOption{grpc.WithInsecure()}
	clientConn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", GrpcPort), opts...)
	if err != nil {
		return nil, err
	}
	return clientConn, nil
}

func ConnectClient(clientConn *grpc.ClientConn) *pb.MONClient {
	monClient := pb.NewMONClient(clientConn)
	return &monClient
}

/*func GetMonClient() pb.MONClient {
	once.Do(func() {
		//grpcPort := config.GetInstance().GetGrpcConfig().Port
		opts := []grpc.DialOption{grpc.WithInsecure()}
		clientConn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", GrpcPort), opts...)
		if err != nil {
			//return err
		}
		grpcMonClient = pb.NewMONClient(clientConn)
	})
	return grpcMonClient
}*/

/*func NewMonitoringRequest(timeout time.Duration) *MonitoringRequest {
	if timeout.Seconds() == 0 {
		timeout = TimeoutMinutes * time.Minute
	}
	mreq := MonitoringRequest{
		Client:  GetMonClient(),
		Timeout: timeout,
		InType:  ConvertType,
		OutType: ConvertType,
	}
	return &mreq
}*/
