package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configures of the application
// The values are read by viper from a config or environment variables

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	//DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_URI"`
	MigrationURL         string        `mapstructure:"MIGRATION_PATH"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	HTTPGinServerAddress string        `mapstructure:"HTTP_GIN_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	TokenSymetricKey     string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuraton   time.Duration `mapstructure:"ACCESS_TOKEN_DURATON"`
	RefreshTokenDuraton  time.Duration `mapstructure:"REFRESH_TOKEN_DURATON"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
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
