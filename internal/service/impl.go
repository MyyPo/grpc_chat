package service

import (
	"github.com/MyyPo/grpc-chat/internal/util"
	glog "google.golang.org/grpc/grpclog"
)

type Implementation struct {
	AuthServer
	ChatServer
	util.TokenManager
}

func NewImplementation(grpcLogger glog.LoggerV2, jwtSignature string) Implementation {
	return Implementation{
		NewAuthServer(grpcLogger),
		NewChatServer(grpcLogger),
		util.NewTokenManager(jwtSignature),
	}
}
