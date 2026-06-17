package config

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	DBHost          string        `mapstructure:"DB_HOST"`
	DBName          string        `mapstructure:"DB_NAME"`
	DBUser          string        `mapstructure:"DB_USER"`
	DBPassword      string        `mapstructure:"DB_PASSWORD"`
	DBPort          string        `mapstructure:"DB_PORT"`
	SecretKey       string        `mapstructure:"SECRET_KEY"`
	Port            string        `mapstructure:"PORT"`
	AccessTokenTTL  time.Duration `mapstructure:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `mapstructure:"REFRESH_TOKEN_TTL"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
	"SECRET_KEY", "PORT", "ACCESS_TOKEN_TTL", "REFRESH_TOKEN_TTL",
}

// LoadConfig loads configuration from the .env file or environment variables.
func LoadConfig() (Config, error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("ACCESS_TOKEN_TTL", "20m")
	viper.SetDefault("REFRESH_TOKEN_TTL", "168h")

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		for _, env := range envs {
			if err := viper.BindEnv(env); err != nil {
				return cfg, err
			}
		}
	}

	if err := viper.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = mapstructure.StringToTimeDurationHookFunc()
	}); err != nil {
		return cfg, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
