package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

type Account struct {
	ID           int
	Email        string
	Password     string
	IMAPHost     string
	IMAPPort     string
	SMTPHost     string
	SMTPPort     string
	ReplySubject string
	ReplyBody    string
	Active       bool
}

func NewPostgresStore(dbURL string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Tạo bảng nếu chưa có
	queries := []string{
		`CREATE TABLE IF NOT EXISTS processed_emails (id TEXT PRIMARY KEY);`,
		`CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			imap_host TEXT DEFAULT 'imap.gmail.com',
			imap_port TEXT DEFAULT '993',
			smtp_host TEXT DEFAULT 'smtp.gmail.com',
			smtp_port TEXT DEFAULT '587',
			reply_subject TEXT DEFAULT 'Auto-Reply',
			reply_body TEXT DEFAULT 'Hello World',
			active BOOLEAN DEFAULT true
		);`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return nil, err
		}
	}

	return &PostgresStore{db: db}, nil
}

// Lấy danh sách accounts đang active
func (s *PostgresStore) GetActiveAccounts() ([]Account, error) {
	rows, err := s.db.Query(`
		SELECT id, email, password, imap_host, imap_port, smtp_host, smtp_port, reply_subject, reply_body, active
		FROM accounts WHERE active = true
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var a Account
		err := rows.Scan(&a.ID, &a.Email, &a.Password, &a.IMAPHost, &a.IMAPPort, &a.SMTPHost, &a.SMTPPort, &a.ReplySubject, &a.ReplyBody, &a.Active)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

// Thêm account mới
func (s *PostgresStore) AddAccount(email, password, imapHost, imapPort, smtpHost, smtpPort, replySubject, replyBody string) error {
	_, err := s.db.Exec(`
		INSERT INTO accounts (email, password, imap_host, imap_port, smtp_host, smtp_port, reply_subject, reply_body)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, email, password, imapHost, imapPort, smtpHost, smtpPort, replySubject, replyBody)
	return err
}

// Xóa account
func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM accounts WHERE id = $1", id)
	return err
}

// Kiểm tra email đã xử lý chưa
func (s *PostgresStore) IsProcessed(id string) bool {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM processed_emails WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// Đánh dấu email đã xử lý
func (s *PostgresStore) MarkAsProcessed(id string) error {
	_, err := s.db.Exec("INSERT INTO processed_emails (id) VALUES ($1)", id)
	return err
}
