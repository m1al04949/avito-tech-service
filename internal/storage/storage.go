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
	ErrSegmentExists     = errors.New("segment exists")
	ErrSegmentNotExists  = errors.New("segment not exists")
	ErrSegmentsNotExists = errors.New("segments not exists")
	ErrUserExists        = errors.New("user exists")
	ErrUserNotExists     = errors.New("user not exists")
	ErrUserDelete        = errors.New("delete user from user segments table")
	ErrSegmentDelete     = errors.New("delete segment from user segments table")
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

	m := &model.Segments{
		SegmentName: segmToSave,
	}

	if err := s.db.QueryRow("SELECT (created_at) FROM segments WHERE segment_name=$1",
		segmToSave).Scan(&m.CreatedAt); err != nil {
		stmt, err := s.db.Prepare("INSERT INTO segments(segment_name) VALUES ($1)")
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		defer stmt.Close()

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

	m := &model.Segments{
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
		defer stmt.Close()

		_, err = stmt.Exec(segmToDelete)
		if err != nil {
			return fmt.Errorf("%s: %w", op, ErrSegmentDelete)
		}
	}

	return nil
}

// Save User
func (s *Storage) SaveUser(userToSave int) error {
	const op = "storage.SaveUser"

	m := &model.Users{
		UserID: userToSave,
	}

	if err := s.db.QueryRow("SELECT (created_at) FROM users WHERE user_id=$1",
		userToSave).Scan(&m.CreatedAt); err != nil {
		stmt, err := s.db.Prepare("INSERT INTO users(user_id) VALUES ($1)")
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		defer stmt.Close()

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

// Delete User
func (s *Storage) DeleteUser(userToDelete int) error {
	const op = "storage.DeleteUser"

	m := &model.Users{
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
		defer stmt.Close()

		_, err = stmt.Exec(userToDelete)
		if err != nil {
			return fmt.Errorf("%s: %w", op, ErrUserDelete)
		}
	}

	return nil
}

// Save Segments for User
func (s *Storage) SaveSegmToUser(user int, segments []string) error {
	const op = "storage.AddToUser"

	var userExists bool

	if err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)",
		user).Scan(&userExists); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !userExists {
		return fmt.Errorf("%s: %w", op, ErrUserNotExists)
	}

	existingSegments := make([]string, 0, len(segments))
	for _, v := range segments {
		var segmentExists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM segments WHERE segment_name=$1)", v).Scan(&segmentExists)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if segmentExists {
			existingSegments = append(existingSegments, v)
		}
	}
	if len(existingSegments) == 0 {
		return fmt.Errorf("%s: %w", op, ErrSegmentsNotExists)
	}

	stmt, err := s.db.Prepare("INSERT INTO user_segments(user_id, segment_name) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	for _, v := range existingSegments {
		_, err := stmt.Exec(user, v)
		if err != nil {
			if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code == "23505" {
				continue
			}
		}
	}

	return nil
}

// Delete Segments for User
func (s *Storage) DeleteSegmFromUser(user int, segments []string) error {
	const op = "storage.deletesegmentsfromuser"

	var userExists bool

	if err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)",
		user).Scan(&userExists); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !userExists {
		return fmt.Errorf("%s: %w", op, ErrUserNotExists)
	}

	existingSegments := make([]string, 0, len(segments))
	for _, v := range segments {
		var segmentExists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM segments WHERE segment_name=$1)", v).Scan(&segmentExists)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if segmentExists {
			existingSegments = append(existingSegments, v)
		}
	}
	if len(existingSegments) == 0 {
		return fmt.Errorf("%s: %w", op, ErrSegmentsNotExists)
	}

	stmt, err := s.db.Prepare("DELETE FROM user_segments WHERE user_id=$1 AND segment_name=$2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	for _, v := range existingSegments {
		_, err := stmt.Exec(user, v)
		if err != nil {
			if sqlErr, ok := err.(*pq.Error); ok && sqlErr.Code == "23505" {
				continue
			}
		}
	}

	return nil
}

// Get User Info
func (s *Storage) GetUser(user int) (segments []string, err error) {
	const op = "storage.getuser"

	var userExists bool

	if err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)",
		user).Scan(&userExists); err != nil {
		return segments, fmt.Errorf("%s: %w", op, err)
	}
	if !userExists {
		return segments, fmt.Errorf("%s: %w", op, ErrUserNotExists)
	}

	rows, err := s.db.Query("SELECT segment_name FROM user_segments WHERE user_id = $1", user)
	if err != nil {
		return segments, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var segmentName string
		if err := rows.Scan(&segmentName); err != nil {
			return segments, fmt.Errorf("%s: %w", op, ErrSegmentsNotExists)
		}
		segments = append(segments, segmentName)
	}
	if err := rows.Err(); err != nil {
		return segments, fmt.Errorf("%s: %w", op, err)
	}

	return segments, nil
}
