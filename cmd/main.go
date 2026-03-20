package main

import (
	"log"
	"auto-reply-service/internal/config"
	"auto-reply-service/internal/service"
	"auto-reply-service/internal/store"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.EmailUser == "" || cfg.EmailPass == "" {
		log.Fatal("❌ Missing EMAIL_USER or EMAIL_PASS in environment variables")
	}

	s, err := store.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		log.Fatalf("❌ Failed to init store: %v", err)
	}

	worker := service.NewWorker(cfg, s)
	worker.Start()
}