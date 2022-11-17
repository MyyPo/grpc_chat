package utils

import (
	"fmt"
	"github.com/spf13/viper"
	glog "google.golang.org/grpc/grpclog"
	"os"
)

type Config struct {
	JWTSignature string `mapstructure:"JWT_SIGNATURE"`

	grpcLog glog.LoggerV2
}

func NewConfig(path string) (*Config, error) {
	var config Config
	config.grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &config, nil
}
