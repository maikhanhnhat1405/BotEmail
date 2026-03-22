package main

import (
	"flag"
	"fmt"
	"log"
	"auto-reply-service/internal/config"
	"auto-reply-service/internal/service"
	"auto-reply-service/internal/store"
)

func main() {
	// Register commands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	email := addCmd.String("email", "", "Email address")
	password := addCmd.String("password", "", "App password")
	replySubject := addCmd.String("reply-subject", "Auto-Reply", "Reply subject")
	replyBody := addCmd.String("reply-body", "Hello World", "Reply body")

	if len(flag.Args()) == 0 {
		runService()
		return
	}

	switch flag.Args()[0] {
	case "add":
		addCmd.Parse(flag.Args()[1:])
		if *email == "" || *password == "" {
			log.Fatal("❌ --email and --password are required")
		}
		addAccount(*email, *password, *replySubject, *replyBody)
	default:
		log.Fatalf("❌ Unknown command: %s", flag.Args()[0])
	}
}

func runService() {
	cfg := config.LoadConfig()
	if cfg.DBURL == "" {
		log.Fatal("❌ Missing DB_URL")
	}

	s, err := store.NewPostgresStore(cfg.DBURL)
	if err != nil {
		log.Fatalf("❌ Failed to init store: %v", err)
	}

	accounts, err := s.GetActiveAccounts()
	if err != nil {
		log.Fatalf("❌ Failed to get accounts: %v", err)
	}
	if len(accounts) == 0 {
		log.Fatal("❌ No accounts found. Run: go run ./cmd/main.go add --email=... --password=...")
	}

	worker := service.NewWorker(cfg, s)
	worker.Start()
}

func addAccount(email, password, replySubject, replyBody string) {
	cfg := config.LoadConfig()
	if cfg.DBURL == "" {
		log.Fatal("❌ Missing DB_URL")
	}

	s, err := store.NewPostgresStore(cfg.DBURL)
	if err != nil {
		log.Fatalf("❌ Failed to init store: %v", err)
	}

	err = s.AddAccount(
		email,
		password,
		"imap.gmail.com",
		"993",
		"smtp.gmail.com",
		"587",
		replySubject,
		replyBody,
	)
	if err != nil {
		log.Fatalf("❌ Failed to add account: %v", err)
	}

	fmt.Printf("✅ Account added: %s\n", email)
}
