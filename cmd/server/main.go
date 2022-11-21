package main

import (
	"context"
	"log"
	"net"

	"github.com/MyyPo/grpc-chat/internal/config"
	"github.com/MyyPo/grpc-chat/internal/repositories"
	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	chatpb "github.com/MyyPo/grpc-chat/pb/chat/v1"

	"github.com/MyyPo/grpc-chat/internal/service"
	"google.golang.org/grpc"
)

func main() {

	// load the env variables
	// config, err := util.NewConfig("./../..")

	// env in docker
	conf, err := config.NewConfig("./")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	authRepo := repositories.NewAuthRepo(db)
	res, err := authRepo.SignUp(context.Background(), &authpb.SignUpRequest{
		Username: "Mykyta1",
		Password: "hello!",
	})
	if err != nil {
		log.Fatalf("repo error: %v", err)
	}
	log.Fatalf(res.Username)

	basicPath := "/chat.v1.BroadcastService/"
	accessibleRoles := map[string][]string{
		basicPath + "ServerMessageStream": {"user"},
		basicPath + "ClientMessage":       {"user"},
	}

	impl := service.NewImplementation(*conf, accessibleRoles)

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
