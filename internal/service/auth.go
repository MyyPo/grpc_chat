package service

import (
	"context"

	"github.com/MyyPo/grpc-chat/internal/repositories"
	"github.com/MyyPo/grpc-chat/internal/util"
	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	glog "google.golang.org/grpc/grpclog"
)

func NewAuthServer(grpcLog glog.LoggerV2, tokenManager util.TokenManager, authRepo repositories.DBAuth) AuthServer {
	return AuthServer{
		authpb.UnimplementedAuthServiceServer{},
		&tokenManager,
		authRepo,
		grpcLog,
	}
}

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	tokenManager *util.TokenManager
	authRepo     repositories.Auth
	grpcLog      glog.LoggerV2
}

func (s *AuthServer) SignUp(ctx context.Context, req *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {
	// !TODO
	user := req.GetUsername()
	s.grpcLog.Info("Sign up attempt with: ", user)

	res, err := s.authRepo.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.tokenManager.GenerateJWT(true, res.UserID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokenManager.GenerateJWT(false, res.UserID)
	if err != nil {
		return nil, err
	}

	return &authpb.SignUpResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServer) SignIn(ctx context.Context, req *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	user := req.GetUsername()
	s.grpcLog.Info("Attempt to log in with: ", user)

	res, err := s.authRepo.SignIn(ctx, req)
	if err != nil {
		return nil, err
	}

	accessToken, _ := s.tokenManager.GenerateJWT(true, res.UserID)
	if err != nil {
		return nil, err
	}
	refreshToken, _ := s.tokenManager.GenerateJWT(false, res.UserID)
	if err != nil {
		return nil, err
	}

	return &authpb.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *AuthServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	tokenPackage, err := s.tokenManager.RerfreshToken(req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &tokenPackage, nil
}
