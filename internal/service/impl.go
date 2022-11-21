package service

import (
	"os"

	"github.com/MyyPo/grpc-chat/internal/config"
	"github.com/MyyPo/grpc-chat/internal/repositories"
	"github.com/MyyPo/grpc-chat/internal/util"
	glog "google.golang.org/grpc/grpclog"
)

type Implementation struct {
	AuthServer
	ChatServer
	AuthInterceptor
}

func NewImplementation(config config.Config, accessibleRoles map[string][]string, authRepo repositories.DBAuth) Implementation {
	grpcLogger := glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
	tokenManager := util.NewTokenManager(config.AccessSignature, config.RefreshSignature, config.AccessTokenDuration, config.RefreshTokenDuration)
	hasher := NewHasher()

	return Implementation{
		NewAuthServer(grpcLogger, tokenManager, hasher, authRepo),
		NewChatServer(grpcLogger),
		NewAuthInterceptor(tokenManager, accessibleRoles),
	}
}
