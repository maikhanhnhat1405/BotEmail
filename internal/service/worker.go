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
	store *store.JSONStore
}

func NewWorker(cfg *config.Config, s *store.JSONStore) *Worker {
	return &Worker{cfg: cfg, store: s}
}

func (w *Worker) Start() {
	log.Println("🚀 Service started. Monitoring mailbox...")
	for {
		w.process()
		time.Sleep(1 * time.Minute)
	}
}

func (w *Worker) process() {
	c, err := email.ConnectIMAP(w.cfg.IMAPHost, w.cfg.IMAPPort, w.cfg.EmailUser, w.cfg.EmailPass)
	if err != nil {
		log.Printf("❌ IMAP Connection error: %v", err)
		return
	}
	defer c.Logout()

	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Printf("❌ Select INBOX error: %v", err)
		return
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, _ := c.Search(criteria)

	if len(ids) == 0 {
		return
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)
	messages := make(chan *imap.Message, len(ids))
	
	go func() {
		c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	for msg := range messages {
		msgID := msg.Envelope.MessageId
		from := msg.Envelope.From[0].Address()

		if w.store.Exists(msgID) {
			continue
		}

		log.Printf("📩 Sending auto-reply to: %s", from)
		err := email.SendReply(from, w.cfg.ReplySubject, w.cfg.ReplyBody,
			w.cfg.SMTPHost, w.cfg.SMTPPort, w.cfg.EmailUser, w.cfg.EmailPass)

		if err == nil {
			w.store.Save(msgID)
			log.Printf("✅ Success: Reply sent to %s", from)
		} else {
			log.Printf("❌ Error sending reply: %v", err)
		}
	}
}