package service

import (
	"context"

	// "github.com/MyyPo/grpc-chat/internal/repositories"
	"github.com/MyyPo/grpc-chat/internal/util"
	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	"google.golang.org/grpc/codes"
	glog "google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

func NewAuthServer(grpcLog glog.LoggerV2, tokenManager util.TokenManager) AuthServer {
	return AuthServer{
		authpb.UnimplementedAuthServiceServer{},
		&tokenManager,
		grpcLog,
	}
}

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	tokenManager *util.TokenManager
	grpcLog      glog.LoggerV2
}

func (s *AuthServer) SignUp(ctx context.Context, req *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {
	// !TODO
	user := req.GetUsername()
	password := req.GetPassword()
	s.grpcLog.Info("Sign up attempt with: ", user)

	return &authpb.SignUpResponse{
		AccessToken:  user,
		RefreshToken: password,
	}, nil
}

func (s *AuthServer) SignIn(ctx context.Context, req *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	// user := req.GetUsername()
	// password := req.GetPassword()
	// s.grpcLog.Info("Attempt to log in with: ", user)

	// repositories.DBAuth.SignIn( ctx, req)

	// if user == "Anon" && password == "Cute" {
	// 	s.grpcLog.Infof("User logged in: %s", user)
	// 	accessToken, _ := s.tokenManager.GenerateJWT(true)
	// 	refreshToken, _ := s.tokenManager.GenerateJWT(false)
	// 	res := &authpb.SignInResponse{
	// 		AccessToken:  accessToken,
	// 		RefreshToken: refreshToken,
	// 	}

	// 	return res, nil
	// }

	return nil, status.Errorf(codes.NotFound, "Not found, login failed")

}

func (s *AuthServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	tokenPackage, err := s.tokenManager.RerfreshToken(req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &tokenPackage, nil
}
