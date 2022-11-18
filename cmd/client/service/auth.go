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

	for {
		username, password := client.getCredentials()

		req := &authpb.SignInRequest{
			Username: username,
			Password: password,
		}
		res, err := client.service.SignIn(context.Background(), req)

		if err != nil {
			fmt.Printf("Error while logging in: %v\n", err)
		} else {

			fmt.Println(res)

			break
		}
	}
}

func (client *AuthClient) getCredentials() (string, string) {
	fmt.Print("Username: ")
	client.scanner.Scan()
	username := client.scanner.Text()

	fmt.Print("Password: ")
	client.scanner.Scan()
	password := client.scanner.Text()

	return username, password
}
