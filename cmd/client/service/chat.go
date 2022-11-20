package client_service

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	// "flag"
	"fmt"
	"sync"
	"time"

	chatpb "github.com/MyyPo/grpc-chat/pb/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatClient struct {
	service chatpb.BroadcastServiceClient
	scanner *bufio.Scanner
	wait    *sync.WaitGroup
}

func NewChatClient(conn *grpc.ClientConn, scanner *bufio.Scanner) *ChatClient {
	service := chatpb.NewBroadcastServiceClient(conn)
	wait := &sync.WaitGroup{}

	return &ChatClient{service, scanner, wait}
}

func (client *ChatClient) ServerMessageStream(ctx context.Context) error {
	var streamerr error

	timestamp := timestamppb.Now()
	// name := flag.String("N", "Anon", "User name")
	// flag.Parse()
	id := sha256.Sum256([]byte(timestamp.String() + "Name"))
	user := &chatpb.User{
		Id:   hex.EncodeToString(id[:]),
		Name: "Name",
	}

	stream, err := client.service.ServerMessageStream(ctx, &chatpb.ServerMessageStreamRequest{
		User:   user,
		Active: true,
	})

	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	client.wait.Add(1)
	go func(stream chatpb.BroadcastService_ServerMessageStreamClient) {
		defer client.wait.Done()

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

func (client *ChatClient) ClientMessage(ctx context.Context) {
	timestamp := timestamppb.Now()
	// name := flag.Parse()
	// name := flag.String("N", "Anon", "User name")
	id := sha256.Sum256([]byte(timestamp.String() + "Name"))

	user := &chatpb.User{
		Id:   hex.EncodeToString(id[:]),
		Name: "Name",
	}

	done := make(chan int)

	client.wait.Add(1)
	go func() {
		defer client.wait.Done()

		for client.scanner.Scan() {
			msg := &chatpb.Message{
				Id:        user.Id,
				Content:   client.scanner.Text(),
				Timestamp: timestamppb.New(time.Now()),
			}

			_, err := client.service.ClientMessage(ctx, msg)
			if err != nil {
				fmt.Printf("Error while sending a message: %v", err)
				break
			}
		}
	}()

	go func() {
		client.wait.Wait()
		close(done)
	}()

	<-done
}
