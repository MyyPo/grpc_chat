package service

import (
	"context"

	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	"google.golang.org/grpc/codes"
	glog "google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

func NewAuthServer(grpcLog glog.LoggerV2) AuthServer {
	return AuthServer{
		authpb.UnimplementedAuthServiceServer{},
		grpcLog,
	}
}

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	grpcLog glog.LoggerV2
}

func (i *Implementation) SignIn(ctx context.Context, req *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	user := req.GetUsername()
	password := req.GetPassword()
	i.ChatServer.grpcLog.Info("Attempt to log in with: ", user)
	if user == "Anon" && password == "Cute" {
		i.ChatServer.grpcLog.Infof("User logged in: %s", user)
		accessToken, _ := i.TokenManager.GenerateJWT(true)
		refreshToken, _ := i.TokenManager.GenerateJWT(false)
		res := &authpb.SignInResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return res, nil
	}

	return nil, status.Errorf(codes.NotFound, "Not found, login failed")

}

func (i *Implementation) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	tokenPackage, err := i.TokenManager.RerfreshToken(req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &tokenPackage, nil
}
