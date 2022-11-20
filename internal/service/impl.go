package service

import (
	"github.com/MyyPo/grpc-chat/internal/util"
	glog "google.golang.org/grpc/grpclog"
	"os"
)

type Implementation struct {
	AuthServer
	ChatServer
	AuthInterceptor
}

func NewImplementation(config util.Config, accessibleRoles map[string][]string) Implementation {
	grpcLogger := glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
	tokenManager := util.NewTokenManager(config.AccessSignature, config.RefreshSignature, config.AccessTokenDuration, config.RefreshTokenDuration)

	return Implementation{
		NewAuthServer(grpcLogger, tokenManager),
		NewChatServer(grpcLogger),
		NewAuthInterceptor(tokenManager, accessibleRoles),
	}
}
