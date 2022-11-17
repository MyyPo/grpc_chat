package client_service

import (
	"bufio"
	"context"
	"fmt"
	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	"google.golang.org/grpc"
)

type AuthClient struct {
	service authpb.AuthServiceClient
	scanner *bufio.Scanner
}

func NewSignInClient(conn *grpc.ClientConn, scanner *bufio.Scanner) *AuthClient {
	service := authpb.NewAuthServiceClient(conn)

	return &AuthClient{service, scanner}
}

func (client *AuthClient) SignIn(ctx context.Context) {

	for client.scanner.Scan() {
		req := &authpb.SignInRequest{
			Username: client.scanner.Text(),
			Password: "lol",
		}
		res, err := client.service.SignIn(context.Background(), req)

		if err != nil {
			fmt.Printf("Error while logging in: %v", err)
		} else {

			fmt.Println(res)

			break
		}
	}
}
