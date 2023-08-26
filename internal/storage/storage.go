package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage struct {
	config *Config
	db     *sql.DB
}

// Get instance
func New(cfgpath, dburl string) *Storage {
	return &Storage{
		config: &Config{
			ConfigPath:  cfgpath,
			DatabaseURL: dburl,
		},
	}
}

// Open connection to DB
func (s *Storage) Open() error {

	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

// Close connection
func (s *Storage) Close() {
	s.db.Close()
}
