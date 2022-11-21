package service

import (
	"context"
	"fmt"

	"github.com/MyyPo/grpc-chat/internal/repositories"
	"github.com/MyyPo/grpc-chat/internal/util"
	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	glog "google.golang.org/grpc/grpclog"
)

func NewAuthServer(grpcLog glog.LoggerV2,
	tokenManager util.TokenManager,
	hasher Hasher,
	authRepo repositories.DBAuth,
) AuthServer {
	return AuthServer{
		authpb.UnimplementedAuthServiceServer{},
		&tokenManager,
		hasher,
		authRepo,
		grpcLog,
	}
}

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	tokenManager *util.TokenManager
	hasher       Hasher
	authRepo     repositories.Auth
	grpcLog      glog.LoggerV2
}

func (s *AuthServer) SignUp(ctx context.Context, req *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {
	username := req.GetUsername()
	s.grpcLog.Info("Sign up attempt with: ", username)

	req.Password, _ = s.hasher.Hash(req.GetPassword())

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
	username := req.GetUsername()
	s.grpcLog.Info("Attempt to log in with: ", username)

	res, err := s.authRepo.SignIn(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("invalid password or username")
	}
	if !s.hasher.IsValid(req.GetPassword(), res.Password) {
		return nil, fmt.Errorf("invalid password or username")
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
