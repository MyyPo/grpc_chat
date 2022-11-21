package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AccessSignature      string        `mapstructure:"ACCESS_SIGNATURE"`
	RefreshSignature     string        `mapstructure:"REFRESH_SIGNATURE"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func NewConfig(path string) (*Config, error) {
	newConf := &Config{}

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
