package main

import (
	// "context"
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

	basicPath := "/chat.v1.BroadcastService/"
	accessibleRoles := map[string][]string{
		basicPath + "ServerMessageStream": {"user"},
		basicPath + "ClientMessage":       {"user"},
	}

	impl := service.NewImplementation(*conf, accessibleRoles, *authRepo)

	// res, err := impl.SignUp(context.Background(), &authpb.SignUpRequest{
	// 	Username: "TestingHasher",
	// 	Password: "TestingHasher",
	// })
	// if err != nil {
	// 	log.Fatalf("error signing UP: %v", err)
	// }
	// res, err := impl.SignIn(context.Background(), &authpb.SignInRequest{
	// 	Username: "TestingHasher",
	// 	Password: "TestingHasher",
	// })
	// if err != nil {
	// 	log.Fatalf("error signing in: %v", err)
	// }
	// log.Fatalln(res)

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
