package store

import (
	"database/sql"
	_ "modernc.org/sqlite" // Driver SQLite thuần Go (không cần CGO)
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Tạo bảng lưu ID email nếu chưa có
	query := `CREATE TABLE IF NOT EXISTS processed_emails (id TEXT PRIMARY KEY);`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) IsProcessed(id string) bool {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM processed_emails WHERE id=?)", id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// Lưu ID email vào DB
func (s *SQLiteStore) MarkAsProcessed(id string) error {
	_, err := s.db.Exec("INSERT INTO processed_emails (id) VALUES (?)", id)
	return err
}