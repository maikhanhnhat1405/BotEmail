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

	s := store.NewJSONStore(cfg.StorePath)
	worker := service.NewWorker(cfg, s)

	worker.Start()
}