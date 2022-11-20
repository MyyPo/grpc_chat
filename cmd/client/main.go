package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	client_service "github.com/MyyPo/grpc-chat/cmd/client/service"
	"google.golang.org/grpc"
)

var scanner *bufio.Scanner

func main() {
	basicPath := "/chat.v1.BroadcastService/"

	authMethods := map[string]bool{
		basicPath + "ServerMessageStream": true,
		basicPath + "ClientMessage":       true,
	}

	scanner = bufio.NewScanner(os.Stdin)

	tempConn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to the server")
	}
	authClient := client_service.NewSignInClient(tempConn, scanner)
	interceptor, err := client_service.NewAuthInterceptor(authClient, authMethods, 10*time.Minute)
	if err != nil {
		log.Fatalf("Failed to initialize interceptors %v", err)
	}

	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(
		"localhost:8080",
		transportOption,
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to the server")
	}
	scanner = bufio.NewScanner(os.Stdin)

	chatClient := client_service.NewChatClient(conn, scanner)

	chatClient.ServerMessageStream(context.Background())
	chatClient.ClientMessage(context.Background())

}
