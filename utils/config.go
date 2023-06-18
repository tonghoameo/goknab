package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configures of the application
// The values are read by viper from a config or environment variables

type Config struct {
	DBDriver           string        `mapstructure:"DB_DRIVER"`
	DBSource           string        `mapstructure:"DB_URI"`
	ServerAddress      string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymetricKey   string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuraton time.Duration `mapstructure:"ACCESS_TOKEN_DURATON"`
	RefreshTokenDuraton time.Duration `mapstructure:"REFRESH_TOKEN_DURATON"`
}

// LoadConfig  read configuration from file or environment variables

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
