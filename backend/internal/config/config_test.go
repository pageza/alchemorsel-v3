package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")

	cfg := Load()

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host 0.0.0.0, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("expected default db host localhost, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected default db port 5432, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "postgres" {
		t.Errorf("expected default db user postgres, got %s", cfg.Database.User)
	}
	if cfg.Database.Password != "" {
		t.Errorf("expected default db password empty, got %s", cfg.Database.Password)
	}
	if cfg.Database.Name != "alchemorsel" {
		t.Errorf("expected default db name alchemorsel, got %s", cfg.Database.Name)
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("expected default sslmode disable, got %s", cfg.Database.SSLMode)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	os.Setenv("APP_HOST", "1.2.3.4")
	os.Setenv("APP_PORT", "9000")
	os.Setenv("DB_HOST", "dbhost")
	os.Setenv("DB_PORT", "notint")
	os.Setenv("DB_USER", "bob")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_NAME", "mydb")
	os.Setenv("DB_SSLMODE", "require")
	defer func() {
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
	}()

	cfg := Load()

	if cfg.Server.Host != "1.2.3.4" {
		t.Errorf("expected host 1.2.3.4, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("expected port 9000, got %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "dbhost" {
		t.Errorf("expected db host dbhost, got %s", cfg.Database.Host)
	}
	// DB_PORT is not int so should fall back to default 5432
	if cfg.Database.Port != 5432 {
		t.Errorf("expected db port 5432 due to fallback, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "bob" {
		t.Errorf("expected db user bob, got %s", cfg.Database.User)
	}
	if cfg.Database.Password != "secret" {
		t.Errorf("expected db password secret, got %s", cfg.Database.Password)
	}
	if cfg.Database.Name != "mydb" {
		t.Errorf("expected db name mydb, got %s", cfg.Database.Name)
	}
	if cfg.Database.SSLMode != "require" {
		t.Errorf("expected sslmode require, got %s", cfg.Database.SSLMode)
	}
}
