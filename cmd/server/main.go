package main

import (
	"fmt"
	"log"
	"net"
	"os"

	authpb "github.com/MyyPo/grpc-chat/pb/auth/v1"
	chatpb "github.com/MyyPo/grpc-chat/pb/chat/v1"
	"github.com/MyyPo/grpc-chat/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2
var config Config

type Config struct {
	JWTSignature string `mapstructure:"JWT_SIGNATURE"`
}

func (config *Config) LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return err
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

func main() {
	config.LoadConfig("./../..")
	fmt.Println(config.JWTSignature)

	chatServer := service.NewChatServer(grpcLog)
	authServer := service.NewAuthServer(grpcLog)

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to create the server %v", err)
	}

	grpcLog.Info("The server successfuly started")

	chatpb.RegisterBroadcastServiceServer(grpcServer, chatServer)
	authpb.RegisterAuthServiceServer(grpcServer, authServer)
	grpcServer.Serve(lis)
}
