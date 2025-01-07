package main

import (
	interactivev1 "github.com/Andras5014/webook/api/proto/gen/interactive/v1"
	"github.com/Andras5014/webook/interactive/grpc"
	grpc2 "google.golang.org/grpc"
	"net"
)

func main() {
	server := grpc2.NewServer()
	intrSvc := &grpc.InteractiveServiceServer{}
	interactivev1.RegisterInteractiveServiceServer(server, intrSvc)
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = server.Serve(l)
	if err != nil {
		panic(err)
	}
}
