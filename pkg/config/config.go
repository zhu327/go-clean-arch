package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration.
type Config struct {
	DBHost          string
	DBName          string
	DBUser          string
	DBPassword      string
	DBPort          string
	DBTimezone      string
	SecretKey       string
	Port            string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", "DB_TIMEZONE",
	"SECRET_KEY", "PORT", "ACCESS_TOKEN_TTL", "REFRESH_TOKEN_TTL",
}

// LoadConfig loads configuration from .env and process environment variables.
// Process environment variables take precedence over values in .env.
func LoadConfig() (Config, error) {
	return loadConfig(".env", os.LookupEnv)
}

func loadConfig(path string, lookupEnv func(string) (string, bool)) (Config, error) {
	values, err := readDotEnv(path)
	if err != nil {
		return Config{}, err
	}
	for _, key := range envs {
		if value, ok := lookupEnv(key); ok {
			values[key] = value
		}
	}
	if values["PORT"] == "" {
		values["PORT"] = "8000"
	}
	if values["ACCESS_TOKEN_TTL"] == "" {
		values["ACCESS_TOKEN_TTL"] = "20m"
	}
	if values["REFRESH_TOKEN_TTL"] == "" {
		values["REFRESH_TOKEN_TTL"] = "168h"
	}
	if values["DB_TIMEZONE"] == "" {
		values["DB_TIMEZONE"] = "Asia/Shanghai"
	}
	return configFromValues(values)
}

func readDotEnv(path string) (map[string]string, error) {
	values := make(map[string]string)
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return values, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open .env: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for line := 1; scanner.Scan(); line++ {
		text := strings.TrimSpace(scanner.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		key, value, ok := strings.Cut(text, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid .env line %d", line)
		}
		values[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), "\"'")
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read .env: %w", err)
	}
	return values, nil
}

func configFromValues(values map[string]string) (Config, error) {
	accessTTL, err := time.ParseDuration(values["ACCESS_TOKEN_TTL"])
	if err != nil || accessTTL <= 0 {
		return Config{}, fmt.Errorf("ACCESS_TOKEN_TTL must be a positive duration")
	}
	refreshTTL, err := time.ParseDuration(values["REFRESH_TOKEN_TTL"])
	if err != nil || refreshTTL <= accessTTL {
		return Config{}, fmt.Errorf("REFRESH_TOKEN_TTL must be greater than ACCESS_TOKEN_TTL")
	}
	for _, key := range []string{"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD"} {
		if strings.TrimSpace(values[key]) == "" {
			return Config{}, fmt.Errorf("%s is required", key)
		}
	}
	if len([]byte(values["SECRET_KEY"])) < 32 {
		return Config{}, fmt.Errorf("SECRET_KEY must be at least 32 bytes")
	}
	if err := validatePort("DB_PORT", values["DB_PORT"]); err != nil {
		return Config{}, err
	}
	if _, err := time.LoadLocation(values["DB_TIMEZONE"]); err != nil {
		return Config{}, fmt.Errorf("DB_TIMEZONE must be a valid IANA timezone: %w", err)
	}
	if err := validatePort("PORT", values["PORT"]); err != nil {
		return Config{}, err
	}
	return Config{
		DBHost: values["DB_HOST"], DBName: values["DB_NAME"], DBUser: values["DB_USER"],
		DBPassword: values["DB_PASSWORD"], DBPort: values["DB_PORT"], DBTimezone: values["DB_TIMEZONE"], SecretKey: values["SECRET_KEY"],
		Port: values["PORT"], AccessTokenTTL: accessTTL, RefreshTokenTTL: refreshTTL,
	}, nil
}

func validatePort(name, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("%s must be between 1 and 65535", name)
	}
	return nil
}
