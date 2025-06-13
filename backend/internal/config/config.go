package config

import (
	"log"
	"os"
)

type Config struct {
	Address string
}

func Load() Config {
	cfg := Config{
		Address: getEnv("APP_ADDRESS", ":8080"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Printf("%s not set, using default %s", key, fallback)
	return fallback
}
