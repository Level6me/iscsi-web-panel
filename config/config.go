package config

import (
	"os"
)

type Config struct {
	ListenAddr string
	JWTSecret  string
	DataDir    string
	DBPath     string
}

func Load() *Config {
	return &Config{
		ListenAddr: getEnv("LISTEN_ADDR", ":3005"),
		JWTSecret:  getEnv("JWT_SECRET", "iscsi-web-panel-secret-key"),
		DataDir:    getEnv("DATA_DIR", "./data"),
		DBPath:     getEnv("DB_PATH", "./data/iscsi.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
