package util

import (
	"github.com/spf13/viper"
)

type config struct {
	JWTSignature string `mapstructure:"JWT_SIGNATURE"`
}

func NewConfig(path string) (*config, error) {
	newConf := &config{}

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(newConf); err != nil {
		return nil, err
	}

	return newConf, nil
}
