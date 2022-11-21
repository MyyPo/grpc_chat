package client_service

import (
	"bufio"
	"context"
	"fmt"
	"time"

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

func (client *AuthClient) SignIn() (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	username, password := client.getCredentials()
	// username, password := "TestingHasher", "TestingHasher"

	req := &authpb.SignInRequest{
		Username: username,
		Password: password,
	}
	res, err := client.service.SignIn(ctx, req)

	if err != nil {
		return "", fmt.Errorf("error while logging in: %v", err)
	} else {

		// fmt.Println(res)
		// time.Sleep(time.Second * 1)
		// refresh, err := client.RefreshToken(res.RefreshToken)
		// if err != nil {
		// 	fmt.Printf("Error while trying to refresh the token: %v\n", err)
		// }
		// fmt.Println(refresh)
		return res.GetAccessToken(), nil
	}
}

// }

func (client *AuthClient) SignUp(username, password string) (*authpb.SignUpResponse, error) {
	req := &authpb.SignUpRequest{
		Username: username,
		Password: password,
	}
	res, err := client.service.SignUp(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *AuthClient) RefreshToken(refreshToken string) (*authpb.RefreshTokenResponse, error) {
	req := &authpb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}
	res, err := client.service.RefreshToken(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
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
