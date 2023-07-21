package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/vcholak/messenger-bot/storage"
)

type SqliteStorage struct {
	db *sql.DB
}

// New creates a new SQLite storage.
func New(path string) (*SqliteStorage, error) {

	sqlite, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := sqlite.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	storage := &SqliteStorage{db: sqlite}
	return storage, nil
}

// Save saves a page to the storage.
func (s *SqliteStorage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, first_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.FirstName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	return nil
}

// PickRandom picks a random page from the storage.
func (s *SqliteStorage) PickRandom(ctx context.Context, firstName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE first_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, firstName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:       url,
		FirstName: firstName,
	}, nil
}

// Remove removes a page from the storage.
func (s *SqliteStorage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND first_name = ?`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.FirstName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}
	return nil
}

// IsExists checks if a page exists in the storage.
func (s *SqliteStorage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND first_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.FirstName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}
	return count > 0, nil
}

// Init creates the pages table if it doesn't exist yet.
func (s *SqliteStorage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, first_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create 'pages' table: %w", err)
	}
	return nil
}
