package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"link-shortener/internal/storage"
	_ "modernc.org/sqlite"
)

type Storage struct {
	DB *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s (opening database): %w", op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
		    id INTEGER PRIMARY KEY,
		    alias TEXT NOT NULL UNIQUE,
		    url TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s (creating table): %w", op, err)
	}

	// Create an index on the `alias` column if it does not exist
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_alias ON links(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s (creating index): %w", op, err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) SaveURL(URL string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveLink"
	stmt, err := s.DB.Prepare("INSERT INTO links (url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExist)
	}

	res, err := stmt.Exec(URL, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get id %w", op, err)
	}

	return int64(id), nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetLink"

	stmt, err := s.DB.Prepare("SELECT url FROM links WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteURL(urlID int64) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.DB.Prepare("DELETE FROM links WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if n == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}
