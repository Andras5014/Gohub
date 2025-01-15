package ioc

import (
	interactivev1 "github.com/Andras5014/webook/api/proto/gen/interactive/v1"
	"github.com/Andras5014/webook/interactive/config"
	grpc2 "github.com/Andras5014/webook/interactive/grpc"
	"github.com/Andras5014/webook/pkg/grpcx"
	"google.golang.org/grpc"
)

func InitGRPCxServer(config *config.Config, interactiveServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	server := grpc.NewServer()
	interactivev1.RegisterInteractiveServiceServer(server, interactiveServer)
	return &grpcx.Server{
		Server: server,
		Addr:   config.Grpc.Addr,
	}
}
