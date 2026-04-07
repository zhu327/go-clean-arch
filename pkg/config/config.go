package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBPort     string `mapstructure:"DB_PORT"`
	SecretKey  string `mapstructure:"SECRET_KEY"`
	Port       string `mapstructure:"PORT"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
	"SECRET_KEY", "PORT",
}

// LoadConfig loads configuration from the .env file or environment variables.
func LoadConfig() (Config, error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")

	viper.SetDefault("PORT", "8000")

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		for _, env := range envs {
			if err := viper.BindEnv(env); err != nil {
				return cfg, err
			}
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
