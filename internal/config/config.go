package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string // PostgreSQL: postgresql://user:pass@host:port/db
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		DBURL: os.Getenv("DB_URL"),
	}
}