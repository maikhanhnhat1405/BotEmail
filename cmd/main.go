package main

import (
	"log"
	"auto-reply-service/internal/config"
	"auto-reply-service/internal/service"
	"auto-reply-service/internal/store"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.DBURL == "" {
		log.Fatal("❌ Missing DB_URL")
	}

	s, err := store.NewPostgresStore(cfg.DBURL)
	if err != nil {
		log.Fatalf("❌ Failed to init store: %v", err)
	}

	// Kiểm tra có account nào chưa
	accounts, err := s.GetActiveAccounts()
	if err != nil {
		log.Fatalf("❌ Failed to get accounts: %v", err)
	}
	if len(accounts) == 0 {
		log.Fatal("❌ No accounts found. Please add account via SQL:")
		log.Fatal(`   INSERT INTO accounts (email, password, imap_host, imap_port, smtp_host, smtp_port, reply_subject, reply_body)
   VALUES ('your_email@gmail.com', 'app_password', 'imap.gmail.com', '993', 'smtp.gmail.com', '587', 'Auto-Reply', 'Hello World');`)
	}

	worker := service.NewWorker(cfg, s)
	worker.Start()
}