package service

import (
	"github.com/MyyPo/grpc-chat/internal/util"
	glog "google.golang.org/grpc/grpclog"
	"os"
)

type Implementation struct {
	AuthServer
	ChatServer
	util.TokenManager
}

func NewImplementation(tokenManager util.TokenManager) Implementation {
	grpcLogger := glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)

	return Implementation{
		NewAuthServer(grpcLogger),
		NewChatServer(grpcLogger),
		tokenManager,
	}
}
