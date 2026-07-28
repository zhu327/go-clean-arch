package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func validValues() map[string]string {
	return map[string]string{
		"DB_HOST":           "db",
		"DB_NAME":           "app",
		"DB_USER":           "app",
		"DB_PASSWORD":       "password",
		"DB_PORT":           "5432",
		"SECRET_KEY":        "12345678901234567890123456789012",
		"PORT":              "8000",
		"ACCESS_TOKEN_TTL":  "20m",
		"REFRESH_TOKEN_TTL": "168h",
	}
}

func TestLoadConfigFromFile_EnvironmentOverridesDotEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("DB_HOST=file-db\nDB_NAME=app\nDB_USER=app\nDB_PASSWORD=password\nDB_PORT=5432\nSECRET_KEY=12345678901234567890123456789012\nPORT=8000\nACCESS_TOKEN_TTL=20m\nREFRESH_TOKEN_TTL=168h\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig(path, func(key string) (string, bool) {
		if key == "DB_HOST" {
			return "environment-db", true
		}
		return "", false
	})
	if err != nil {
		t.Fatalf("loadConfig returned error: %v", err)
	}
	if cfg.DBHost != "environment-db" {
		t.Errorf("DBHost = %q, want environment-db", cfg.DBHost)
	}
}

func TestLoadConfigFromFile_RejectsInvalidConfiguration(t *testing.T) {
	for name, change := range map[string]func(map[string]string){
		"missing database host":               func(v map[string]string) { v["DB_HOST"] = "" },
		"database port out of range":          func(v map[string]string) { v["DB_PORT"] = "65536" },
		"weak secret":                         func(v map[string]string) { v["SECRET_KEY"] = "too-short" },
		"server port out of range":            func(v map[string]string) { v["PORT"] = "0" },
		"non-positive access TTL":             func(v map[string]string) { v["ACCESS_TOKEN_TTL"] = "0s" },
		"refresh TTL shorter than access TTL": func(v map[string]string) { v["REFRESH_TOKEN_TTL"] = "10m" },
	} {
		t.Run(name, func(t *testing.T) {
			values := validValues()
			change(values)
			if _, err := configFromValues(values); err == nil {
				t.Fatal("configFromValues succeeded, want validation error")
			}
		})
	}
}

func TestConfigFromValues_ParsesValidConfiguration(t *testing.T) {
	cfg, err := configFromValues(validValues())
	if err != nil {
		t.Fatalf("configFromValues returned error: %v", err)
	}
	if cfg.AccessTokenTTL != 20*time.Minute || cfg.RefreshTokenTTL != 168*time.Hour {
		t.Errorf("unexpected TTLs: access=%v refresh=%v", cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	}
}
