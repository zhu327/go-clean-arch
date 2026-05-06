package config

import (
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func TestConfig_DefaultTokenTTLs(t *testing.T) {
	viper.Reset()
	viper.SetDefault("ACCESS_TOKEN_TTL", "20m")
	viper.SetDefault("REFRESH_TOKEN_TTL", "168h")

	var cfg Config
	err := viper.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = mapstructure.StringToTimeDurationHookFunc()
	})
	if err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	expectedAccess := 20 * time.Minute
	expectedRefresh := 168 * time.Hour

	if cfg.AccessTokenTTL != expectedAccess {
		t.Errorf("expected AccessTokenTTL=%v, got %v", expectedAccess, cfg.AccessTokenTTL)
	}
	if cfg.RefreshTokenTTL != expectedRefresh {
		t.Errorf("expected RefreshTokenTTL=%v, got %v", expectedRefresh, cfg.RefreshTokenTTL)
	}
}

func TestConfig_ParsesTokenTTLs(t *testing.T) {
	viper.Reset()
	viper.Set("ACCESS_TOKEN_TTL", "20m")
	viper.Set("REFRESH_TOKEN_TTL", "168h")

	var cfg Config
	err := viper.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = mapstructure.StringToTimeDurationHookFunc()
	})
	if err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if cfg.AccessTokenTTL != 20*time.Minute {
		t.Errorf("expected AccessTokenTTL=20m, got %v", cfg.AccessTokenTTL)
	}
	if cfg.RefreshTokenTTL != 168*time.Hour {
		t.Errorf("expected RefreshTokenTTL=168h, got %v", cfg.RefreshTokenTTL)
	}
}
