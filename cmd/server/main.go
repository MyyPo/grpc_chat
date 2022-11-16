package main

import (
	"log"
	"net"
	"os"

	chatpb "github.com/MyyPo/grpc-chat/pb/chat/v1"
	"github.com/MyyPo/grpc-chat/service"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

func main() {

	server := service.NewChatServer(grpcLog)

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to create the server %v", err)
	}

	grpcLog.Info("The server successfuly started")

	chatpb.RegisterBroadcastServiceServer(grpcServer, server)
	grpcServer.Serve(lis)
}
