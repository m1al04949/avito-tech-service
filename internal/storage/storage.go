package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/m1al04949/avito-tech-service/internal/model"
)

type Storage struct {
	config *Config
	db     *sql.DB
}

var (
	ErrSegmentExists    = errors.New("segment exists")
	ErrSegmentNotExists = errors.New("segment not exists")
	ErrUserExists       = errors.New("user exists")
	ErrUserNotExists    = errors.New("user not exists")
)

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

func (s *Storage) CreateTabs() error {
	const op = "storage.CreateTabs"

	_, err := s.db.Exec(`
	    CREATE TABLE IF NOT EXISTS segments(
		segment_name TEXT NOT NULL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT current_timestamp);
		`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
		user_id INT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT current_timestamp);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS user_segments(
		user_id INT REFERENCES users(user_id),
		segment_name TEXT REFERENCES segments(segment_name),
		PRIMARY KEY (user_id, segment_name));
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Save Segment
func (s *Storage) SaveSegm(segmToSave string) error {
	const op = "storage.SaveSegm"

	m := &model.Segment{
		SegmentName: segmToSave,
	}

	if err := s.db.QueryRow("SELECT (created_at) FROM segments WHERE segment_name=$1",
		segmToSave).Scan(&m.CreatedAt); err != nil {
		stmt, err := s.db.Prepare("INSERT INTO segments(segment_name) VALUES ($1)")
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = stmt.Exec(segmToSave)
		if err != nil {
			if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code == "23505" {
				return fmt.Errorf("%s: %w, created at %s", op, ErrSegmentExists, m.CreatedAt)
			}
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		return fmt.Errorf("%s: %w, created at %s", op, ErrSegmentExists, m.CreatedAt)
	}

	return nil
}

// Delete Segment
func (s *Storage) DeleteSegm(segmToDelete string) error {
	const op = "storage.DeleteSegm"

	m := &model.Segment{
		SegmentName: segmToDelete,
	}

	if err := s.db.QueryRow("SELECT (created_at) FROM segments WHERE segment_name=$1",
		segmToDelete).Scan(&m.CreatedAt); err != nil {
		return fmt.Errorf("%s: %w", op, ErrSegmentNotExists)
	} else {
		stmt, err := s.db.Prepare("DELETE FROM segments WHERE segment_name=$1")
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = stmt.Exec(segmToDelete)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

// Save User
func (s *Storage) SaveUser(userToSave int) error {
	const op = "storage.SaveUser"

	m := &model.User{
		UserID: userToSave,
	}

	if err := s.db.QueryRow("SELECT (created_at) FROM users WHERE user_id=$1",
		userToSave).Scan(&m.CreatedAt); err != nil {
		stmt, err := s.db.Prepare("INSERT INTO users(user_id) VALUES ($1)")
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = stmt.Exec(userToSave)
		if err != nil {
			if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code == "23505" {
				return fmt.Errorf("%s: %w, created at %s", op, ErrUserExists, m.CreatedAt)
			}
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		return fmt.Errorf("%s: %w, created at %s", op, ErrUserExists, m.CreatedAt)
	}

	return nil
}

// Delete User and Segments for him
func (s *Storage) DeleteUser(userToDelete int) error {
	const op = "storage.DeleteUser"

	m := &model.User{
		UserID: userToDelete,
	}

	if err := s.db.QueryRow("SELECT (created_at) FROM users WHERE user_id=$1",
		userToDelete).Scan(&m.CreatedAt); err != nil {
		return fmt.Errorf("%s: %w", op, ErrUserNotExists)
	} else {
		stmt, err := s.db.Prepare("DELETE FROM users WHERE user_id=$1")
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = stmt.Exec(userToDelete)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

// Save Segments for User
func (s *Storage) SaveSegmToUser(userToSave string) error {

	return nil
}
