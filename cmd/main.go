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

	// Thêm account mặc định nếu chưa có
	accounts, _ := s.GetActiveAccounts()
	if len(accounts) == 0 {
		log.Println("📝 No accounts found. Adding default account...")
		err := s.AddAccount(
			"quaniudau812006@gmail.com",
			"hhzpnjkljchiixnr",
			"imap.gmail.com",
			"993",
			"smtp.gmail.com",
			"587",
			"Auto-Reply",
			"Hello World",
		)
		if err != nil {
			log.Fatalf("❌ Failed to add default account: %v", err)
		}
		log.Println("✅ Default account added!")
	}

	worker := service.NewWorker(cfg, s)
	worker.Start()
}