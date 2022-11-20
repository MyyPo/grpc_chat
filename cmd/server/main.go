package main

import (
	"log"
	"net"

	"github.com/MyyPo/grpc-chat/internal/util"
	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	chatpb "github.com/MyyPo/grpc-chat/pb/chat/v1"

	"github.com/MyyPo/grpc-chat/internal/service"
	"google.golang.org/grpc"
)

func main() {

	// load the env variables
	config, err := util.NewConfig("./../..")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	impl := service.NewImplementation(*config)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(impl.AuthInterceptor.Unary()),
		grpc.StreamInterceptor(impl.AuthInterceptor.Stream()),
	)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to create the server %v", err)
	}

	log.Println("The server successfuly started")

	chatpb.RegisterBroadcastServiceServer(grpcServer, &impl)
	authpb.RegisterAuthServiceServer(grpcServer, &impl)
	grpcServer.Serve(lis)
}
