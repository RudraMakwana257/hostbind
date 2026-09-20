package registry

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

type Registry struct {
	db *sql.DB
}

func Open() (*Registry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home dir: %w", err)
	}

	dir := filepath.Join(home, ".hostbind")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("could not create hostbind dir: %w", err)
	}

	dbPath := filepath.Join(dir, "registry.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	r := &Registry{db: db}
	if err := r.migrate(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Registry) Close() error {
	return r.db.Close()
}
