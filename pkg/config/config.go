package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBPort     string `mapstructure:"DB_PORT"`
	SECRET_KEY string `mapstructure:"SECRET_KEY"`
}

// LoadConfig loads the configuration from the .env file
var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
	"SECRET_KEY",
}

func LoadConfig() (config Config, err error) {
	// read .env file
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	if err != nil {
		// binding
		for _, env := range envs {
			if err := viper.BindEnv(env); err != nil {
				return config, err
			}
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(config); err != nil {
		return config, err
	}

	return config, nil
}
