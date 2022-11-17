package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	chatpb "github.com/MyyPo/grpc-chat/pb/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var chatClient chatpb.BroadcastServiceClient
var authClient authpb.AuthServiceClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *chatpb.User) error {
	var streamerr error

	stream, err := chatClient.ServerMessageStream(context.Background(), &chatpb.ServerMessageStreamRequest{
		User:   user,
		Active: true,
	})

	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	wait.Add(1)
	go func(stream chatpb.BroadcastService_ServerMessageStreamClient) {
		defer wait.Done()

		for {
			msg, err := stream.Recv()

			if err != nil {
				streamerr = fmt.Errorf("failed to read a message: %v", err)
				break
			}

			fmt.Printf("%v : %s\n", msg.Id, msg.Content)
		}
	}(stream)

	return streamerr
}

// func login(req *authpb.SignInRequest) (*authpb.SignInResponse, error) {
// func login(req *authpb.SignInRequest) {
func login() {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		req := &authpb.SignInRequest{
			Username: scanner.Text(),
			Password: "lol",
		}
		res, err := authClient.SignIn(context.Background(), req)

		if err != nil {
			fmt.Printf("Error while logging in: %v", err)
		} else {

			fmt.Println(res)

			break
		}
	}
}

func main() {
	timestamp := time.Now()

	done := make(chan int)

	name := flag.String("N", "Anon", "User name")
	flag.Parse()

	id := sha256.Sum256([]byte(timestamp.String() + *name))

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to the server")
	}

	authClient = authpb.NewAuthServiceClient(conn)
	// cred := &authpb.SignInRequest{
	// 	Username: "Anon",
	// 	Password: "lol",
	// }
	// login(cred)
	login()

	chatClient = chatpb.NewBroadcastServiceClient(conn)
	user := &chatpb.User{
		Id:   hex.EncodeToString(id[:]),
		Name: *name,
	}

	connect(user)

	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &chatpb.Message{
				Id:        user.Id,
				Content:   scanner.Text(),
				Timestamp: timestamppb.New(time.Now()),
			}

			_, err := chatClient.ClientMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error while sending a message: %v", err)
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
