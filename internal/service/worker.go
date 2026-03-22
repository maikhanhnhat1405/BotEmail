package service

import (
	"log"
	"time"
	"auto-reply-service/internal/config"
	"auto-reply-service/internal/email"
	"auto-reply-service/internal/store"
	"github.com/emersion/go-imap"
)

type Worker struct {
	cfg   *config.Config
	store *store.PostgresStore
}

func NewWorker(cfg *config.Config, s *store.PostgresStore) *Worker {
	return &Worker{cfg: cfg, store: s}
}

func (w *Worker) Start() {
	log.Println("🚀 Service started. Monitoring mailbox...")
	for {
		w.processAllAccounts()
		time.Sleep(1 * time.Minute)
	}
}

func (w *Worker) processAllAccounts() {
	accounts, err := w.store.GetActiveAccounts()
	if err != nil {
		log.Printf("❌ Failed to get accounts: %v", err)
		return
	}

	if len(accounts) == 0 {
		log.Println("⚠️ No active accounts found")
		return
	}

	log.Printf("📋 Processing %d account(s)...", len(accounts))

	for _, acc := range accounts {
		w.processAccount(acc)
	}
}

func (w *Worker) processAccount(acc store.Account) {
	log.Printf("🔄 Checking inbox: %s", acc.Email)

	c, err := email.ConnectIMAP(acc.IMAPHost, acc.IMAPPort, acc.Email, acc.Password)
	if err != nil {
		log.Printf("❌ IMAP Connection error for %s: %v", acc.Email, err)
		return
	}
	defer c.Logout()

	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Printf("❌ Select INBOX error for %s: %v", acc.Email, err)
		return
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, _ := c.Search(criteria)

	if len(ids) == 0 {
		log.Printf("📭 No unread emails for: %s", acc.Email)
		return
	}

	log.Printf("📬 Found %d unread email(s) for: %s", len(ids), acc.Email)

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)
	messages := make(chan *imap.Message, len(ids))

	go func() {
		c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	for msg := range messages {
		msgID := msg.Envelope.MessageId
		from := msg.Envelope.From[0].Address()

		if w.store.IsProcessed(msgID) {
			log.Printf("⏭️  Skipping already processed: %s", msgID)
			continue
		}

		log.Printf("📩 Sending auto-reply to: %s", from)
		err := email.SendReply(from, acc.ReplySubject, acc.ReplyBody,
			acc.SMTPHost, acc.SMTPPort, acc.Email, acc.Password)

		if err == nil {
			w.store.MarkAsProcessed(msgID)
			log.Printf("✅ Success: Reply sent to %s", from)
		} else {
			log.Printf("❌ Error sending reply: %v", err)
		}
	}
}