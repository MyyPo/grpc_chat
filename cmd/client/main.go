package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"google.golang.org/grpc/credentials/insecure"

	client_service "github.com/MyyPo/grpc-chat/cmd/client/service"
	"google.golang.org/grpc"
)

var scanner *bufio.Scanner

func main() {

	scanner = bufio.NewScanner(os.Stdin)

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to the server")
	}
	scanner = bufio.NewScanner(os.Stdin)

	authClient := client_service.NewSignInClient(conn, scanner)

	authClient.SignIn(context.Background())

	chatClient := client_service.NewChatClient(conn, scanner)

	chatClient.ServerMessageStream(context.Background())
	chatClient.ClientMessage(context.Background())

}
