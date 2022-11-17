package service

import (
	"context"

	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	"google.golang.org/grpc/codes"
	glog "google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

func NewAuthServer(grpcLog glog.LoggerV2) *AuthServer {
	return &AuthServer{
		authpb.UnimplementedAuthServiceServer{},
		grpcLog,
	}
}

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	grpcLog glog.LoggerV2
}

func (s *AuthServer) SignIn(ctx context.Context, req *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	user := req.GetUsername()
	s.grpcLog.Info("Attempt to log in with: ", user)
	if user == "Anon" {
		s.grpcLog.Info("anon logged in")
		res := &authpb.SignInResponse{
			Status:       "success",
			AccessToken:  "placeholder",
			RefreshToken: "placeholder",
		}

		return res, nil
	}

	return nil, status.Errorf(codes.NotFound, "Not found, login failed")

}
