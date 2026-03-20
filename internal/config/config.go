package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	IMAPHost     string
	IMAPPort     string
	SMTPHost     string
	SMTPPort     string
	EmailUser    string
	EmailPass    string
	ReplySubject string
	ReplyBody    string
	DBPath       string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		IMAPHost:     getEnv("IMAP_HOST", "imap.gmail.com"),
		IMAPPort:     getEnv("IMAP_PORT", "993"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		EmailUser:    os.Getenv("EMAIL_USER"),
		EmailPass:    os.Getenv("EMAIL_PASS"),
		ReplySubject: getEnv("REPLY_SUBJECT", "Auto-Reply: Acknowledgement"),
		ReplyBody:    getEnv("REPLY_BODY", "Hello World"),
		DBPath:       getEnv("DB_PATH", "data/emails.db"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}